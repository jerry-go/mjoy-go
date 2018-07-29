package variabledb

import (
	"mjoy.io/utils/database"
	"mjoy.io/common/types"
	"mjoy.io/trie"
	"sync"
	"github.com/hashicorp/golang-lru"
	"fmt"
)

// Trie cache generation limit after which to evic trie nodes from memory.
var MaxTrieCacheGen = uint16(120)

const (
	// Number of past tries to keep. This value is chosen such that
	// reasonable chain reorg depths will hit an existing trie.
	maxPastTries = 12

	// Number of codehash->size associations to keep.
	valuesSizeCacheSize = 100000  //100 000
)

// Database wraps access to tries and contract variable.
type Database interface {
	// Accessing tries:
	// OpenTrie opens the main account trie.
	// OpenStorageTrie opens the storage trie of an account.
	openTrie(root types.Hash) (Trie, error)
	openGroupTrie(addrHash, root types.Hash) (Trie, error)

	// Accessing action result:
	actionResult(keyHash, valHash types.Hash) ([]byte, error)
	actionResultSize(keyHash, valHash types.Hash) (int, error)

	// CopyTrie returns an independent copy of the given trie.
	copyTrie(Trie) Trie
}

// Trie is a Mjoy Merkle Trie.
type Trie interface {
	TryGet(key []byte) ([]byte, error)
	TryUpdate(key, value []byte) error
	TryDelete(key []byte) error
	CommitTo(trie.DatabaseWriter) (types.Hash, error)
	Hash() types.Hash
	NodeIterator(startKey []byte) trie.NodeIterator
	GetKey([]byte) []byte
}

type cachingDB struct {
	db		      database.IDatabase
	mu            sync.Mutex

	pastTries     []*trie.SecureTrie
	valuesSizeCache *lru.Cache
}

// NewDatabase creates a backing store for state. The returned database is safe for
// concurrent use and retains cached trie nodes in memory.
func NewDatabase(db database.IDatabase) Database {
	csc, _ := lru.New(valuesSizeCacheSize)
	return &cachingDB{db: db, valuesSizeCache: csc}
}

func (db *cachingDB) openTrie(root types.Hash) (Trie, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i := len(db.pastTries) - 1; i >= 0; i-- {
		if db.pastTries[i].Hash() == root {
			return cachedTrie{db.pastTries[i].Copy(), db}, nil
		}
	}
	tr, err := trie.NewSecure(root, db.db, MaxTrieCacheGen)
	if err != nil {
		return nil, err
	}
	return cachedTrie{tr, db}, nil
}

func (db *cachingDB) pushTrie(t *trie.SecureTrie) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.pastTries) >= maxPastTries {
		copy(db.pastTries, db.pastTries[1:])
		db.pastTries[len(db.pastTries)-1] = t
	} else {
		db.pastTries = append(db.pastTries, t)
	}
}

func (db *cachingDB) openGroupTrie(addrHash, root types.Hash) (Trie, error) {
	return trie.NewSecure(root, db.db, 0)
}

func (db *cachingDB) copyTrie(t Trie) Trie {
	switch t := t.(type) {
	case cachedTrie:
		return cachedTrie{t.SecureTrie.Copy(), db}
	case *trie.SecureTrie:
		return t.Copy()
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
}

func (db *cachingDB) actionResult(keyHash, valHash types.Hash) ([]byte, error) {
	val, err := db.db.Get(keyHash[:])
	if err == nil {
		db.valuesSizeCache.Add(valHash, len(val))
	}
	return val, err
}

func (db *cachingDB) actionResultSize(keyHash, valHash types.Hash) (int, error) {
	if cached, ok := db.valuesSizeCache.Get(valHash); ok {
		return cached.(int), nil
	}
	val, err := db.actionResult(keyHash, valHash)
	if err == nil {
		db.valuesSizeCache.Add(valHash, len(val))
	}
	return len(val), err
}

// cachedTrie inserts its trie into a cachingDB on commit.
type cachedTrie struct {
	*trie.SecureTrie
	db *cachingDB
}

func (m cachedTrie) CommitTo(dbw trie.DatabaseWriter) (types.Hash, error) {
	root, err := m.SecureTrie.CommitTo(dbw)
	if err == nil {
		m.db.pushTrie(m.SecureTrie)
	}
	return root, err
}
