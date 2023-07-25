package launcher

import (
	"github.com/artheranet/arthera-node/params"
	"github.com/artheranet/lachesis/hash"
	ethparams "github.com/ethereum/go-ethereum/params"

	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
)

var (
	Bootnodes = map[string][]string{
		"main": {},
		"test": {
			"enode://4dbc94a60d0d5c91b0fcafd8dd931bb77a2de8b269c80a58da676af3a74fcf9fa5457c536aea40544080780a99b0dcf6629f34f0974d21da7a4c2f62a0074eec@167.235.203.218:6534",
			"enode://eae25d66e8fedd6e32801f627ae8ff35b74478460c3ea5d773bcf1417dc88e75b8284eb8ee9aa2c5d450f785591b63b440ea708a590d355ed57e1f51c7fc3082@135.181.90.176:6534",
			"enode://8cb1c9e3d93a88a9abb56aaa9c4fdd85bf6c90d6145236c696a6fc334890d334e462c406c6e72b2b463ef41844fc065f769fbc628a6371b3d0c70e3636d272bc@5.78.68.153:6534",
		},
	}

	mainnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		NetworkID:   params.MainNetworkID,
		NetworkName: "main",
	}

	testnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash("0x4288e5d835b1c94747ef6a6fb0366ff0dfacd76ea6c418c1028b3fc18b17474d"),
		NetworkID:   params.TestNetworkID,
		NetworkName: "test",
	}

	AllowedArtheraGenesis = []GenesisTemplate{
		{
			Name:   "Mainnet",
			Header: mainnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				genesisstore.BlocksSection: hash.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				genesisstore.EvmSection:    hash.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			},
		},

		{
			Name:   "Testnet",
			Header: testnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash("0x68f968d90f7e1be36f0470799df69f5d879bfb623d53a80a072d5f072af946dc"),
				genesisstore.BlocksSection: hash.HexToHash("0x783e96ca1b7331e889e93864269c5faeffa7dd9f443776ac63ac8605972c16d6"),
				genesisstore.EvmSection:    hash.HexToHash("0xcb0a423dda114ba5a15d0133d28fb80aac34f70fae9fbc2ae75f85455eda1432"),
			},
		},
	}
)

func overrideParams() {
	ethparams.MainnetBootnodes = []string{}
	ethparams.RopstenBootnodes = []string{}
	ethparams.RinkebyBootnodes = []string{}
	ethparams.GoerliBootnodes = []string{}
}
