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
func NewStateTransition(statedb *state.StateDB, msg Message, coinBase types.Address) *StateTransition {
	return &StateTransition{
		msg:      msg,
		actions:msg.Actions(),
		statedb:  statedb,
		coinBase: coinBase,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
func ApplyMessage(statedb *state.StateDB, msg Message, coinBase types.Address) ([]byte, bool, error) {
	return NewStateTransition(statedb, msg, coinBase).TransitionDb()
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

// TransitionDb will transition the state by applying the current message and
// returning the result. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, failed bool, err error) {
	if err = st.preCheck(); err != nil {
		return
	}

	return  []byte{1,2,3},true , nil

	//if false{
	//
	//	msg := st.msg
	//	sender := st.from() // err checked in preCheck
	//
	//	contractCreation := msg.To() == nil
	//
	//
	//	// Snapshot !!!!!!!!!!!!!!!!!
	//	revid := st.statedb.Snapshot()
	//	if contractCreation {
	//		// TODO:
	//
	//		// RevertToSnapshot !!!!!!!!!!!!!!!!!
	//		st.statedb.RevertToSnapshot(revid)
	//		logger.Warnf("Not support create contraction.")
	//		return nil, true, fmt.Errorf("Not support create contraction.")
	//	} else {
	//		// TODO:
	//		logger.Debugf("Just process simple transaction.")
	//
	//		// Increment the nonce for the next transaction
	//		st.statedb.SetNonce(sender, st.statedb.GetNonce(sender)+1)
	//		// skip, direct success
	//
	//		// TODO: for test
	//		{
	//			fee := new(big.Int).Div(st.msg.Value(), new(big.Int).SetUint64(1000))
	//			if fee.Cmp(new(big.Int).SetUint64(1)) < 0 {
	//				fee.SetUint64(1)
	//			}
	//			if st.statedb.GetBalance(sender).Cmp(fee.Add(fee, st.value)) < 0 {
	//				// RevertToSnapshot !!!!!!!!!!!!!!!!!
	//				st.statedb.RevertToSnapshot(revid)
	//				return nil, true, fmt.Errorf("Insufficient balance(addr: %x).", sender)
	//			}
	//
	//			st.statedb.AddBalance(*msg.To(), st.value)
	//		}
	//	}
	//	/*if vmerr != nil {
	//		logger.Debug("VM returned with error", "err", vmerr)
	//		// The only possible consensus-error would be if there wasn't
	//		// sufficient balance to make the transfer happen. The first
	//		// balance transfer may never fail.
	//		if vmerr == vm.ErrInsufficientBalance {
	//			return nil, 0, false, vmerr
	//		}
	//	}*/
	//
	//	// TODO: just for test, fee per transaction
	//	fee := new(big.Int).Div(st.msg.Value(), new(big.Int).SetUint64(1000))
	//	if fee.Cmp(new(big.Int).SetUint64(1)) < 0 {
	//		fee.SetUint64(1)
	//	}
	//	st.statedb.AddBalance(st.coinBase, fee)
	//
	//	st.statedb.SubBalance(sender, fee.Add(fee,st.value))
	//
	//	return ret, false, err
	//
	//}

}
