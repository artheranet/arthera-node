package params

var (
	thousand = uint64(1000)
	million  = thousand * thousand

	MaxGasForHasActiveSubscription = 500 * thousand
	MaxGasForReduceBalance         = 1 * million
)
