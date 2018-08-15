package db

import (
	"github.com/golovers/leveltable"
)

// NewLDBDatabase returns a LevelDB wrapped object.
func NewLDBDatabase(file string, cache int, handles int) (Database, error) {
	return leveltable.NewLevelTableDB(file, cache, handles)
}
