package db

// IdealBatchSize Code using batches should try to add this much data to the batch.
// The value was determined empirically.
const IdealBatchSize = 100 * 1024

var db Database

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Putter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close()
	NewBatch() Batch
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type Batch interface {
	Putter
	ValueSize() int // amount of data in the batch
	Write() error
	Reset() // Reset resets the batch for reuse
}

// SetDatabase set database to be used
func SetDatabase(database Database) {
	db = database
}

func Put(key []byte, value []byte) error {
	return db.Put(key, value)
}

func Get(key []byte) ([]byte, error) {
	return db.Get(key)
}

func Has(key []byte) (bool, error) {
	return db.Has(key)
}

func Delete(key []byte) error {
	return db.Delete(key)
}

func Close() {
	db.Close()
}

func NewBatch() Batch {
	return db.NewBatch()
}
