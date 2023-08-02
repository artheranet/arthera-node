package params

import (
	utils2 "github.com/artheranet/arthera-node/utils"
	"math/big"
)

type GenesisValidator struct {
	Addr    string
	Pubkey  string
	Stake   *big.Int
	Balance *big.Int
}

type GenesisAccount struct {
	Addr    string
	Balance *big.Int
}

var (
	DevnetValidators = []GenesisValidator{
		{
			Addr:    "0xC34ad0296f606749Ff8F77b75A7d577eB2CDF846",
			Pubkey:  "0xc004bf689e0aa508fc18c9820348cea64cc8b3b3dff85af513fef6309a514c21b33d96e6113904e21c49e012cb73c46d1e5b8ab7cad64131b27a8578d9f87a298f49",
			Stake:   utils2.ToArt(100_000), // min stake StakerConstants.sol -> minSelfStake
			Balance: utils2.ToArt(10_000),
		},
		{
			Addr:    "0xE3E641Dad7ac5477482Ab338FD535EAA1071201F",
			Pubkey:  "0xc0046218198298ade0acaecde7816c1513c40c359673b516449f4e383d87fa53b54c245a75ed98629f2a35eae140306d7d81b6ba33feec81f0b98ddb0b529c48db32",
			Stake:   utils2.ToArt(100_000), // min stake StakerConstants.sol -> minSelfStake
			Balance: utils2.ToArt(10_000),
		},
		{
			Addr:    "0x131BA6dA27B622De8771b63D40915F2530fCa0BD",
			Pubkey:  "0xc004c0cc8ddc257ed1aadd6aec58c40592c00ce653c1e96c5856f3cbb57371b01a7e8c0f2b76255271004673930ebcbc798c360f8df39c08e831bc815ee9f13dc6a6",
			Stake:   utils2.ToArt(100_000), // min stake StakerConstants.sol -> minSelfStake
			Balance: utils2.ToArt(10_000),
		},
	}

	DevnetAccounts = []GenesisAccount{
		{
			Addr:    "0x0c08A529D58152A01d20b46B28DEEB7a4075104A",
			Balance: utils2.ToArt(60_000_000),
		},
		{
			Addr:    "0x9C9994Bc1F7086633fCeA94f94c88251B75CA384",
			Balance: utils2.ToArt(60_000_000),
		},
		{
			Addr:    "0x821b76971bD47770d790F2707823ba29c19bb225",
			Balance: utils2.ToArt(60_000_000),
		},
	}

	TestnetValidators = []GenesisValidator{
		{
			Addr:    "0x7a97E50436a074ADDB9A51D50Fbd35ADAFE88442",
			Pubkey:  "0xc0041d7405a8bc7dabf1e397e6689ff09482466aea9d3a716bf1dd4fd971c22d035d8d939c88764136a3213106282887f9005b5addf23af781302a0119400706996e",
			Stake:   utils2.ToArt(1_000_000), // min stake StakerConstants.sol -> minSelfStake
			Balance: utils2.ToArt(0),
		},
		{
			Addr:    "0xfE8301b91A8Eb4734ed954f8E2FB84c2F72Cef8a",
			Pubkey:  "0xc004a61ec5eb3cf8d6b399ff56682b95277337b601fb31e1a254dd451101b8aafb0218d428fc814faee132aabcc17b3dd39fa35dfce2d5ce29d6bd05615bbd571016",
			Stake:   utils2.ToArt(1_000_000), // min stake StakerConstants.sol -> minSelfStake
			Balance: utils2.ToArt(0),
		},
		{
			Addr:    "0xF51e935061731a129765ff63b3Af0Adb5e4486aC",
			Pubkey:  "0xc004c39c38dc49cc4c9b64ea9d817545e713635f808d692f2f500ad801e002c50987e15cf4d9419731adf4cd83edf2207a806685cb2b75c3027d2dcdd78ec126f430",
			Stake:   utils2.ToArt(1_000_000),
			Balance: utils2.ToArt(0), // min stake StakerConstants.sol -> minSelfStake
		},
	}

	TestnetAccounts = []GenesisAccount{
		{
			Addr:    "0x40bd65cfc4D95844704F4b2a2c46a60f6d6CE766",
			Balance: utils2.ToArt(10_000_000),
		},
		{
			Addr:    "0x35E58946b74fDbD9032aed876FC58629A6e65E79",
			Balance: utils2.ToArt(10_000_000),
		},
		{
			Addr:    "0x846032c611697818a31cC090D436664b263C6E54",
			Balance: utils2.ToArt(10_000_000),
		},
	}

	MainnetValidators = []GenesisValidator{}
	MainnetAccounts   = []GenesisAccount{}
)
