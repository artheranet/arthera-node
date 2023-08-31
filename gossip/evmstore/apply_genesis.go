package evmstore

import (
	"github.com/artheranet/arthera-node/utils/adapters/ethdb2kvdb"
	"github.com/artheranet/arthera-node/utils/dbutil/autocompact"
	"github.com/artheranet/lachesis/kvdb/batched"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/artheranet/arthera-node/genesis"
)

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(g genesis.Genesis) (err error) {
	db := batched.Wrap(autocompact.Wrap(autocompact.Wrap(ethdb2kvdb.Wrap(s.EvmDb), 1*opt.GiB), 16*opt.GiB))
	g.RawEvmItems.ForEach(func(key, value []byte) bool {
		err = db.Put(key, value)
		if err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return db.Write()
}

func (s *Store) WrapTablesAsBatched() (unwrap func()) {
	origTables := s.table

	batchedTxs := batched.Wrap(s.table.Txs)
	s.table.Txs = batchedTxs

	batchedTxPositions := batched.Wrap(s.table.TxPositions)
	s.table.TxPositions = batchedTxPositions

	unwrapLogs := s.EvmLogs.WrapTablesAsBatched()

	batchedReceipts := batched.Wrap(autocompact.Wrap(autocompact.Wrap(s.table.Receipts, opt.GiB), 16*opt.GiB))
	s.table.Receipts = batchedReceipts
	return func() {
		_ = batchedTxs.Flush()
		_ = batchedTxPositions.Flush()
		_ = batchedReceipts.Flush()
		unwrapLogs()
		s.table = origTables
	}
}
