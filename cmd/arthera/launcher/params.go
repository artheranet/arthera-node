package launcher

import (
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/ethereum/go-ethereum/params"

	"github.com/artheranet/arthera-node/opera"
	"github.com/artheranet/arthera-node/opera/genesis"
	"github.com/artheranet/arthera-node/opera/genesisstore"
)

var (
	Bootnodes = map[string][]string{
		"main": {
			"enode://cab84d0b50098ddc7429da3fceb34fa2eadb694cb45bd8e0ff869b2eadb6c0af92a0347235674bb349a8535fa435c0135034bd20b47394c5306d83b83c718631@34.242.220.16:5050",
			"enode://128f3c3d79c195d9612690e83991aaac575bedd4381b8b0df12fea1570c9c1c798fb22f9a893fa93d108c16d062356bf6e90e6a310007affdeeeff72f288bc3a@3.35.200.210:5050",
			"enode://8a0819390475f3dd0dc1056a9304ab7080894031d13abff6f0dbc9fd620896cdd1167d2c00dd51ae8491d434403a6a2717adc00c62993d6698c328091a8c47e4@3.35.200.210:5050",
		},
		"test": {
			"enode://563b30428f48357f31c9d4906ca2f3d3815d663b151302c1ba9d58f3428265b554398c6fabf4b806a49525670cd9e031257c805375b9fdbcc015f60a7943e427@3.213.142.230:7946",
			"enode://8b53fe4410cde82d98d28697d56ccb793f9a67b1f8807c523eadafe96339d6e56bc82c0e702757ac5010972e966761b1abecb4935d9a86a9feed47e3e9ba27a6@3.227.34.226:7946",
			"enode://1703640d1239434dcaf010541cafeeb3c4c707be9098954c50aa705f6e97e2d0273671df13f6e447563e7d3a7c7ffc88de48318d8a3cc2cc59d196516054f17e@52.72.222.228:7946",
		},
	}

	mainnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash("0x4a53c5445584b3bfc20dbfb2ec18ae20037c716f3ba2d9e1da768a9deca17cb4"),
		NetworkID:   opera.MainNetworkID,
		NetworkName: "main",
	}

	testnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash("0xe5f49e8fc4aa98637e388612eee3fc807c503a2258925e8695f790da678f94c5"),
		NetworkID:   opera.TestNetworkID,
		NetworkName: "test",
	}

	AllowedOperaGenesis = []GenesisTemplate{
		{
			Name:   "Mainnet",
			Header: mainnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash("0x78f9a29621bedba751074ecc810c4f44af13effa76355a41d29663eb89ff9554"),
				genesisstore.BlocksSection: hash.HexToHash("0x29caa135cc7c4d7ec6fd38d473334057e8b25db5375b0f9848f566dac7da4d38"),
				genesisstore.EvmSection:    hash.HexToHash("0x9f87b77720ea44dedbafd0f1908d9d9f26c790ad9996fd8f3216bfe054224796"),
			},
		},

		{
			Name:   "Testnet",
			Header: testnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash("0xba2d180d8d978509eb10dce539b5998b408824595c10b7dd635fead6a53fa290"),
				genesisstore.BlocksSection: hash.HexToHash("0xdb774b0378602202558c7ebada4a0ada5662fa43046def9677c7fdbb29d37c29"),
				genesisstore.EvmSection:    hash.HexToHash("0xf9d0b0a942d56b6e3901d445726cb8e441387258e085ecb7a6f6de8529b1032f"),
			},
		},
	}
)

func overrideParams() {
	params.MainnetBootnodes = []string{}
	params.RopstenBootnodes = []string{}
	params.RinkebyBootnodes = []string{}
	params.GoerliBootnodes = []string{}
}
