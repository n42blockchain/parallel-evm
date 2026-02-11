package main

import (
	"flag"
	_ "net/http/pprof"
	"os"
	"path"

	"github.com/c2h5oh/datasize"
	libdir "github.com/ledgerwatch/erigon-lib/common/datadir"

	"starlink-world/erigon-evm/eth/ethconfig"
	"starlink-world/erigon-evm/internal/app"
	"starlink-world/erigon-evm/log"
	"starlink-world/erigon-evm/params"
	"starlink-world/erigon-evm/stage"
)

var (
	datadir          = "/Users/mac/devel/500s"
	startBlockNumber = uint64(0000000)
	stopBlockNumber  = uint64(4990000)
	chaindata        = path.Join(datadir, "chaindata")
	batchSizeStr     = "512M"
	chainName        string
)

func main() {
	pathp := flag.String("p", "", "datadir")
	Startnum := flag.Uint64("b", 0, "begin block number")
	Stopnum := flag.Uint64("e", 0, "end block number")
	flag.Parse()
	if *pathp != "" {
		datadir = *pathp
		chaindata = path.Join(datadir, "chaindata")
	}
	if *Startnum != 0 {
		startBlockNumber = *Startnum
	}
	if *Stopnum != 0 {
		stopBlockNumber = *Stopnum
	}

	dirs := libdir.New(datadir)
	logger := log.New()
	logger.Info("Build info", "git_branch", params.GitBranch, "git_tag", params.GitTag, "git_commit", params.GitCommit)
	logger.Info("Program starting", "args", os.Args)
	db := app.OpenDB(dirs.Chaindata, logger, true, 2)
	defer db.Close()
	ctx, _ := app.RootContext()
	var batchSize datasize.ByteSize
	app.Must(batchSize.UnmarshalText([]byte(batchSizeStr)))
	pm, engine, chainConfig, vmConfig := app.NewSync(ctx, db, nil, chainName, batchSizeStr)
	cfg := stage.StageExecuteBlocksCfg(db, pm, batchSize, nil, chainConfig, engine, vmConfig, nil, false, false, false, dirs, app.GetBlockReader(db, datadir), nil, nil, ethconfig.Sync{}, nil)
	stage.SpawnExecuteBlocksStage(startBlockNumber, stopBlockNumber, ctx, cfg)
}
