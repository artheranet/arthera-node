package utils

import (
	"github.com/artheranet/lachesis/kvdb"
	"github.com/artheranet/lachesis/kvdb/table"
)

func NewTableOrSelf(db kvdb.Store, prefix []byte) kvdb.Store {
	if len(prefix) == 0 {
		return db
	}
	return table.New(db, prefix)
}
