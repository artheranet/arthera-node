package fake

import (
	"crypto/ecdsa"
	"github.com/artheranet/arthera-node/contracts/driver"
	"github.com/artheranet/arthera-node/genesis/builder"
	"github.com/artheranet/arthera-node/params"
	"math/big"
	"time"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
	"github.com/artheranet/arthera-node/inter"
	"github.com/artheranet/arthera-node/inter/validatorpk"
	"github.com/artheranet/arthera-node/internal/evmcore"
)

var (
	FakeGenesisTime = inter.Timestamp(1608600000 * time.Second)
)

// FakeKey gets n-th fake private key.
func FakeKey(n idx.ValidatorID) *ecdsa.PrivateKey {
	return evmcore.FakeKey(int(n))
}

func FakeGenesisStore(num idx.Validator, balance, stake *big.Int) *genesisstore.Store {
	return FakeGenesisStoreWithRulesAndStart(num, balance, stake, params.FakeNetRules(), 2, 1)
}

func FakeGenesisStoreWithRulesAndStart(num idx.Validator, balance, stake *big.Int, rules params.ProtocolRules, epoch idx.Epoch, block idx.Block) *genesisstore.Store {
	genesisBuilder := builder.NewGenesisBuilder(memorydb.NewProducer(""))

	validators := GetFakeValidators(num)

	// add balances to validators
	var delegations []driver.Delegation
	for _, val := range validators {
		genesisBuilder.AddBalance(val.Address, balance)
		delegations = append(delegations, driver.Delegation{
			Address:            val.Address,
			ValidatorID:        val.ID,
			Stake:              stake,
			LockedStake:        new(big.Int),
			LockupFromEpoch:    0,
			LockupEndTime:      0,
			LockupDuration:     0,
			EarlyUnlockPenalty: new(big.Int),
			Rewards:            new(big.Int),
		})
	}

	genesisBuilder.DeployBaseContracts()
	genesisBuilder.InitializeEpoch(block, epoch, rules, FakeGenesisTime)

	var owner common.Address
	if num != 0 {
		owner = validators[0].Address
	}

	blockProc := builder.DefaultBlockProc()
	genesisTxs := genesisBuilder.GetGenesisTxs(epoch-2, validators, genesisBuilder.TotalSupply(), delegations, owner)
	err := genesisBuilder.ExecuteGenesisTxs(blockProc, genesisTxs)
	if err != nil {
		panic(err)
	}

	return genesisBuilder.Build(genesis.Header{
		GenesisID:   genesisBuilder.CurrentHash(),
		NetworkID:   rules.NetworkID,
		NetworkName: rules.Name,
	})
}

func GetFakeValidators(num idx.Validator) genesis.Validators {
	validators := make(genesis.Validators, 0, num)

	for i := idx.ValidatorID(1); i <= idx.ValidatorID(num); i++ {
		key := FakeKey(i)
		addr := crypto.PubkeyToAddress(key.PublicKey)
		pubkeyraw := crypto.FromECDSAPub(&key.PublicKey)
		validators = append(validators, genesis.Validator{
			ID:      i,
			Address: addr,
			PubKey: validatorpk.PubKey{
				Raw:  pubkeyraw,
				Type: validatorpk.Types.Secp256k1,
			},
			CreationTime:     FakeGenesisTime,
			CreationEpoch:    0,
			DeactivatedTime:  0,
			DeactivatedEpoch: 0,
			Status:           0,
		})
	}

	return validators
}
