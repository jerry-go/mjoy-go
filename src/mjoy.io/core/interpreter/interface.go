package interpreter

import (
	"mjoy.io/params"
	"mjoy.io/common/types"
	"mjoy.io/core/state"
	"mjoy.io/utils/crypto"
	"errors"
	"mjoy.io/core/transaction"
)

type ActionResult struct {
	Key []byte
	Val []byte
}

type MemDatabase struct {
	Key []byte
	Val []byte
}

type ActionResults []ActionResult


var (
	ErrContractAddressCollision = errors.New("contract address collision")
	ErrContractCodeSizeTooLong = errors.New("contract address code size too long")
)

type vms struct {}
// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

func  Create(sender types.Address, stateDb *state.StateDB, actions transaction.ActionSlice) ( actionReuslts []ActionResult, contractAddr types.Address, err error) {

	// Ensure there's no existing contract already at the designated address
	nonce := stateDb.GetNonce(sender)
	stateDb.SetNonce(sender, nonce+1)

	contractAddr = crypto.CreateAddress(sender, nonce)
	contractHash := stateDb.GetCodeHash(contractAddr)
	if stateDb.GetNonce(contractAddr) != 0 || (contractHash != (types.Hash{}) && contractHash != emptyCodeHash) {
		return nil, types.Address{}, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := stateDb.Snapshot()
	stateDb.CreateAccount(contractAddr)
	stateDb.SetNonce(contractAddr, 1)

	if len(actions) != 2 {
		stateDb.RevertToSnapshot(snapshot)
		return nil, types.Address{}, ErrContractCodeSizeTooLong
	}

	//todo fee transfer

	//2. save contract code
	if len(actions[1].Params) > params.MaxCodeSize {
		stateDb.RevertToSnapshot(snapshot)
		return nil, types.Address{}, ErrContractCodeSizeTooLong
	}
	stateDb.SetCode(contractAddr, actions[1].Params)

	return ActionResults{}, contractAddr, nil
}


