package autocompact

import (
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"sync"
)

type Store struct {
	kvdb.Store
	minKey  []byte
	maxKey  []byte
	written uint64
	limit   uint64
	compMu  sync.Mutex
}

type Batch struct {
	kvdb.Batch
	written uint64
	minKey  []byte
	maxKey  []byte
	onWrite func(key []byte, size uint64, force bool)
}

func Wrap(s kvdb.Store, limit uint64) *Store {
	return &Store{
		Store: s,
		limit: limit,
	}
}
