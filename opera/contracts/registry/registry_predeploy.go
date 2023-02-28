package registry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GetContractBin is the Staking contract genesis implementation bin code
// Has to be compiled with flag bin-runtime
// Built from arthera-contracts main, solc 0.5.17+commit.d19bba13, optimize-runs 10000
func GetContractBin() []byte {
	return hexutil.MustDecode("0x6080604052348015600f57600080fd5b506004361060285760003560e01c8063669d8d4514602d575b600080fd5b605d60048036036020811015604157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166086565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b73ffffffffffffffffffffffffffffffffffffffff908116600090815260208190526040902054169056fea265627a7a7231582001f8c4db9a32f0a05afa99168b5c7b887d07b7a8abe3a6b519ff3901776314ce64736f6c63430005110032")
}

// ContractAddress is the Registry contract address
var ContractAddress = common.HexToAddress("0xfc00face00000000000000000000000000000002")
