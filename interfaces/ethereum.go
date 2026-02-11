// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package interfaces

import (
	"context"
	"errors"
	"math/big"

	"github.com/holiman/uint256"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	types2 "github.com/ledgerwatch/erigon-lib/types"
	"starlink-world/erigon-evm/core/types"
)

// NotFound is returned by API methods if the requested item does not exist.
var NotFound = errors.New("not found")

// Subscription represents an event subscription where events are
// delivered on a data channel.
type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}

// ChainEthReader provides access to the blockchain for high-level Ethereum API.
type ChainEthReader interface {
	BlockByHash(ctx context.Context, hash libcommon.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByHash(ctx context.Context, hash libcommon.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionCount(ctx context.Context, blockHash libcommon.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash libcommon.Hash, index uint) (*types.Transaction, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (Subscription, error)
}

// TransactionReader provides access to past transactions and their receipts.
type TransactionReader interface {
	TransactionByHash(ctx context.Context, txHash libcommon.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash libcommon.Hash) (*types.Receipt, error)
}

// ChainStateReader wraps access to the state trie of the canonical blockchain.
type ChainStateReader interface {
	BalanceAt(ctx context.Context, account libcommon.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account libcommon.Address, key libcommon.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account libcommon.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account libcommon.Address, blockNumber *big.Int) (uint64, error)
}

// SyncProgress gives progress indications when the node is synchronising with the Ethereum network.
type SyncProgress struct {
	StartingBlock uint64
	CurrentBlock  uint64
	HighestBlock  uint64
	PulledStates  uint64
	KnownStates   uint64
}

// ChainSyncReader wraps access to the node's current sync status.
type ChainSyncReader interface {
	SyncProgress(ctx context.Context) (*SyncProgress, error)
}

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From       libcommon.Address
	To         *libcommon.Address
	Gas        uint64
	GasPrice   *uint256.Int
	Value      *uint256.Int
	Data       []byte
	FeeCap     *uint256.Int
	Tip        *uint256.Int
	AccessList types2.AccessList
}

// ContractCaller provides contract calls.
type ContractCaller interface {
	CallContract(ctx context.Context, call CallMsg, blockNumber *big.Int) ([]byte, error)
}

// FilterQuery contains options for contract log filtering.
type FilterQuery struct {
	BlockHash *libcommon.Hash
	FromBlock *big.Int
	ToBlock   *big.Int
	Addresses []libcommon.Address
	Topics    [][]libcommon.Hash
}

// LogFilterer provides access to contract log events.
type LogFilterer interface {
	FilterLogs(ctx context.Context, q FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q FilterQuery, ch chan<- types.Log) (Subscription, error)
}

// TransactionSender wraps transaction sending.
type TransactionSender interface {
	SendTransaction(ctx context.Context, tx *types.Transaction) error
}

// GasPricer wraps the gas price oracle.
type GasPricer interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

// PendingStateReader provides access to the pending state.
type PendingStateReader interface {
	PendingBalanceAt(ctx context.Context, account libcommon.Address) (*big.Int, error)
	PendingStorageAt(ctx context.Context, account libcommon.Address, key libcommon.Hash) ([]byte, error)
	PendingCodeAt(ctx context.Context, account libcommon.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account libcommon.Address) (uint64, error)
	PendingTransactionCount(ctx context.Context) (uint, error)
}

// PendingContractCaller can be used to perform calls against the pending state.
type PendingContractCaller interface {
	PendingCallContract(ctx context.Context, call CallMsg) ([]byte, error)
}

// GasEstimator wraps EstimateGas.
type GasEstimator interface {
	EstimateGas(ctx context.Context, call CallMsg) (uint64, error)
}

// PendingStateEventer provides access to real time notifications about changes to the pending state.
type PendingStateEventer interface {
	SubscribePendingTransactions(ctx context.Context, ch chan<- *types.Transaction) (Subscription, error)
}
