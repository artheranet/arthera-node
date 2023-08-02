package params

import "github.com/ethereum/go-ethereum/common"

var (
	Bootnodes = map[string][]string{
		"main": {},
		"test": {
			"enode://4dbc94a60d0d5c91b0fcafd8dd931bb77a2de8b269c80a58da676af3a74fcf9fa5457c536aea40544080780a99b0dcf6629f34f0974d21da7a4c2f62a0074eec@167.235.203.218:6534",
			"enode://eae25d66e8fedd6e32801f627ae8ff35b74478460c3ea5d773bcf1417dc88e75b8284eb8ee9aa2c5d450f785591b63b440ea708a590d355ed57e1f51c7fc3082@135.181.90.176:6534",
			"enode://8cb1c9e3d93a88a9abb56aaa9c4fdd85bf6c90d6145236c696a6fc334890d334e462c406c6e72b2b463ef41844fc065f769fbc628a6371b3d0c70e3636d272bc@5.78.68.153:6534",
		},
		"dev": {
			"enode://2195b2c0ca9695cff2fe84851e6213f2a231d21d056572f6e9bc52b800299beea846ac47ab37256200638aa574fb0387df006638eb0026eaac8ef8a5a3f7b604@49.13.19.99:6535",
			"enode://a50aca72719cec5e2a1f8cd2361adb8eb363c99968be9b633277e20fd86b702a7b3cc7fecaefaf99200023808a01aba98301b16a931081f49b877d48563db510@168.119.169.170:6535",
			"enode://03c117cb3bca6902c31923d99e1cc6afa8aa16cfbd9d1daa262da8b38b96ba2e716e100d1f58b2f6bc26d1b332dd750ea813dd83291038b41659a97adaf88e16@65.108.216.225:6535",
		},
	}
)

var (
	MainnetGenesisID            = "0x0000000000000000000000000000000000000000000000000000000000000000"
	MainnetGenesisEpochsSection = "0x0000000000000000000000000000000000000000000000000000000000000000"
	MainnetGenesisBlocksSection = "0x0000000000000000000000000000000000000000000000000000000000000000"
	MainnetGenesisEvmSection    = "0x0000000000000000000000000000000000000000000000000000000000000000"

	TestnetGenesisID            = "0x4288e5d835b1c94747ef6a6fb0366ff0dfacd76ea6c418c1028b3fc18b17474d"
	TestnetGenesisEpochsSection = "0x68f968d90f7e1be36f0470799df69f5d879bfb623d53a80a072d5f072af946dc"
	TestnetGenesisBlocksSection = "0x783e96ca1b7331e889e93864269c5faeffa7dd9f443776ac63ac8605972c16d6"
	TestnetGenesisEvmSection    = "0xcb0a423dda114ba5a15d0133d28fb80aac34f70fae9fbc2ae75f85455eda1432"

	DevnetGenesisID            = "0xe37b5f4893fa4d446c1d2e716f76bd6257590814a4ec22a49786d381e206de39"
	DevnetGenesisEpochsSection = "0x419f919b566250e9f4b37edb8d12e9df6a42d57ebf2683587a2b34398f9dd59c"
	DevnetGenesisBlocksSection = "0x39f9cb2c2bc9ad3e753490f0eacf7a3a3badba95377d35a5f26549f45a8c55b8"
	DevnetGenesisEvmSection    = "0xd0103758494c8d0468b81fe0961902329931a0f763f04f6993fc824decf3b0e3"
)

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
