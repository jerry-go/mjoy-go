package state

import (
	"mjoy.io/trie"
	"mjoy.io/common/types"
	"bytes"
	"github.com/tinylib/msgp/msgp"
)

// NewStateSync create a new state trie download scheduler.
func NewStateSync(root types.Hash, database trie.DatabaseReader) *trie.TrieSync {
	var syncer *trie.TrieSync
	callback := func(leaf []byte, parent types.Hash) error {
		var obj Account
		byteBuf := bytes.NewBuffer(leaf)
		if err := msgp.Decode(byteBuf, &obj); err != nil {
			return err
		}
		syncer.AddSubTrie(obj.Root, 64, parent, nil)
		syncer.AddRawEntry(types.BytesToHash(obj.CodeHash), 64, parent)
		return nil
	}
	syncer = trie.NewTrieSync(root, database, callback)
	return syncer
}
