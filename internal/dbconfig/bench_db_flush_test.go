package dbconfig

import (
	"github.com/artheranet/arthera-node/genesis/fake"
	"io/ioutil"
	"os"
	"testing"

	"github.com/artheranet/lachesis/abft"
	"github.com/artheranet/lachesis/hash"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/artheranet/lachesis/utils/cachescale"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artheranet/arthera-node/gossip"
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/artheranet/arthera-node/internal/vecmt"
	"github.com/artheranet/arthera-node/utils"
)

func BenchmarkFlushDBs(b *testing.B) {
	dir := tmpDir("flush_bench")
	defer os.RemoveAll(dir)
	genStore := fake.FakeGenesisStore(1, utils.ToArt(1), utils.ToArt(1))
	g := genStore.Genesis()
	_, _, store, s2, _, closeDBs := MakeEngine(dir, &g, Configs{
		Opera:         gossip.DefaultConfig(cachescale.Identity),
		OperaStore:    gossip.DefaultStoreConfig(cachescale.Identity),
		Lachesis:      abft.DefaultConfig(),
		LachesisStore: abft.DefaultStoreConfig(cachescale.Identity),
		VectorClock:   vecmt.DefaultConfig(cachescale.Identity),
		DBs:           DefaultDBsConfig(cachescale.Identity.U64, 512),
	})
	defer closeDBs()
	defer store.Close()
	defer s2.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		n := idx.Block(0)
		randUint32s := func() []uint32 {
			arr := make([]uint32, 128)
			for i := 0; i < len(arr); i++ {
				arr[i] = uint32(i) ^ (uint32(n) << 16) ^ 0xd0ad884e
			}
			return []uint32{uint32(n), uint32(n) + 1, uint32(n) + 2}
		}
		for !store.IsCommitNeeded() {
			store.SetBlock(n, &inter.Block{
				Time:        inter.Timestamp(n << 32),
				Atropos:     hash.Event{},
				Events:      hash.Events{},
				Txs:         []common.Hash{},
				InternalTxs: []common.Hash{},
				SkippedTxs:  randUint32s(),
				GasUsed:     uint64(n) << 24,
				Root:        hash.Hash{},
			})
			n++
		}
		b.StartTimer()
		err := store.Commit()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func tmpDir(name string) string {
	dir, err := ioutil.TempDir("", name)
	if err != nil {
		panic(err)
	}
	return dir
}
