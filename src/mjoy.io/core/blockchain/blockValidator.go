package blockchain

import (
	"fmt"
	"mjoy.io/core/state"
	"errors"
	"mjoy.io/core/transaction"
	"mjoy.io/consensus"
	"mjoy.io/core/blockchain/block"
	"mjoy.io/core"
	"mjoy.io/params"
)

// BlockValidator implements Validator.
type BlockValidator struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(config *params.ChainConfig, blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody verifies the the block
// header's transaction root. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(blk *block.Block) error {
	// Check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(blk.Hash()) {
		return core.ErrKnownBlock
	}
	if !v.bc.HasBlockAndState(blk.ParentHash()) {
		if !v.bc.HasBlock(blk.ParentHash(), blk.NumberU64()-1) {
			return errors.New("unknown ancestor")
		}
		return errors.New("pruned ancestor")
	}
	// Header validity is known at this point
	header := blk.Header()

	if hash := block.DeriveSha(blk.Transactions()); hash != header.TxHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash.String(), header.TxHash.String())
	}
	return nil
}

// ValidateState validates the various changes that happen after a state
// transition, such as the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(blk, parent *block.Block, statedb *state.StateDB, receipts transaction.Receipts) error {
	header := blk.Header()

	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := transaction.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := block.DeriveSha(receipts)
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash.String(), receiptSha.String())
	}
	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(); header.StateHash != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.StateHash.String(), root.String())
	}
	return nil
}