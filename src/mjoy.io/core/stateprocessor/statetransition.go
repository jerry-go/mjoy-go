////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: statetransition.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package stateprocessor

import (
	"mjoy.io/common/types"
	"mjoy.io/core/state"
	"mjoy.io/core"
	"mjoy.io/core/transaction"
	"mjoy.io/core/interpreter"
	"mjoy.io/utils/crypto"
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	msg        Message
	actions     []transaction.Action
	statedb    *state.StateDB
	coinBase   types.Address
	Cache      *DbCache
}

// Message represents a message sent to a contract.
type Message interface {
	From() types.Address
	To() *types.Address
	Actions()[]transaction.Action
	Nonce() uint64
	CheckNonce() bool
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(statedb *state.StateDB, msg Message, coinBase types.Address, cache *DbCache) *StateTransition {
	return &StateTransition{
		msg:      msg,
		actions:  msg.Actions(),
		statedb:  statedb,
		coinBase: coinBase,
		Cache : cache,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
func ApplyMessage(statedb *state.StateDB, msg Message, coinBase types.Address, cache *DbCache) ([]byte, bool, error) {
	return NewStateTransition(statedb, msg, coinBase, cache).TransitionDb()
}

func (st *StateTransition) from() types.Address {
	f := st.msg.From()
	if !st.statedb.Exist(f) {
		st.statedb.CreateAccount(f)
	}
	return f
}

func (st *StateTransition) to() types.Address {
	if st.msg == nil {
		return types.Address{}
	}
	to := st.msg.To()
	if to == nil {
		return types.Address{} // contract creation
	}

	reference := *to
	if !st.statedb.Exist(*to) {
		st.statedb.CreateAccount(*to)
	}
	return reference
}

func (st *StateTransition) preCheck() error {
	msg := st.msg
	sender := st.from()

	// Make sure this transaction's nonce is correct
	if msg.CheckNonce() {
		nonce := st.statedb.GetNonce(sender)
		if nonce < msg.Nonce() {
			return core.ErrNonceTooHigh
		} else if nonce > msg.Nonce() {
			return core.ErrNonceTooLow
		}
	}
	return nil
}


// make log  function
func MakeLog(address types.Address, results interpreter.ActionResults, blockNumber uint64) *transaction.Log {
	topics := []types.Hash{}
	data := [][]byte{}

	for _, result := range results {
		topics = append(topics,types.BytesToHash(result.Key))
		data = append(data, result.Key)
	}

	return &transaction.Log{
		Address:     address,
		Topics:      topics,
		Data:        data,
		BlockNumber:   blockNumber,
	}
}

// TransitionDb will transition the state by applying the current message and
// returning the result. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, failed bool, err error) {
	if err = st.preCheck(); err != nil {
		return
	}

	//return  []byte{1,2,3},true , nil
	msg := st.msg
	sender := st.from() // err checked in preCheck

	contractCreation := msg.To() == nil

	// Snapshot !!!!!!!!!!!!!!!!!
	snapshot := st.statedb.Snapshot()
	results := interpreter.ActionResults{}
	contractAddr := types.Address{}
	if contractCreation {
		results, contractAddr, err = interpreter.Create(sender, st.statedb, st.actions)
	} else {
		// TODO:
		logger.Debugf("Just process simple transaction.")

	}

	if err != nil {
		st.statedb.RevertToSnapshot(snapshot)
		return nil, true, err
	}

	for _, result := range results {
		storgageKey := append(contractAddr.Bytes(), result.Key...)

		//1, collect results for block producer future write level db
		st.Cache.Cache[string(storgageKey)] = interpreter.MemDatabase{
			contractAddr,
			result.Key,
			result.Val}

		//2. make log for receipt
		//todo here need vm ctx
		log := MakeLog(contractAddr,results, 0)
		st.statedb.AddLog(log)

		//3, change statedb storage
		storageKeyHash := crypto.Keccak256Hash(storgageKey)
		storageValHash := crypto.Keccak256Hash(result.Val)
		st.statedb.SetState(contractAddr, storageKeyHash, storageValHash)
	}

	return ret, false, err

}
