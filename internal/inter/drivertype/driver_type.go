package drivertype

import (
	"math/big"

	"github.com/artheranet/lachesis/inter/idx"

	"github.com/artheranet/arthera-node/internal/inter/validatorpk"
)

var (
	// DoublesignBit is set if validator has a confirmed pair of fork events
	DoublesignBit = uint64(1 << 7)
	OkStatus      = uint64(0)
)

// Validator is the node-side representation of Driver validator
type Validator struct {
	Weight *big.Int
	PubKey validatorpk.PubKey
}

// ValidatorAndID is pair Validator + ValidatorID
type ValidatorAndID struct {
	ValidatorID idx.ValidatorID
	Validator   Validator
}
