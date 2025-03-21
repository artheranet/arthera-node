package brstream

import (
	"fmt"

	"github.com/artheranet/lachesis/inter/idx"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/artheranet/lachesis/gossip/basestream"
)

type Request struct {
	Session   Session
	Limit     Metric
	Type      basestream.RequestType
	MaxChunks uint32
}

type Response struct {
	SessionID uint32
	Done      bool
	Payload   []rlp.RawValue
}

type Session struct {
	ID    uint32
	Start Locator
	Stop  Locator
}

type Locator idx.Block

func (l Locator) Compare(b basestream.Locator) int {
	if l == b.(Locator) {
		return 0
	}
	if l < b.(Locator) {
		return -1
	}
	return 1
}

func (l Locator) Inc() basestream.Locator {
	return l + 1
}

type Payload struct {
	Items []rlp.RawValue
	Keys  []Locator
	Size  uint64
}

func (p *Payload) AddFullBlockRecords(id Locator, brsB rlp.RawValue) {
	p.Items = append(p.Items, brsB)
	p.Keys = append(p.Keys, id)
	p.Size += uint64(len(brsB))
}

func (p Payload) Len() int {
	return len(p.Keys)
}

func (p Payload) TotalSize() uint64 {
	return p.Size
}

func (p Payload) TotalMemSize() int {
	return int(p.Size) + len(p.Keys)*32
}

type Metric struct {
	Num  idx.Block
	Size uint64
}

func (m Metric) String() string {
	return fmt.Sprintf("{Num=%d,Size=%d}", m.Num, m.Size)
}
