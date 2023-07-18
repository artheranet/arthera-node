package asyncflushproducer

import (
	"github.com/artheranet/lachesis/kvdb"
)

type store struct {
	kvdb.Store
	CloseFn func() error
}

func (s *store) Close() error {
	return s.CloseFn()
}
