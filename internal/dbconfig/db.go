package dbconfig

import (
	"github.com/artheranet/arthera-node/utils/dbutil/dbcounter"
	"github.com/artheranet/lachesis/kvdb/leveldb"
	"io"
	"os"
	"path"
	"strings"

	"github.com/artheranet/lachesis/hash"
	"github.com/artheranet/lachesis/inter/dag"
	"github.com/artheranet/lachesis/kvdb"
	"github.com/artheranet/lachesis/kvdb/flaggedproducer"
	"github.com/artheranet/lachesis/kvdb/flushable"
	"github.com/artheranet/lachesis/kvdb/multidb"
	"github.com/artheranet/lachesis/kvdb/pebble"
	"github.com/artheranet/lachesis/utils/fmtfilter"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"

	"github.com/artheranet/arthera-node/gossip"
	"github.com/artheranet/arthera-node/utils/dbutil/asyncflushproducer"
)

type DBsConfig struct {
	Routing      RoutingConfig
	RuntimeCache DBsCacheConfig
	GenesisCache DBsCacheConfig
}

type DBCacheConfig struct {
	Cache   uint64
	Fdlimit uint64
}

type DBsCacheConfig struct {
	Table map[string]DBCacheConfig
}

func SupportedDBs(chaindataDir string, cfg DBsCacheConfig) (map[multidb.TypeName]kvdb.IterableDBProducer, map[multidb.TypeName]kvdb.FullDBProducer) {
	if chaindataDir == "inmemory" || chaindataDir == "" {
		chaindataDir, _ = os.MkdirTemp("", "arthera-tmp")
	}
	cacher, err := dbCacheFdlimit(cfg)
	if err != nil {
		utils.Fatalf("Failed to create DB cacher: %v", err)
	}

	leveldbFsh := dbcounter.Wrap(leveldb.NewProducer(path.Join(chaindataDir, "leveldb-fsh"), cacher), true)
	pebbleFsh := dbcounter.Wrap(pebble.NewProducer(path.Join(chaindataDir, "pebble-fsh"), cacher), true)
	pebbleFlg := dbcounter.Wrap(pebble.NewProducer(path.Join(chaindataDir, "pebble-flg"), cacher), true)
	pebbleDrc := dbcounter.Wrap(pebble.NewProducer(path.Join(chaindataDir, "pebble-drc"), cacher), true)

	if metrics.Enabled {
		leveldbFsh = WrapDatabaseWithMetrics(leveldbFsh)
		pebbleFsh = WrapDatabaseWithMetrics(pebbleFsh)
		pebbleFlg = WrapDatabaseWithMetrics(pebbleFlg)
		pebbleDrc = WrapDatabaseWithMetrics(pebbleDrc)
	}

	return map[multidb.TypeName]kvdb.IterableDBProducer{
			"leveldb-fsh": leveldbFsh,
			"pebble-fsh":  pebbleFsh,
			"pebble-flg":  pebbleFlg,
			"pebble-drc":  pebbleDrc,
		}, map[multidb.TypeName]kvdb.FullDBProducer{
			"leveldb-fsh": flushable.NewSyncedPool(leveldbFsh, FlushIDKey),
			"pebble-fsh":  asyncflushproducer.Wrap(flushable.NewSyncedPool(pebbleFsh, FlushIDKey), 200000),
			"pebble-flg":  flaggedproducer.Wrap(pebbleFlg, FlushIDKey),
			"pebble-drc":  &DummyScopedProducer{pebbleDrc},
		}
}

func dbCacheFdlimit(cfg DBsCacheConfig) (func(string) (int, int), error) {
	fmts := make([]func(req string) (string, error), 0, len(cfg.Table))
	fmtsCaches := make([]DBCacheConfig, 0, len(cfg.Table))
	exactTable := make(map[string]DBCacheConfig, len(cfg.Table))
	// build scanf filters
	for name, cache := range cfg.Table {
		if !strings.ContainsRune(name, '%') {
			exactTable[name] = cache
		} else {
			fn, err := fmtfilter.CompileFilter(name, name)
			if err != nil {
				return nil, err
			}
			fmts = append(fmts, fn)
			fmtsCaches = append(fmtsCaches, cache)
		}
	}
	return func(name string) (int, int) {
		// try exact match
		if cache, ok := cfg.Table[name]; ok {
			return int(cache.Cache), int(cache.Fdlimit)
		}
		// try regexp
		for i, fn := range fmts {
			if _, err := fn(name); err == nil {
				return int(fmtsCaches[i].Cache), int(fmtsCaches[i].Fdlimit)
			}
		}
		// default
		return int(cfg.Table[""].Cache), int(cfg.Table[""].Fdlimit)
	}, nil
}

func isEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

func dropAllDBs(chaindataDir string) {
	_ = os.RemoveAll(chaindataDir)
}

func dropAllDBsIfInterrupted(chaindataDir string) {
	if isInterrupted(chaindataDir) {
		log.Info("Restarting genesis processing")
		dropAllDBs(chaindataDir)
	}
}

type GossipStoreAdapter struct {
	*gossip.Store
}

func (g *GossipStoreAdapter) GetEvent(id hash.Event) dag.Event {
	e := g.Store.GetEvent(id)
	if e == nil {
		return nil
	}
	return e
}

func MakeDBDirs(chaindataDir string) {
	dbs, _ := SupportedDBs(chaindataDir, DBsCacheConfig{})
	for typ := range dbs {
		if err := os.MkdirAll(path.Join(chaindataDir, string(typ)), 0700); err != nil {
			utils.Fatalf("Failed to create chaindata/leveldb directory: %v", err)
		}
	}
}

type DummyScopedProducer struct {
	kvdb.IterableDBProducer
}

func (d DummyScopedProducer) NotFlushedSizeEst() int {
	return 0
}

func (d DummyScopedProducer) Flush(_ []byte) error {
	return nil
}

func (d DummyScopedProducer) Initialize(_ []string, flushID []byte) ([]byte, error) {
	return flushID, nil
}

func (d DummyScopedProducer) Close() error {
	return nil
}
