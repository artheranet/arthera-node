package native_sub

import (
	"bytes"
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	PlanIdCounterSlot = iota
	PriceProviderSlot
	PlanIdToPlanSlot
)

var (
	newSubscriptionEvent                    = crypto.Keccak256Hash([]byte("NewSubscription(uint256,uint256,uint256,uint256,uint256)"))
	renewSubscriptionEvent                  = crypto.Keccak256Hash([]byte("RenewSubscription(uint256,uint256,uint256,uint256,uint256)"))
	switchPlanEvent                         = crypto.Keccak256Hash([]byte("SwitchPlan(uint256,uint256,uint256)"))
	terminateSubscriptionEvent              = crypto.Keccak256Hash([]byte("TerminateSubscription(uint256)"))
	whitelisterAddedEvent                   = crypto.Keccak256Hash([]byte("WhitelisterAdded(uint256,address)"))
	whitelisterRemovedEvent                 = crypto.Keccak256Hash([]byte("WhitelisterRemoved(uint256,address)"))
	whitelistContractAccountSubscriberEvent = crypto.Keccak256Hash([]byte("WhitelistContractAccountSubscriber(uint256,address)"))
	blacklistContractAccountSubscriberEvent = crypto.Keccak256Hash([]byte("BlacklistContractAccountSubscriber(uint256,address)"))
	planCreatedEvent                        = crypto.Keccak256Hash([]byte("PlanCreated(uint256,string,string,uint256,uint256,uint256,uint,uint256,bool)"))
	planUpdatedEvent                        = crypto.Keccak256Hash([]byte("PlanUpdated(uint256,string,string,uint256,uint256,uint256,uint,uint256)"))
	planActivatedEvent                      = crypto.Keccak256Hash([]byte("PlanActivated(uint256)"))
	planDeactivatedEvent                    = crypto.Keccak256Hash([]byte("PlanDeactivated(uint256)"))

	newEOASubscriptionMethodID        []byte
	newContractSubscriptionMethodID   []byte
	renewSubscriptionMethodID         []byte
	switchPlanMethodID                []byte
	terminateSubscriptionMethodID     []byte
	addWhitelisterMethodID            []byte
	removeWhitelisterMethodID         []byte
	whitelistAccountMethodID          []byte
	blacklistAccountMethodID          []byte
	getContractSubscriptionIdMethodID []byte
	isWhitelistedMethodID             []byte
	getCapRemainingMethodID           []byte
	getCapWindowMethodID              []byte
	getBalanceMethodID                []byte
	getStartTimeMethodID              []byte
	getEndTimeMethodID                []byte
	hasActiveSubscriptionMethodID     []byte
	hasSubscriptionMethodID           []byte
	getSubscriptionDataMethodID       []byte
	getSubscriptionTokenIdMethodID    []byte
	getSubscriberByIdMethodID         []byte
	createPlanMethodID                []byte
	updatePlanMethodID                []byte
	getPlansMethodID                  []byte
	getPlanMethodID                   []byte
	setActiveMethodID                 []byte
	setPriceProviderMethodID          []byte
	getPriceProviderMethodID          []byte
	priceInAAMethodID                 []byte
)

func init() {
	for name, constID := range map[string]*[]byte{
		"balanceOf":                 &balanceOfMethodID,
		"ownerOf":                   &ownerOfMethodID,
		"safeTransferFrom":          &safeTransferFromMethodID,
		"transferFrom":              &transferFromMethodID,
		"approve":                   &approveMethodID,
		"getApproved":               &getApprovedMethodID,
		"name":                      &nameMethodID,
		"symbol":                    &symbolMethodID,
		"tokenURI":                  &tokenURIMethodID,
		"newEOASubscription":        &newEOASubscriptionMethodID,
		"newContractSubscription":   &newContractSubscriptionMethodID,
		"renewSubscription":         &renewSubscriptionMethodID,
		"switchPlan":                &switchPlanMethodID,
		"terminateSubscription":     &terminateSubscriptionMethodID,
		"addWhitelister":            &addWhitelisterMethodID,
		"removeWhitelister":         &removeWhitelisterMethodID,
		"whitelistAccount":          &whitelistAccountMethodID,
		"blacklistAccount":          &blacklistAccountMethodID,
		"isWhitelisted":             &isWhitelistedMethodID,
		"getContractSubscriptionId": &getContractSubscriptionIdMethodID,
		"getCapRemaining":           &getCapRemainingMethodID,
		"getCapWindow":              &getCapWindowMethodID,
		"getBalance":                &getBalanceMethodID,
		"getStartTime":              &getStartTimeMethodID,
		"getEndTime":                &getEndTimeMethodID,
		"hasActiveSubscription":     &hasActiveSubscriptionMethodID,
		"hasSubscription":           &hasSubscriptionMethodID,
		"getSubscriptionData":       &getSubscriptionDataMethodID,
		"getSubscriptionTokenId":    &getSubscriptionTokenIdMethodID,
		"getSubscriberById":         &getSubscriberByIdMethodID,
		"createPlan":                &createPlanMethodID,
		"updatePlan":                &updatePlanMethodID,
		"getPlans":                  &getPlansMethodID,
		"getPlan":                   &getPlanMethodID,
		"setActive":                 &setActiveMethodID,
		"setPriceProvider":          &setPriceProviderMethodID,
		"getPriceProvider":          &getPriceProviderMethodID,
		"priceInAA":                 &priceInAAMethodID,
	} {
		method, exist := abis.ISubscribers.Methods[name]
		if !exist {
			panic("unknown ISubscribers method")
		}

		*constID = make([]byte, len(method.ID))
		copy(*constID, method.ID)
	}
}

type PreCompiledContract struct{}

func (_ PreCompiledContract) Run(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if len(input) < 4 {
		return nil, 0, vm.ErrExecutionReverted
	}

	if evm.StateDB.GetCodeSize(contracts.SubscribersSmartContractAddress) == 0 {
		evm.StateDB.SetCode(contracts.SubscribersSmartContractAddress, []byte{0})
	}

	methodId := input[:4]
	args := input[4:]

	erc721 := ERC721{
		evm:                 evm,
		name:                "Arthera Subscription",
		symbol:              "ASUB",
		tokenUri:            "",
		beforeTokenTransfer: beforeNftTransfer,
		contractAddress:     contracts.SubscribersSmartContractAddress,
	}

	result, leftoverGas, err := erc721.ProcessMethod(methodId, caller, args, suppliedGas)
	if err != nil || result != nil {
		return result, leftoverGas, err
	}

	if bytes.Equal(methodId, newEOASubscriptionMethodID) {
		return newEOASubscription(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, newContractSubscriptionMethodID) {
		return newContractSubscription(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, renewSubscriptionMethodID) {
		return renewSubscription(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, switchPlanMethodID) {
		return switchPlan(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, terminateSubscriptionMethodID) {
		return terminateSubscription(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, addWhitelisterMethodID) {
		return addWhitelister(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, removeWhitelisterMethodID) {
		return removeWhitelister(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, whitelistAccountMethodID) {
		return whitelistAccount(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, blacklistAccountMethodID) {
		return blacklistAccount(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, isWhitelistedMethodID) {
		return isWhitelisted(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getContractSubscriptionIdMethodID) {
		return getContractSubscriptionId(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getCapRemainingMethodID) {
		return getCapRemaining(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getCapWindowMethodID) {
		return getCapWindow(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getBalanceMethodID) {
		return getBalance(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getStartTimeMethodID) {
		return getStartTime(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getEndTimeMethodID) {
		return getEndTime(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, hasActiveSubscriptionMethodID) {
		return hasActiveSubscription(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, hasSubscriptionMethodID) {
		return hasSubscription(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getSubscriptionDataMethodID) {
		return getSubscriptionData(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getSubscriptionTokenIdMethodID) {
		return getSubscriptionTokenId(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getSubscriberByIdMethodID) {
		return getSubscriberById(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, createPlanMethodID) {
		return createPlan(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, updatePlanMethodID) {
		return updatePlan(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getPlansMethodID) {
		return getPlans(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getPlanMethodID) {
		return getPlan(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, setActiveMethodID) {
		return setPlanActive(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, setPriceProviderMethodID) {
		return setPriceProvider(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getPriceProviderMethodID) {
		return getPriceProvider(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, priceInAAMethodID) {
		return priceInAA(evm, caller, args, suppliedGas)
	} else {
		return nil, 0, vm.ErrExecutionReverted
	}
}

func beforeNftTransfer(state vm.StateDB, address common.Address, address2 common.Address, b *big.Int) []byte {
	/**
	if (to != address(0)) {
		require(_subscriptionsById[getSubscriptionTokenId(to)].planId == 0, "the receiver already has a subscription");
	}
	*/
	return nil
}

// function newEOASubscription(uint256 planId) payable
func newEOASubscription(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function newContractSubscription(address _contract, uint256 planId) payable
func newContractSubscription(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function renewSubscription() payable;
func renewSubscription(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function switchPlan(uint256 newPlanId) payable;
func switchPlan(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function terminateSubscription()
func terminateSubscription(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function addWhitelister(address _contract, address whitelister)
func addWhitelister(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function removeWhitelister(address _contract, address whitelister)
func removeWhitelister(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function whitelistAccount(address _contract, address account)
func whitelistAccount(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function blacklistAccount(address _contract, address account)
func blacklistAccount(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function isWhitelisted(address subscriber, address account) returns (bool)
func isWhitelisted(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getContractSubscriptionId(address _contract) returns (uint256);
func getContractSubscriptionId(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getCapRemaining(address subscriber) returns (uint256);
func getCapRemaining(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getCapWindow(address subscriber) returns (uint256);
func getCapWindow(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getBalance(address subscriber) returns (uint256)
func getBalance(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getStartTime(address subscriber) returns (uint256)
func getStartTime(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getEndTime(address subscriber) returns (uint256)
func getEndTime(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function hasActiveSubscription(address subscriber) returns (bool)
func hasActiveSubscription(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function hasSubscription(address subscriber) returns (bool)
func hasSubscription(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getSubscriptionData(address subscriber) returns (Subscription memory)
func getSubscriptionData(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getSubscriptionTokenId(address subscriber) returns (uint256);
func getSubscriptionTokenId(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getSubscriberById(uint256 id) returns (address);
func getSubscriberById(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getPlans() returns (Plan[] memory)
func getPlans(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getPlan(uint256 id) returns (Plan memory)
func getPlan(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function setPriceProvider(address provider) external
func setPriceProvider(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function getPriceProvider(address provider) external returns (address)
func getPriceProvider(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function priceInAA(address provider) external returns (uint256)
func priceInAA(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function createPlan(
//
//	    string calldata name,
//	    string calldata description,
//	    uint256 duration,
//	    uint256 units,
//	    uint256 usdPrice,
//	    uint capFrequency,
//	    uint256 capUnits,
//	    bool forContract
//	)
func createPlan(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function updatePlan(
//
//	    uint256 planId,
//	    string calldata name,
//	    string calldata description,
//	    uint256 duration,
//	    uint256 units,
//	    uint256 usdPrice,
//	    uint capFrequency,
//	    uint256 capUnits
//	)
func updatePlan(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

// function setActive(uint256 planId, bool active)
func setPlanActive(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return nil, 0, vm.ErrExecutionReverted
}

func _getPlanById(state vm.StateDB, planId *big.Int) (*SubscriptionPlan, error) {
	return nil, nil
}

func _setPlanById(state vm.StateDB, planId *big.Int, plan SubscriptionPlan) error {
	return nil
}
