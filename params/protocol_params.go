package params

var (
	thousand = uint64(1000)
	million  = thousand * thousand

	MaxGasForHasActiveSubscription = 500 * thousand
	MaxGasForDebitSubscription     = 500 * thousand
	MaxGasForCreditSubscription    = 500 * thousand
	MaxGasForGetSub                = 500 * thousand
)
