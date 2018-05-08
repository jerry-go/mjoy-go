package consensus

import (
	"mjoy.io/core/blockchain/block"
	"mjoy.io/common/types"
	"mjoy.io/core/state"
	"mjoy.io/core/transaction"
	"math/big"
)

type Engine_empty struct {
}


func (empty *Engine_empty) Author(header *block.Header) (types.Address, error) {
	return header.Coinbase, nil
}

func (empty *Engine_empty) VerifyHeader(chain ChainReader, header *block.Header, seal bool) error {
	return nil
}

func (empty *Engine_empty) VerifyHeaders(chain ChainReader, headers []*block.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort, results := make(chan struct{}), make(chan error, len(headers))
	for i := 0; i < len(headers); i++ {
		results <- nil
	}
	return abort, results
}

func (empty *Engine_empty) VerifySeal(chain ChainReader, header *block.Header) error {
	return nil
}

func (empty *Engine_empty) Prepare(chain ChainReader, header *block.Header) error {
	return nil
}

func (empty *Engine_empty) Finalize(chain ChainReader, header *block.Header, state *state.StateDB, txs []*transaction.Transaction, receipts []*transaction.Receipt) (*block.Block, error) {
	reward := big.NewInt(5e+18)
	state.AddBalance(header.Coinbase, reward)
	header.StateHash = state.IntermediateRoot()
	return block.NewBlock(header, txs, receipts), nil
}

func (empty *Engine_empty) Seal(chain ChainReader, block *block.Block, stop <-chan struct{}) (*block.Block, error){
	header := block.Header()
	return block.WithSeal(header), nil
}
