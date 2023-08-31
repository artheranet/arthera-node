package threads

import "github.com/artheranet/lachesis/kvdb"

type limitedProducer struct {
	kvdb.FullDBProducer
}

type limitedStore struct {
	kvdb.Store
}

func Limited(dbs kvdb.FullDBProducer) kvdb.FullDBProducer {
	return &limitedProducer{dbs}
}

func (p *limitedProducer) OpenDB(name string) (kvdb.Store, error) {
	s, err := p.FullDBProducer.OpenDB(name)
	return &limitedStore{s}, err
}
