package database

type IDatabaseGetter interface {
	Get(key []byte) ([]byte, error)
}

type IDatabasePutter interface {
	Put(key []byte, value []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type IDatabase interface {
	IDatabaseGetter
	IDatabasePutter
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close()
	NewBatch() IBatch
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type IBatch interface {
	IDatabasePutter
	ValueSize() int 	// amount of data in the batch
	Write() error
	Reset() 			// Reset resets the batch for reuse
}
