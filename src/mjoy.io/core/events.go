package core

import (
	"mjoy.io/core/blockchain/block"
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
)

// TxPreEvent is posted when a transaction enters the transaction pool.
type TxPreEvent struct{ Tx *transaction.Transaction}

// PendingLogsEvent is posted pre producing and notifies of pending logs.
type PendingLogsEvent struct {
	Logs []*transaction.Log
}

// PendingStateEvent is posted pre producing and notifies of pending state changes.
type PendingStateEvent struct{}

// NewProducedBlockEvent is posted when a block has been imported.
type NewProducedBlockEvent struct{ Block *block.Block }

// RemovedTransactionEvent is posted when a reorg happens
type RemovedTransactionEvent struct{ Txs transaction.Transactions }

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent struct{ Logs []*transaction.Log }

type ChainEvent struct {
	Block *block.Block
	Hash  types.Hash
	Logs  []*transaction.Log
}

type ChainSideEvent struct {
	Block *block.Block
}

type ChainHeadEvent struct{ Block *block.Block }
