package stateprocessor

import (
	"mjoy.io/core/blockchain/block"
	"mjoy.io/core/state"
	"mjoy.io/params"
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
	"mjoy.io/utils/bloom"
	"fmt"
	"mjoy.io/consensus"
)

type IChainForState interface {
	consensus.ChainReader
}

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	cs     IChainForState      // chain interface for state processor
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, cs IChainForState, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		cs: cs,
		engine: engine,
		config: config,
	}
}

// Process processes the state changes according to the Mjoy rules by running
// the transaction messages using the statedb and applying any rewards to
// the processor (coinbase).
//
// Process returns the receipts and logs accumulated during the process.
// If any of the transactions failed  it will return an error.
func (p *StateProcessor) Process(block *block.Block, statedb *state.StateDB) (transaction.Receipts, []*transaction.Log, error) {
	var (
		receipts transaction.Receipts
		header   = block.Header()
		allLogs  []*transaction.Log
	)

	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, err := ApplyTransaction(p.config, nil, statedb, header, tx)
		if err != nil {
			logger.Errorf("ApplyTransacton Wrong.....:",err.Error())

			return nil, nil, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}

	// TODO: need to be compeleted, now skip this step
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	if p.engine != nil {
		p.engine.Finalize(p.cs, header, statedb, block.Transactions(), receipts)
	}

	return receipts, allLogs, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, author *types.Address, statedb *state.StateDB, header *block.Header, tx *transaction.Transaction) (*transaction.Receipt, error) {
	msg, err := tx.AsMessage(transaction.MakeSigner(config, &header.Number.IntVal))
	if err != nil {
		return nil, err
	}
	
	// Apply the transaction to the current state (included in the env)
	if author == nil {
		author = &header.Coinbase
	}
	_, failed, err := ApplyMessage(statedb, msg, *author)
	if err != nil {
		return nil, err
	}
	// Update the state with pending changes
	var root []byte
	statedb.Finalise(true)

	// Create a new receipt for the transaction, storing the intermediate root  by the tx
	// based on the mip phase, we're passing wether the root touch-delete accounts.
	receipt := transaction.NewReceipt(root, failed)
	receipt.TxHash = tx.Hash()
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		// TODO:
		logger.Warnf("Not support to create contract!\n")
		return nil, fmt.Errorf("Not support to create contract!")

		//receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())

	topics := []bloom.BloomByte{}
	for _, log := range receipt.Logs {
		topics = append(topics, log.Address)
		for _, topic := range log.Topics {
			topics = append(topics, topic)
		}
	}
	receipt.Bloom = bloom.CreateBloom(topics)

	return receipt, err
}
