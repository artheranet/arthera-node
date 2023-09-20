package abis

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

var (
	NodeDriver           *abi.ABI = mustParseAbi("NodeDriver", NodeDriverStr)
	NetworkInitializer   *abi.ABI = mustParseAbi("NetworkInitializer", NetworkInitializerStr)
	EVMWriter            *abi.ABI = mustParseAbi("EVMWriter", EVMWriterStr)
	Subscribers          *abi.ABI = mustParseAbi("Subscribers", SubscribersStr)
	PayAsYouGoGasRewards *abi.ABI = mustParseAbi("PayAsYouGoGasRewards", PayAsYouGoGasRewardsStr)
	IERC20WithMetadata   *abi.ABI = mustParseAbi("IERC20WithMetadata", IERC20WithMetadataStr)
	Staking              *abi.ABI = mustParseAbi("Staking", StakingStr)
)

func mustParseAbi(name, abiStr string) *abi.ABI {
	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		panic(fmt.Sprintf("Error reading ABI %s err=%s", name, err))
	}
	return &parsedAbi
}
