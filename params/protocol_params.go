package params

import "github.com/ethereum/go-ethereum/common"

var (
	Bootnodes = map[string][]string{
		"main": {
			"enode://c4138ac86c0cd7607231601afeb24d2d0b7aaf4c8e1de28978aab6fdbca1af2ec41307b33529e64da72b9e592782840aa1552bda5b1f700eefbdba2a662c4ecc@bootnode1.arthera.net:6533",
			"enode://91a72a01ba3eb7994c681d2e0212a8239d51fffe9dfdb38af373af2cb9612bf1f2f9514fea1b099a2d1e2bb83afa9c4487af7c7f05bcb5f0c6c63f8f757c7fcd@bootnode2.arthera.net:6533",
			"enode://c3c05932654c6a92a8f544857491073f96715fa394476edee07931bf4da00d76eeb1b256d0a68d0404a8b37e8dec7e038c18f1a4cdbde0da91ef45d04c243de1@bootnode3.arthera.net:6533",
			"enode://0ccae9dd3d8f6033929c9cf8371c4dd96f27c602e2cf9b363196a22557f8f36ca02b6ee64d20d230cb6b60d4a473066f9a4e6b0551d6ca67c797250c9bb3db20@bootnode4.arthera.net:6533",
			"enode://946ec02182a11adb2c93851303ca8c6a6ef5041d84e7f695f3dd6a574f1891b879c5bde1c19cafa1ad91a9057e7d6830b8def76dede45bda350a9deb4e0c152c@bootnode5.arthera.net:6533",
		},
		"test": {
			"enode://b8069d60bd2c6992b3c53740aec0dd21fe8737921094f857f24e73a926148b7a3da7f4259b956da94b8d55f3a606f91ae8ac571e3078bf1759e78ea1be02fc2b@57.129.12.242:6534",
			"enode://bfd1f3076b9d8df8ea24874ad588ba6ac5e6ba2e9029bd4e6748f9e38d4503c62aeb649bd6df842dc1c79f60919156cf472d5cf1a4b6ae60288cf21523f5ea87@51.178.46.242:6534",
			"enode://09294690223b242ef4225a2b8fec8f0d1183d94d4135ae573ab50b1f6c593b5fc4cf81a1c4a7d27debd07c53eb371fef08b63b12b80b48ddcbe3c5fb9bbb8ab3@57.129.13.140:6534",
		},
		"dev": {
			"enode://0c202f27f0b3f0cff2542a639216bfac72c7900f4949bb30ac165c1021cfa905d8b6b4a714314456af3a01a4df1617c66bda9f4f35ea99da8049a1b76d2725d2@167.235.247.140:6535",
			"enode://8ac57f12ca2962454f2196c372d9708a8926dd012c6ef901d128d721faeb8efbbb89c9e81d1c472ac00bbbf01b5680626990ef2caedbf317b51faec1c13b246f@49.13.80.149:6535",
			"enode://f5415312d0e205956e89fcd8d3aabe437d15f0d328153df4604ea135e8a0955ff40a439e2b8432c837e30b572c0ec8e30ad3e42e5f6a1e1515553160cbc6f220@159.69.240.199:6535",
		},
	}
)

var (
	MainnetGenesisID            = "0x993417afae15968aa185376537e2c633844b5e624444812b64e0241a3289d8be"
	MainnetGenesisEpochsSection = "0x3ab151817b1a2f204394ab9ff1fbb98a3fd633cb1baa613d5c62b60b03b3e655"
	MainnetGenesisBlocksSection = "0xc4c14fcc2367bee44d422f35b82d643eff26038ffa251ea5ca83fc28e0090805"
	MainnetGenesisEvmSection    = "0x104d1761f3ab6f6be312e2adeb2a6640048104dbe647d17018d23ea0ab2cfe9d"

	TestnetGenesisID            = "0x9b213544c349bd0e44b6890f7db7e40c64298925dbb9a69529e9fd0f31fbf337"
	TestnetGenesisEpochsSection = "0xd7e51f497d950b35f7e2d03b383d25c204b7780ce1822dcf70404c040338292f"
	TestnetGenesisBlocksSection = "0xd1d7afa34a7f09983dab3b5c5bb9697302f476502ad9440d04118fe3868eeea7"
	TestnetGenesisEvmSection    = "0x6babd9234ed59ca9c7ee5eb285a048d1b916e736f46f043ea1d7b920c44f798a"

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
