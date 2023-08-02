package launcher

import (
	"github.com/artheranet/arthera-node/params"
	"github.com/artheranet/lachesis/hash"
	ethparams "github.com/ethereum/go-ethereum/params"

	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
)

var (
	mainnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash(params.MainnetGenesisID),
		NetworkID:   params.MainNetworkID,
		NetworkName: "main",
	}

	testnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash(params.TestnetGenesisID),
		NetworkID:   params.TestNetworkID,
		NetworkName: "test",
	}

	devnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash(params.DevnetGenesisID),
		NetworkID:   params.DevNetworkID,
		NetworkName: "dev",
	}

	AllowedArtheraGenesis = []GenesisTemplate{
		{
			Name:   "Mainnet",
			Header: mainnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash(params.MainnetGenesisEpochsSection),
				genesisstore.BlocksSection: hash.HexToHash(params.MainnetGenesisBlocksSection),
				genesisstore.EvmSection:    hash.HexToHash(params.MainnetGenesisEvmSection),
			},
		},

		{
			Name:   "Testnet",
			Header: testnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash(params.TestnetGenesisEpochsSection),
				genesisstore.BlocksSection: hash.HexToHash(params.TestnetGenesisBlocksSection),
				genesisstore.EvmSection:    hash.HexToHash(params.TestnetGenesisEvmSection),
			},
		},

		{
			Name:   "Devnet",
			Header: devnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash(params.DevnetGenesisEpochsSection),
				genesisstore.BlocksSection: hash.HexToHash(params.DevnetGenesisBlocksSection),
				genesisstore.EvmSection:    hash.HexToHash(params.DevnetGenesisEvmSection),
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
