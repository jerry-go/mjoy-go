package state

import (
	"math/big"
	"mjoy.io/common/types"
)

type journalEntry interface {
	undo(*StateDB)
}

type journal []journalEntry

type (
	// Changes to the account trie.
	createObjectChange struct {
		account *types.Address
	}
	resetObjectChange struct {
		prev *stateObject
	}
	suicideChange struct {
		account     *types.Address
		prev        bool // whether account had already suicided
		prevbalance *big.Int
	}

	// Changes to individual accounts.
	balanceChange struct {
		account *types.Address
		prev    *big.Int
	}
	nonceChange struct {
		account *types.Address
		prev    uint64
	}
	storageChange struct {
		account       *types.Address
		key, prevalue types.Hash
	}
	codeChange struct {
		account            *types.Address
		prevcode, prevhash []byte
	}

	// Changes to other state values.
	refundChange struct {
		prev uint64
	}
	addLogChange struct {
		txhash types.Hash
	}
	addPreimageChange struct {
		hash types.Hash
	}
	touchChange struct {
		account   *types.Address
		prev      bool
		prevDirty bool
	}
)

func (ch createObjectChange) undo(s *StateDB) {
	delete(s.stateObjects, *ch.account)
	delete(s.stateObjectsDirty, *ch.account)
}

func (ch resetObjectChange) undo(s *StateDB) {
	s.setStateObject(ch.prev)
}

func (ch suicideChange) undo(s *StateDB) {
	obj := s.getStateObject(*ch.account)
	if obj != nil {
		obj.suicided = ch.prev
		obj.setBalance(types.NewBigInt(*ch.prevbalance))
	}
}

var ripemd = types.HexToAddress("0000000000000000000000000000000000000003")

func (ch touchChange) undo(s *StateDB) {
	if !ch.prev && *ch.account != ripemd {
		s.getStateObject(*ch.account).touched = ch.prev
		if !ch.prevDirty {
			delete(s.stateObjectsDirty, *ch.account)
		}
	}
}

func (ch balanceChange) undo(s *StateDB) {
	s.getStateObject(*ch.account).setBalance(types.NewBigInt(*ch.prev))
}

func (ch nonceChange) undo(s *StateDB) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch codeChange) undo(s *StateDB) {
	s.getStateObject(*ch.account).setCode(types.BytesToHash(ch.prevhash), ch.prevcode)
}

func (ch storageChange) undo(s *StateDB) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevalue)
}

func (ch refundChange) undo(s *StateDB) {
	s.refund = ch.prev
}

func (ch addLogChange) undo(s *StateDB) {
	logs := s.logs[ch.txhash]
	if len(logs) == 1 {
		delete(s.logs, ch.txhash)
	} else {
		s.logs[ch.txhash] = logs[:len(logs)-1]
	}
	s.logSize--
}

func (ch addPreimageChange) undo(s *StateDB) {
	delete(s.preimages, ch.hash)
}
