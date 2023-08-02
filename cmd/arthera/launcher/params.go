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
		"dev": {
			"enode://2195b2c0ca9695cff2fe84851e6213f2a231d21d056572f6e9bc52b800299beea846ac47ab37256200638aa574fb0387df006638eb0026eaac8ef8a5a3f7b604@49.13.19.99:6535",
			"enode://a50aca72719cec5e2a1f8cd2361adb8eb363c99968be9b633277e20fd86b702a7b3cc7fecaefaf99200023808a01aba98301b16a931081f49b877d48563db510@168.119.169.170:6535",
			"enode://03c117cb3bca6902c31923d99e1cc6afa8aa16cfbd9d1daa262da8b38b96ba2e716e100d1f58b2f6bc26d1b332dd750ea813dd83291038b41659a97adaf88e16@65.108.216.225:6535",
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

	devnetHeader = genesis.Header{
		GenesisID:   hash.HexToHash("0x23254992e8c6d5e3c69c0218268cd702de79b69c3471045f031dd703a8ecca81"),
		NetworkID:   params.DevNetworkID,
		NetworkName: "dev",
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

		{
			Name:   "Devnet",
			Header: devnetHeader,
			Hashes: genesis.Hashes{
				genesisstore.EpochsSection: hash.HexToHash("0xd5307bcd76919384992cdc041d2bf6cff80c4358424028a3e5ddfabab76a3630"),
				genesisstore.BlocksSection: hash.HexToHash("0x1451833512ca813c33fedce593ce60dbf17ea6e9a8e746fc1597d8b9e6dd7c43"),
				genesisstore.EvmSection:    hash.HexToHash("0xd1c78b1d64284fae15ff18fa87d006db36317675747e3e47ab8c171c48ba3b11"),
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
