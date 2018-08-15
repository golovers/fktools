package db

import (
	"github.com/golovers/leveltable"
)

// IdealBatchSize Code using batches should try to add this much data to the batch.
// The value was determined empirically.
const IdealBatchSize = 100 * 1024

var db Database

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	leveltable.Database
}

// SetDatabase set database to be used
func SetDatabase(database Database) {
	db = database
}

func Table(name string) Database {
	return db.Table(name)
}

func Close() {
	db.Close()
}
