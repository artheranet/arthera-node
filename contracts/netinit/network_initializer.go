package netinit

import (
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/utils"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func InitializeAll(
	sealedEpoch idx.Epoch,
	totalSupply *big.Int,
	owner common.Address,
) []byte {
	data, _ := abis.NetworkInitializer.Pack(
		"initializeAll",
		utils.U64toBig(uint64(sealedEpoch)),
		totalSupply,
		owner,
	)
	return data
}
