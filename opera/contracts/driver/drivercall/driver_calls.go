package drivercall

import (
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/artheranet/arthera-node/opera/contracts/abis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/artheranet/arthera-node/inter"
	"github.com/artheranet/arthera-node/opera"
	"github.com/artheranet/arthera-node/opera/genesis/gpos"
	"github.com/artheranet/arthera-node/utils"
)

type Delegation struct {
	Address            common.Address
	ValidatorID        idx.ValidatorID
	Stake              *big.Int
	LockedStake        *big.Int
	LockupFromEpoch    idx.Epoch
	LockupEndTime      idx.Epoch
	LockupDuration     uint64
	EarlyUnlockPenalty *big.Int
	Rewards            *big.Int
}

// Methods

func SealEpochValidators(_validators []idx.ValidatorID) []byte {
	validators := make([]*big.Int, len(_validators))
	for i, v := range _validators {
		validators[i] = utils.U64toBig(uint64(v))
	}
	data, _ := abis.NodeDriver.Pack("sealEpochValidators", validators)
	return data
}

type ValidatorEpochMetric struct {
	Missed          opera.BlocksMissed
	Uptime          inter.Timestamp
	OriginatedTxFee *big.Int
}

func SealEpoch(metrics []ValidatorEpochMetric) []byte {
	offlineTimes := make([]*big.Int, len(metrics))
	offlineBlocks := make([]*big.Int, len(metrics))
	uptimes := make([]*big.Int, len(metrics))
	originatedTxFees := make([]*big.Int, len(metrics))
	log.Debug("------------ Sealing Epoch -----------")
	for i, m := range metrics {
		offlineTimes[i] = utils.U64toBig(uint64(m.Missed.Period.Unix()))
		offlineBlocks[i] = utils.U64toBig(uint64(m.Missed.BlocksNum))
		uptimes[i] = utils.U64toBig(uint64(m.Uptime.Unix()))
		originatedTxFees[i] = m.OriginatedTxFee

		log.Debug("---> Epoch Validator",
			"index", i,
			"offlineTimes", uint64(m.Missed.Period.Unix()),
			"offlineBlocks", uint64(m.Missed.BlocksNum),
			"uptimes", uint64(m.Uptime.Unix()),
			"originatedTxFees", m.OriginatedTxFee.String(),
		)
	}

	data, _ := abis.NodeDriver.Pack("sealEpoch", offlineTimes, offlineBlocks, uptimes, originatedTxFees)
	return data
}

func SetGenesisValidator(v gpos.Validator) []byte {
	data, _ := abis.NodeDriver.Pack(
		"setGenesisValidator",
		v.Address,
		utils.U64toBig(uint64(v.ID)),
		v.PubKey.Bytes(),
		utils.U64toBig(v.Status),
		utils.U64toBig(uint64(v.CreationEpoch)),
		utils.U64toBig(uint64(v.CreationTime.Unix())),
		utils.U64toBig(uint64(v.DeactivatedEpoch)),
		utils.U64toBig(uint64(v.DeactivatedTime.Unix())),
	)
	return data
}

func SetGenesisDelegation(d Delegation) []byte {
	data, _ := abis.NodeDriver.Pack(
		"setGenesisDelegation",
		d.Address,
		utils.U64toBig(uint64(d.ValidatorID)),
		d.Stake,
		d.LockedStake,
		utils.U64toBig(uint64(d.LockupFromEpoch)),
		utils.U64toBig(uint64(d.LockupEndTime)),
		utils.U64toBig(d.LockupDuration),
		d.EarlyUnlockPenalty,
		d.Rewards,
	)
	return data
}

func DeactivateValidator(validatorID idx.ValidatorID, status uint64) []byte {
	data, _ := abis.NodeDriver.Pack("deactivateValidator", utils.U64toBig(uint64(validatorID)), utils.U64toBig(status))
	return data
}
