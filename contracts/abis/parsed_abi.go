package abis

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

var (
	NodeDriver           = mustParseAbi("NodeDriver", NodeDriverStr)
	NetworkInitializer   = mustParseAbi("NetworkInitializer", NetworkInitializerStr)
	EVMWriter            = mustParseAbi("EVMWriter", EVMWriterStr)
	ISubscribers         = mustParseAbi("ISubscribers", ISubscribersStr)
	PayAsYouGoGasRewards = mustParseAbi("PayAsYouGoGasRewards", PayAsYouGoGasRewardsStr)
	NativeToken          = mustParseAbi("NativeToken", NativeTokenStr)
	Staking              = mustParseAbi("Staking", StakingStr)
)

func mustParseAbi(name, abiStr string) *abi.ABI {
	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		panic(fmt.Sprintf("Error reading ABI %s err=%s", name, err))
	}
	return &parsedAbi
}
