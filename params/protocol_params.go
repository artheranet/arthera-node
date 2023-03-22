package params

var (
	thousand = uint64(1000)
	million  = thousand * thousand

	MaxGasForHasActiveSubscription = 1 * million
	MaxGasForDebitSubscription     = 1 * million
	MaxGasForCreditSubscription    = 1 * million
	MaxGasForGetSub                = 1 * million
)
