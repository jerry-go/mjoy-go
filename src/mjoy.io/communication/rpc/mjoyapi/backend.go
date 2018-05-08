// Package mjoyapi implements the general MJOY API functions.
package mjoyapi

import (
	"context"

	"mjoy.io/node/services/mjoy/downloader"
	"mjoy.io/utils/event"
	"mjoy.io/accounts"
	"mjoy.io/core/state"
	"mjoy.io/core"
	"mjoy.io/params"
	"mjoy.io/utils/database"
	"mjoy.io/communication/rpc"
	"mjoy.io/core/blockchain/block"
	"mjoy.io/common/types"
	"mjoy.io/core/transaction"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	// General mjoy API
	Downloader() *downloader.Downloader
	ProtocolVersion() int
	ChainDb() database.IDatabase
	EventMux() *event.TypeMux
	AccountManager() *accounts.Manager

	// BlockChain API
	SetHead(number uint64)
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Header, error)
	BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*block.Block, error)
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *block.Header, error)
	GetBlock(ctx context.Context, blockHash types.Hash) (*block.Block, error)
	GetReceipts(ctx context.Context, blockHash types.Hash) (transaction.Receipts, error)
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription

	// TxPool API
	SendTx(ctx context.Context, signedTx *transaction.Transaction) error
	GetPoolTransactions() (transaction.Transactions, error)
	GetPoolTransaction(txHash types.Hash) *transaction.Transaction
	GetPoolNonce(ctx context.Context, addr types.Address) (uint64, error)
	Stats() (pending int, queued int)
	TxPoolContent() (map[types.Address]transaction.Transactions, map[types.Address]transaction.Transactions)
	SubscribeTxPreEvent(chan<- core.TxPreEvent) event.Subscription

	ChainConfig() *params.ChainConfig
	CurrentBlock() *block.Block
}

func GetAPIs(apiBackend Backend) []rpc.API {
	nonceLock := new(AddrLocker)
	return []rpc.API{
		{
			Namespace: "mjoy",
			Version:   "1.0",
			Service:   NewPublicMjoyAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "mjoy",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "mjoy",
			Version:   "1.0",
			Service:   NewPublicTransactionPoolAPI(apiBackend, nonceLock),
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "mjoy",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(apiBackend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(apiBackend, nonceLock),
			Public:    false,
		},
	}
}
