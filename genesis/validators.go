package genesis

import (
	"crypto/ecdsa"
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artheranet/arthera-node/internal/inter/validatorpk"
	"github.com/artheranet/lachesis/inter/idx"
)

type (
	// Validator is a helper structure to define genesis validators
	Validator struct {
		ID               idx.ValidatorID
		Address          common.Address
		PubKey           validatorpk.PubKey
		PrivKey          *ecdsa.PrivateKey
		CreationTime     inter.Timestamp
		CreationEpoch    idx.Epoch
		DeactivatedTime  inter.Timestamp
		DeactivatedEpoch idx.Epoch
		Status           uint64
	}

	Validators []Validator
)

// Map converts Validators to map
func (gv Validators) Map() map[idx.ValidatorID]Validator {
	validators := map[idx.ValidatorID]Validator{}
	for _, val := range gv {
		validators[val.ID] = val
	}
	return validators
}

// PubKeys returns not sorted genesis pub keys
func (gv Validators) PubKeys() []validatorpk.PubKey {
	res := make([]validatorpk.PubKey, 0, len(gv))
	for _, v := range gv {
		res = append(res, v.PubKey)
	}
	return res
}

// Addresses returns not sorted genesis addresses
func (gv Validators) Addresses() []common.Address {
	res := make([]common.Address, 0, len(gv))
	for _, v := range gv {
		res = append(res, v.Address)
	}
	return res
}
