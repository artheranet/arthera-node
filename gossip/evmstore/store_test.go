package evmstore

import (
	"github.com/artheranet/lachesis/kvdb/memorydb"
)

func cachedStore() *Store {
	cfg := LiteStoreConfig()

	return NewStore(memorydb.NewProducer(""), cfg)
}

func nonCachedStore() *Store {
	cfg := StoreConfig{}

	return NewStore(memorydb.NewProducer(""), cfg)
}
