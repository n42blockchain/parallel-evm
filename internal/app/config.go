package app

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/c2h5oh/datasize"
	chainp "github.com/ledgerwatch/erigon-lib/chain"
	"github.com/ledgerwatch/erigon-lib/kv"

	"starlink-world/erigon-evm/consensus"
	"starlink-world/erigon-evm/consensus/ethash"
	"starlink-world/erigon-evm/consensus/serenity"
	"starlink-world/erigon-evm/core"
	"starlink-world/erigon-evm/core/types"
	"starlink-world/erigon-evm/core/vm"
	"starlink-world/erigon-evm/eth/ethconfig"
	"starlink-world/erigon-evm/ethdb/prune"
	"starlink-world/erigon-evm/interfaces"
	kv2 "starlink-world/erigon-evm/kv/mdbx"
	"starlink-world/erigon-evm/log"
	"starlink-world/erigon-evm/params"
	"starlink-world/erigon-evm/turbo/snapshotsync"
	"starlink-world/erigon-evm/turbo/snapshotsync/snap"
)

func OpenKV(label kv.Label, logger log.Logger, path string, exclusive bool, databaseVerbosity int) kv.RwDB {
	opts := kv2.NewMDBX(logger).Path(path).Label(label)
	if exclusive {
		opts = opts.Exclusive()
	}
	if databaseVerbosity != -1 {
		opts = opts.DBVerbosity(kv.DBVerbosityLvl(databaseVerbosity))
	}
	return opts.MustOpen()
}

func OpenDB(path string, logger log.Logger, applyMigrations bool, databaseVerbosity int) kv.RwDB {
	label := kv.ChainDB
	db := OpenKV(label, logger, path, false, databaseVerbosity)
	if applyMigrations {
		// Migration logic placeholder
	}
	return db
}

func RootContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		ch := make(chan os.Signal, 1)
		defer close(ch)

		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case sig := <-ch:
			log.Info("Got interrupt, shutting down...", "sig", sig)
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

var openBlockReaderOnce sync.Once
var blockReaderSingleton interfaces.FullBlockReader

func GetBlockReader(db kv.RoDB, datadir string) interfaces.FullBlockReader {
	openBlockReaderOnce.Do(func() {
		blockReaderSingleton = snapshotsync.NewBlockReader()
		if sn := AllSnapshots(db, datadir); sn.Cfg().Enabled {
			x := snapshotsync.NewBlockReaderWithSnapshots(sn, false)
			blockReaderSingleton = x
		}
	})
	return blockReaderSingleton
}

var openSnapshotOnce sync.Once
var allSnapshotsSingleton *snapshotsync.RoSnapshots

func AllSnapshots(db kv.RoDB, datadir string) *snapshotsync.RoSnapshots {
	openSnapshotOnce.Do(func() {
		var useSnapshots bool
		_ = db.View(context.Background(), func(tx kv.Tx) error {
			var err error
			useSnapshots, err = snap.Enabled(tx)
			return err
		})
		snapCfg := ethconfig.NewSnapCfg(useSnapshots, true, true)
		allSnapshotsSingleton = snapshotsync.NewRoSnapshots(snapCfg, filepath.Join(datadir, "snapshots"))
		if useSnapshots {
			if err := allSnapshotsSingleton.ReopenFolder(); err != nil {
				panic(err)
			}
		}
	})
	return allSnapshotsSingleton
}

func ByChain(chain string) (*types.Genesis, *chainp.Config) {
	var chainConfig *chainp.Config
	var genesis *types.Genesis
	if chain == "" {
		chainConfig = params.MainnetChainConfig
		genesis = core.MainnetGenesisBlock()
	} else {
		chainConfig = params.ChainConfigByChainName(chain)
		genesis = core.GenesisBlockByChainName(chain)
	}
	return genesis, chainConfig
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func NewSync(ctx context.Context, db kv.RwDB, miningConfig *params.MiningConfig, chainName string, batchSizeStr string) (prune.Mode, consensus.Engine, *chainp.Config, *vm.Config) {
	var pm prune.Mode
	var err error
	if err = db.View(context.Background(), func(tx kv.Tx) error {
		pm, err = prune.Get(tx)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
	vmConfig := &vm.Config{}

	genesis, chainConfig := ByChain(chainName)
	var engine consensus.Engine
	engine = ethash.NewFaker()
	engine = serenity.New(engine)

	chainConfig, _, genesisErr := core.CommitGenesisBlock(db, genesis, "")
	if _, ok := genesisErr.(*chainp.ConfigCompatError); genesisErr != nil && !ok {
		panic(genesisErr)
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	var batchSize datasize.ByteSize
	Must(batchSize.UnmarshalText([]byte(batchSizeStr)))

	cfg := ethconfig.Defaults
	cfg.Prune = pm
	cfg.BatchSize = batchSize
	if miningConfig != nil {
		cfg.Miner = *miningConfig
	}

	return pm, engine, chainConfig, vmConfig
}
