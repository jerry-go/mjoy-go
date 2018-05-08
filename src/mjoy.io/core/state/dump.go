package state

import (
	"encoding/json"
	"fmt"
	"mjoy.io/common/types/util"
	"mjoy.io/common/types"
	"mjoy.io/trie"
	"bytes"
	"github.com/tinylib/msgp/msgp"
)

type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

func (self *StateDB) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data Account
		byteBuf := bytes.NewBuffer(it.Value)
		if err := msgp.Decode(byteBuf, &data); err!= nil{
			panic(err)
		}

		obj := newObject(nil, types.BytesToAddress(addr), data, nil)
		account := DumpAccount{
			Balance:  data.Balance.IntVal.String(),
			Nonce:    data.Nonce,
			Root:     util.Bytes2Hex(data.Root[:]),
			CodeHash: util.Bytes2Hex(data.CodeHash),
			Code:     util.Bytes2Hex(obj.Code(self.db)),
			Storage:  make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.Storage[util.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = util.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[util.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}
