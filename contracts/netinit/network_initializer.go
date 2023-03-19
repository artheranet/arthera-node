package netinit

import (
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// GetContractBin is NetworkInitializer contract genesis implementation bin code
// Built from arthera-contracts main, solc 0.8.17, optimize-runs 10000
func GetContractBin() []byte {
	return hexutil.MustDecode("0x608060405234801561001057600080fd5b506004361061002b5760003560e01c806333beb9b914610030575b600080fd5b61004361003e366004610336565b610045565b005b6040517f485cc95500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152858116602483015286169063485cc95590604401600060405180830381600087803b1580156100b657600080fd5b505af11580156100ca573d6000803e3d6000fd5b50506040517fc0c53b8b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a81166004830152888116602483015284811660448301528916925063c0c53b8b9150606401600060405180830381600087803b15801561014757600080fd5b505af115801561015b573d6000803e3d6000fd5b50506040517f019e2729000000000000000000000000000000000000000000000000000000008152600481018c9052602481018b905273ffffffffffffffffffffffffffffffffffffffff898116604483015284811660648301528a16925063019e27299150608401600060405180830381600087803b1580156101de57600080fd5b505af11580156101f2573d6000803e3d6000fd5b50506040517f485cc95500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a8116600483015284811660248301528616925063485cc9559150604401600060405180830381600087803b15801561026757600080fd5b505af115801561027b573d6000803e3d6000fd5b50506040517f485cc95500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff898116600483015284811660248301528516925063485cc9559150604401600060405180830381600087803b1580156102f057600080fd5b505af1158015610304573d6000803e3d6000fd5b50600092505050ff5b803573ffffffffffffffffffffffffffffffffffffffff8116811461033157600080fd5b919050565b60008060008060008060008060006101208a8c03121561035557600080fd5b8935985060208a0135975061036c60408b0161030d565b965061037a60608b0161030d565b955061038860808b0161030d565b945061039660a08b0161030d565b93506103a460c08b0161030d565b92506103b260e08b0161030d565b91506103c16101008b0161030d565b9050929598509295985092959856fea264697066735822122008df5ce025831c07ff99d0bc457dde29b9367d018587a6d14d599f67004962eb64736f6c63430008110033")
}

// ContractAddress is the NetworkInitializer contract address
var ContractAddress = common.HexToAddress("0xd1005eed00000000000000000000000000000000")

func InitializeAll(
	sealedEpoch idx.Epoch,
	totalSupply *big.Int,
	sfcAddr common.Address,
	driverAuthAddr common.Address,
	driverAddr common.Address,
	evmWriterAddr common.Address,
	validatorInfoAddr common.Address,
	subscriptionRegistry common.Address,
	owner common.Address,
) []byte {
	data, _ := abis.NetworkInitializer.Pack(
		"initializeAll",
		utils.U64toBig(uint64(sealedEpoch)),
		totalSupply,
		sfcAddr,
		driverAuthAddr,
		driverAddr,
		evmWriterAddr,
		validatorInfoAddr,
		subscriptionRegistry,
		owner,
	)
	return data
}