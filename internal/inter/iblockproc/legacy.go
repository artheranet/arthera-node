package iblockproc

import (
	"github.com/artheranet/lachesis/hash"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/artheranet/lachesis/inter/pos"

	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/artheranet/arthera-node/params"
)

type ValidatorEpochStateV0 struct {
	GasRefund      uint64
	PrevEpochEvent hash.Event
}

type EpochStateV0 struct {
	Epoch          idx.Epoch
	EpochStart     inter.Timestamp
	PrevEpochStart inter.Timestamp

	EpochStateRoot hash.Hash

	Validators        *pos.Validators
	ValidatorStates   []ValidatorEpochStateV0
	ValidatorProfiles ValidatorProfiles

	Rules params.ProtocolRules
}
