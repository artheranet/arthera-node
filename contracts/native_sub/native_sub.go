package native_sub

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Subscription struct {
	Id           uint64
	PlanId       uint64
	Balance      uint64
	StartTime    uint64
	EndTime      uint64
	LastCapReset uint64
	PeriodUsage  uint64
}

type SubscriptionPlan struct {
	PlanId       uint64
	Name         string
	Description  string
	Duration     uint64
	Units        uint64
	Price        uint64
	CapFrequency uint8
	CapUnits     uint64
	ForContract  bool
	Active       bool
}

const (
	CAP_FREQUENCY_NONE = iota
)

var (
	subscriptionPlanAbi = abi.Arguments{
		{Type: abis.AbiUint256},
		{Type: abis.AbiString},
		{Type: abis.AbiString},
		{Type: abis.AbiUint256},
		{Type: abis.AbiUint256},
		{Type: abis.AbiUint256},
		{Type: abis.AbiUint8},
		{Type: abis.AbiUint256},
		{Type: abis.AbiBool},
		{Type: abis.AbiBool},
	}
)

func RetrieveSubscriptionPlan(state vm.StateDB, planId uint64) SubscriptionPlan {
	bigPlanId := new(big.Int).SetUint64(planId)
	storageKey := crypto.Keccak256Hash(
		common.LeftPadBytes(bigPlanId.Bytes(), 32),
		common.LeftPadBytes(new(big.Int).SetUint64(PlanIdToPlanSlot).Bytes(), 32),
	)
	data := state.GetState(contracts.SubscribersSmartContractAddress, storageKey)
	arrayLen := data.Big()
	hashes := make([]common.Hash, arrayLen.Uint64())
	for i := uint64(0); i < arrayLen.Uint64(); i++ {
		skey := crypto.Keccak256Hash(
			storageKey.Bytes(),
			new(big.Int).SetUint64(i).Bytes(),
		)
		hashes[i] = state.GetState(contracts.SubscribersSmartContractAddress, skey)
	}
	return unpackSubscriptionPlan(hashes)
}

func unpackSubscriptionPlan(hashes []common.Hash) SubscriptionPlan {
	buffer := new(bytes.Buffer)
	for _, hash := range hashes {
		buffer.Write(hash.Bytes())
	}
	var plan SubscriptionPlan
	var nameLen = uint64(0)
	var descLen = uint64(0)
	binary.Read(buffer, binary.LittleEndian, &plan.PlanId)
	binary.Read(buffer, binary.LittleEndian, &nameLen)
	var nameBytes = make([]byte, nameLen)
	binary.Read(buffer, binary.LittleEndian, &nameBytes)
	plan.Name = string(nameBytes)
	binary.Read(buffer, binary.LittleEndian, &descLen)
	var descBytes = make([]byte, descLen)
	binary.Read(buffer, binary.LittleEndian, &descBytes)
	plan.Description = string(descBytes)
	binary.Read(buffer, binary.LittleEndian, &plan.Duration)
	binary.Read(buffer, binary.LittleEndian, &plan.Units)
	binary.Read(buffer, binary.LittleEndian, &plan.Price)
	binary.Read(buffer, binary.LittleEndian, &plan.CapFrequency)
	binary.Read(buffer, binary.LittleEndian, &plan.CapUnits)
	binary.Read(buffer, binary.LittleEndian, &plan.ForContract)
	binary.Read(buffer, binary.LittleEndian, &plan.Active)
	return plan
}

func StoreSubscriptionPlan(state vm.StateDB, plan SubscriptionPlan) {
	hashes := packSubscriptionPlan(plan)
	planId := new(big.Int).SetUint64(plan.PlanId)
	storageKey := crypto.Keccak256Hash(
		common.LeftPadBytes(planId.Bytes(), 32),
		common.LeftPadBytes(new(big.Int).SetUint64(PlanIdToPlanSlot).Bytes(), 32),
	)
	arrayLen := new(big.Int).SetUint64(uint64(len(hashes)))
	fmt.Println(arrayLen.String())
	state.SetState(contracts.SubscribersSmartContractAddress, storageKey, common.BigToHash(arrayLen))

	for idx, hash := range hashes {
		skey := crypto.Keccak256Hash(
			storageKey.Bytes(),
			new(big.Int).SetUint64(uint64(idx)).Bytes(),
		)
		state.SetState(contracts.SubscribersSmartContractAddress, skey, hash)
	}
}

func packSubscriptionPlan(plan SubscriptionPlan) []common.Hash {
	buffer := new(bytes.Buffer)
	var nameLen = uint64(len(plan.Name))
	var descLen = uint64(len(plan.Description))
	binary.Write(buffer, binary.LittleEndian, plan.PlanId)
	binary.Write(buffer, binary.LittleEndian, nameLen)
	binary.Write(buffer, binary.LittleEndian, []byte(plan.Name))
	binary.Write(buffer, binary.LittleEndian, descLen)
	binary.Write(buffer, binary.LittleEndian, []byte(plan.Description))
	binary.Write(buffer, binary.LittleEndian, plan.Duration)
	binary.Write(buffer, binary.LittleEndian, plan.Units)
	binary.Write(buffer, binary.LittleEndian, plan.Price)
	binary.Write(buffer, binary.LittleEndian, plan.CapFrequency)
	binary.Write(buffer, binary.LittleEndian, plan.CapUnits)
	binary.Write(buffer, binary.LittleEndian, plan.ForContract)
	binary.Write(buffer, binary.LittleEndian, plan.Active)

	rem := buffer.Len() % 32
	slices := (buffer.Len() - rem) / 32
	if rem > 0 {
		slices++
	}
	hashes := make([]common.Hash, slices)
	for i := 0; i < slices; i++ {
		hashes[i] = common.BytesToHash(buffer.Next(32))
	}
	return hashes
}

func GetSubscriptionData(state vm.StateDB, address common.Address) *Subscription {
	return nil
}

func HasActiveSubscription(state vm.StateDB, subscriber common.Address) bool {
	return false
}

func IsWhitelisted(state vm.StateDB, subscriber common.Address, account common.Address) bool {
	return false
}

func GetCapRemaining(state vm.StateDB, address common.Address) *big.Int {
	return nil
}

func DebitSubscription(state vm.StateDB, target common.Address, units *big.Int) *big.Int {
	if units.BitLen() == 0 {
		return big.NewInt(0)
	}

	return nil
}

func CreditSubscription(state vm.StateDB, target common.Address, units *big.Int) {
	if units.BitLen() == 0 {
		return
	}
}

func SubscriptionDataValid(sub *Subscription) bool {
	return sub != nil && sub.Id > 0 && sub.PlanId > 0 && sub.StartTime > 0 && sub.EndTime > 0
}
