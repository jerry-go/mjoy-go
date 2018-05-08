package database

type table struct {
	db     IDatabase
	prefix string
}

// NewTable returns a Database object that prefixes all keys with a given
// string.
func NewTable(db IDatabase, prefix string) IDatabase {
	return &table{
		db:     db,
		prefix: prefix,
	}
}

func (tb *table) Put(key []byte, value []byte) error {
	return tb.db.Put(append([]byte(tb.prefix), key...), value)
}

func (tb *table) Has(key []byte) (bool, error) {
	return tb.db.Has(append([]byte(tb.prefix), key...))
}

func (tb *table) Get(key []byte) ([]byte, error) {
	return tb.db.Get(append([]byte(tb.prefix), key...))
}

func (tb *table) Delete(key []byte) error {
	return tb.db.Delete(append([]byte(tb.prefix), key...))
}

func (dt *table) Close() {
	// Do nothing; don't close the underlying DB.
}
type tableBatch struct {
	batch  IBatch
	prefix string
}

// NewTableBatch returns a Batch object which prefixes all keys with a given string.
func NewTableBatch(db IDatabase, prefix string) IBatch{
	return &tableBatch{db.NewBatch(), prefix}
}

func (tb *table) NewBatch() IBatch {
	return &tableBatch{tb.db.NewBatch(), tb.prefix}
}

func (tbatch *tableBatch) Put(key, value []byte) error {
	return tbatch.batch.Put(append([]byte(tbatch.prefix), key...), value)
}

func (tbatch *tableBatch) Write() error {
	return tbatch.batch.Write()
}

func (tbatch *tableBatch) ValueSize() int {
	return tbatch.batch.ValueSize()
}

func (tbatch *tableBatch) Reset() {
	tbatch.batch.Reset()
}