package params

import "github.com/ethereum/go-ethereum/common"

var (
	ZeroAddress = common.Address{}
	thousand    = uint64(1000)
	million     = thousand * thousand

	MaxGasForHasActiveSubscription = 500 * thousand
	MaxGasForDebitSubscription     = 500 * thousand
	MaxGasForCreditSubscription    = 500 * thousand
	MaxGasForGetSub                = 500 * thousand
	MaxGasForIsWhitelisted         = 500 * thousand
	MaxGasForSetOwnerOfContract    = 500 * thousand
	MaxGasForAddReward             = 500 * thousand
)
