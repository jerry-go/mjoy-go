package common

import (
	"context"
	"math/big"
	"mjoy.io/common/types"

)

type Subscription interface {
	// Unsubscribe cancels the sending of events to the data channel
	// and closes the error channel.
	Unsubscribe()
	// Err returns the subscription error channel. The error channel receives
	// a value if there is an issue with the subscription (e.g. the network connection
	// delivering the events has been closed). Only one value will ever be sent.
	// The error channel is closed by Unsubscribe.
	Err() <-chan error

}


// ChainStateReader wraps access to the state trie of the canonical blockchain. Note that
// implementations of the interface may be unable to return state values for old blocks.
// In many cases, using CallContract can be preferable to reading raw contract storage.
type ChainStateReader interface {
	BalanceAt(ctx context.Context, account types.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account types.Address, key types.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account types.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account types.Address, blockNumber *big.Int) (uint64, error)
}


// SyncProgress gives progress indications when the node is synchronising with
// the Mjoy network.
type SyncProgress struct {
	StartingBlock uint64 // Block number where sync began
	CurrentBlock  uint64 // Current block number where sync is at
	HighestBlock  uint64 // Highest alleged block number in the chain
	PulledStates  uint64 // Number of state trie entries already downloaded
	KnownStates   uint64 // Total number of state trie entries known about
}
