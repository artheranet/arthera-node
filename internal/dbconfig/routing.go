package dbconfig

import (
	"fmt"
	"github.com/artheranet/arthera-node/utils/dbutil/threads"

	"github.com/artheranet/lachesis/kvdb"
	"github.com/artheranet/lachesis/kvdb/cachedproducer"
	"github.com/artheranet/lachesis/kvdb/multidb"
	"github.com/artheranet/lachesis/kvdb/skipkeys"
)

type RoutingConfig struct {
	Table map[string]multidb.Route
}

func MakeMultiProducer(rawProducers map[multidb.TypeName]kvdb.IterableDBProducer, scopedProducers map[multidb.TypeName]kvdb.FullDBProducer, cfg RoutingConfig) (kvdb.FullDBProducer, error) {
	cachedProducers := make(map[multidb.TypeName]kvdb.FullDBProducer)
	var flushID []byte
	var err error
	for typ, producer := range scopedProducers {
		flushID, err = producer.Initialize(rawProducers[typ].Names(), flushID)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing databases: %v. Try to use 'db heal' to recover", err)
		}
		cachedProducers[typ] = cachedproducer.WrapAll(producer)
	}

	p, err := makeMultiProducer(cachedProducers, cfg)
	return threads.CountedFullDBProducer(p), err
}

func MakeDirectMultiProducer(rawProducers map[multidb.TypeName]kvdb.IterableDBProducer, cfg RoutingConfig) (kvdb.FullDBProducer, error) {
	dproducers := map[multidb.TypeName]kvdb.FullDBProducer{}
	for typ, producer := range rawProducers {
		dproducers[typ] = &DummyScopedProducer{producer}
	}
	return MakeMultiProducer(rawProducers, dproducers, cfg)
}

func makeMultiProducer(scopedProducers map[multidb.TypeName]kvdb.FullDBProducer, cfg RoutingConfig) (kvdb.FullDBProducer, error) {
	multi, err := multidb.NewProducer(scopedProducers, cfg.Table, TablesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to construct multidb: %v", err)
	}

	err = multi.Verify()
	if err != nil {
		return nil, fmt.Errorf("incompatible chainstore DB layout: %v. Try to use 'db transform' to recover", err)
	}
	return skipkeys.WrapAllProducer(multi, MetadataPrefix), nil
}
