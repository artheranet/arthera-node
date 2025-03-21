package ibr

import (
	"github.com/artheranet/lachesis/common/bigendian"
	"github.com/artheranet/lachesis/hash"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/artheranet/arthera-node/internal/inter"
)

type LlrBlockVote struct {
	Atropos      hash.Event
	Root         hash.Hash
	TxHash       hash.Hash
	ReceiptsHash hash.Hash
	Time         inter.Timestamp
	GasUsed      uint64
}

type LlrFullBlockRecord struct {
	Atropos  hash.Event
	Root     hash.Hash
	Txs      types.Transactions
	Receipts []*types.ReceiptForStorage
	Time     inter.Timestamp
	GasUsed  uint64
}

type LlrIdxFullBlockRecord struct {
	LlrFullBlockRecord
	Idx idx.Block
}

func (bv LlrBlockVote) Hash() hash.Hash {
	return hash.Of(bv.Atropos.Bytes(), bv.Root.Bytes(), bv.TxHash.Bytes(), bv.ReceiptsHash.Bytes(), bv.Time.Bytes(), bigendian.Uint64ToBytes(bv.GasUsed))
}

func (br LlrFullBlockRecord) Hash() hash.Hash {
	return LlrBlockVote{
		Atropos:      br.Atropos,
		Root:         br.Root,
		TxHash:       inter.CalcTxHash(br.Txs),
		ReceiptsHash: inter.CalcReceiptsHash(br.Receipts),
		Time:         br.Time,
		GasUsed:      br.GasUsed,
	}.Hash()
}
