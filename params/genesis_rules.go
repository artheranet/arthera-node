package params

import (
	utils2 "github.com/artheranet/arthera-node/utils"
	"math/big"
)

type GenesisValidator struct {
	ID      uint32
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
			ID:      1,
			Addr:    "0xC34ad0296f606749Ff8F77b75A7d577eB2CDF846",
			Pubkey:  "0xc004bf689e0aa508fc18c9820348cea64cc8b3b3dff85af513fef6309a514c21b33d96e6113904e21c49e012cb73c46d1e5b8ab7cad64131b27a8578d9f87a298f49",
			Stake:   utils2.ToArt(5_000_000),
			Balance: utils2.ToArt(60_000_000),
		},
		{
			ID:      2,
			Addr:    "0xE3E641Dad7ac5477482Ab338FD535EAA1071201F",
			Pubkey:  "0xc0046218198298ade0acaecde7816c1513c40c359673b516449f4e383d87fa53b54c245a75ed98629f2a35eae140306d7d81b6ba33feec81f0b98ddb0b529c48db32",
			Stake:   utils2.ToArt(5_000_000),
			Balance: utils2.ToArt(60_000_000),
		},
		{
			ID:      3,
			Addr:    "0x131BA6dA27B622De8771b63D40915F2530fCa0BD",
			Pubkey:  "0xc004c0cc8ddc257ed1aadd6aec58c40592c00ce653c1e96c5856f3cbb57371b01a7e8c0f2b76255271004673930ebcbc798c360f8df39c08e831bc815ee9f13dc6a6",
			Stake:   utils2.ToArt(5_000_000),
			Balance: utils2.ToArt(60_000_000),
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
			ID:      1,
			Addr:    "0x7a97E50436a074ADDB9A51D50Fbd35ADAFE88442",
			Pubkey:  "0xc0041d7405a8bc7dabf1e397e6689ff09482466aea9d3a716bf1dd4fd971c22d035d8d939c88764136a3213106282887f9005b5addf23af781302a0119400706996e",
			Stake:   utils2.ToArt(20_000_000),
			Balance: utils2.ToArt(100_000),
		},
		{
			ID:      2,
			Addr:    "0xfE8301b91A8Eb4734ed954f8E2FB84c2F72Cef8a",
			Pubkey:  "0xc004a61ec5eb3cf8d6b399ff56682b95277337b601fb31e1a254dd451101b8aafb0218d428fc814faee132aabcc17b3dd39fa35dfce2d5ce29d6bd05615bbd571016",
			Stake:   utils2.ToArt(20_000_000),
			Balance: utils2.ToArt(100_000),
		},
		{
			ID:      3,
			Addr:    "0xF51e935061731a129765ff63b3Af0Adb5e4486aC",
			Pubkey:  "0xc004c39c38dc49cc4c9b64ea9d817545e713635f808d692f2f500ad801e002c50987e15cf4d9419731adf4cd83edf2207a806685cb2b75c3027d2dcdd78ec126f430",
			Stake:   utils2.ToArt(20_000_000),
			Balance: utils2.ToArt(100_000),
		},
	}

	TestnetAccounts = []GenesisAccount{
		{
			Addr:    "0x40bd65cfc4D95844704F4b2a2c46a60f6d6CE766",
			Balance: utils2.ToArt(60_000_000),
		},
		{
			Addr:    "0x35E58946b74fDbD9032aed876FC58629A6e65E79",
			Balance: utils2.ToArt(60_000_000),
		},
		{
			Addr:    "0x846032c611697818a31cC090D436664b263C6E54",
			Balance: utils2.ToArt(60_000_000),
		},
	}

	MainnetValidators = []GenesisValidator{
		{
			// Bootnode 1
			ID:      1,
			Addr:    "0x1f25524b9D23320c40CE0C6ffc7E2D5630db5Ceb",
			Pubkey:  "0xc00418209c3d5503b479c383c63456d80e10c6d8844254064b01eb90c13a1486d2071a48a05cb1ab7d93ef80896f9e876430efe9235d80a3ae2f9b2b37a354766944",
			Stake:   utils2.ToArt(3_250_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Bootnode 2
			ID:      2,
			Addr:    "0x37F272A6eFd66AF2d89b16675DBB4410C5fD2AFA",
			Pubkey:  "0xc00463ba1938f62777bcdd5f551cd07a3643c84487b614cd58693690adfcb0c51efa3c49f6ffbccd16e0a95b3a013f7ba9ed80906bb12c51281416c19a77d70c92e3",
			Stake:   utils2.ToArt(3_250_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Bootnode 3
			ID:      3,
			Addr:    "0xBc02ABfA05e2251017C5FD4fD33e59667CC872D3",
			Pubkey:  "0xc004e7a7d0ab393725cc0a581e82d44b3015c3d915c320273eade7fe2b508e360586d44d7c82b8e7c644acf7a92919d587fccd618339a9321de9f0b85d3bef627c5b",
			Stake:   utils2.ToArt(3_250_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Bootnode 4
			ID:      4,
			Addr:    "0x8F69AB452EE5A3c3cD91d15BBc339165A508AceB",
			Pubkey:  "0xc00477944c51faeadeb1063bb60ed4084b42bcf2241446408e121424dc9e4acf36a97e41fc2739025a53e53ebfa98303b1f00640be7c9b677ce85614383de0734c4a",
			Stake:   utils2.ToArt(3_250_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Bootnode 5
			ID:      5,
			Addr:    "0x688Dc34D7176387596EF4D1da78AC5c1FafBF12D",
			Pubkey:  "0xc0043a0e7634b9fd9176866bb83395a757a52e77ee00dad25df447b352c2b2aa1d7c8b6009ac365cf0b66447a5fed118196a26368a84f8876826416c775277f64f61",
			Stake:   utils2.ToArt(3_250_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Arthera Validator 1
			ID:      6,
			Addr:    "0xe837B818796f3657A988f39CAFfE75203c38b205",
			Pubkey:  "0xc0043e927eedf3ef43d38e614dd55c8308f50d92d095eb1778fada8244810bdc5a5612277e0240e8f4d1bf69309e11b2cae3619ca2b74551838827aadfbfe9a60a8f",
			Stake:   utils2.ToArt(2_150_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Arthera Validator 2
			ID:      7,
			Addr:    "0xB5b5dfB5e3277D79B8E59389cecC74D46B2Bb607",
			Pubkey:  "0xc0043fbb08c7bbcac82e2722b678d0740c0f5b58ea465d54e7d5954d19eacb3bb18aa0a6e6f0f9db43c60487b25193bb28318a1010e96b024602f400c974fd5262d4",
			Stake:   utils2.ToArt(2_150_000),
			Balance: utils2.ToArt(10_000),
		},
		{
			// Arthera Validator 3
			ID:      8,
			Addr:    "0xFea412F092122461C61F7A1b849C5bbd88e08e5A",
			Pubkey:  "0xc00485f89b7f12d5050ef72c3d8d96af8eb0f1ff1ada3e913d505a30d739954398408d65d398385ee2f6582aac887720f3855148a871ee138b2295e8b4e022550cee",
			Stake:   utils2.ToArt(2_150_000),
			Balance: utils2.ToArt(10_000),
		},
	}
	MainnetAccounts = []GenesisAccount{
		{
			// Pre-Seed
			Addr:    "0x461292FD4Cc5598938c500065b59045cB6F441A8",
			Balance: utils2.ToArt(2_000_000),
		},
		{
			// Seed
			Addr:    "0xe2618bF31A16Bba5D296C118F8D75bD6C4a75dBc",
			Balance: utils2.ToArt(12_000_000),
		},
		{
			// Perpetual Security Fund
			Addr:    "0xe723a256D8912615A47Ef84F0830b05b22a7e9F5",
			Balance: utils2.ToArt(8_000_000),
		},
		{
			// Insurance Fund
			Addr:    "0x6A412042A05167CdB80d49f2cdfEfb6240C158bA",
			Balance: utils2.ToArt(8_000_000),
		},
		{
			// Team
			Addr:    "0xcABa6265997076a95FC81ba19E5Ba22D38F72e00",
			Balance: utils2.ToArt(16_000_000),
		},
		{
			// Advisors
			Addr:    "0xEEaBCe1f04F29dB5cf25275C3b7475E79fC7568d",
			Balance: utils2.ToArt(4_000_000),
		},
		{
			// Marketing
			Addr:    "0xCc688B6A11271f54e84534d12A51Db04fa677C5f",
			Balance: utils2.ToArt(36_000_000),
		},
		{
			// Growth Reserve
			Addr:    "0xd949738946DB18e65DEf26a2A3528A14a5417AdF",
			Balance: utils2.ToArt(32_000_000),
		},
		{
			// DEX & CEX Liquidity
			Addr:    "0x625d81DFc1a2e6A0Bb87e8eD4ed81d54235fBc23",
			Balance: utils2.ToArt(22_000_000),
		},
		{
			// Project Treasury
			Addr:    "0xe37CdbfFa2a302e608b4f10FC832a5314Ed1529f",
			Balance: utils2.ToArt(24_000_000),
		},
		{
			// Web2 Grants
			Addr:    "0x6d8d46c01b50931239f2371BF289D38F173085Fd",
			Balance: utils2.ToArt(8_000_000),
		},
		{
			// Liquidity Pools
			Addr:    "0xBaCF674C76C945B024ACfC7EeAa4a75528BcebF3",
			Balance: utils2.ToArt(30_000_000),
		},
		{
			// Genesis Validator Delegation
			Addr:    "0x9D986C0c1F931EB700161aaD57e9D799FCEA48b1",
			Balance: utils2.ToArt(7_220_000),
		},
	}
)
