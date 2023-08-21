package fake

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/artheranet/arthera-node/contracts/driver"
	"github.com/artheranet/arthera-node/genesis/builder"
	"github.com/artheranet/arthera-node/params"
	"github.com/artheranet/arthera-node/utils"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"

	"github.com/artheranet/lachesis/inter/idx"
	"github.com/artheranet/lachesis/kvdb/memorydb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
	"github.com/artheranet/arthera-node/internal/evmcore"
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/artheranet/arthera-node/internal/inter/validatorpk"
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
	if num > 10 {
		log.Error("Too many validators, maximum is 10", "num", num)
	}

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

	log.Info("------------ Fakenet Info -----------")
	log.Info("---> Validators:")
	for _, val := range validators {
		log.Info("Validator", "ID", val.ID, "Address", val.Address.String(), "Public Key",
			val.PubKey.String(), "Private key", hex.EncodeToString(val.PrivKey.D.Bytes()), "Balance", utils.WeiToArt(balance), "Stake", utils.WeiToArt(stake))
	}
	for i := 100; i < 110; i++ {
		key := evmcore.FakeKey(i)
		addr := crypto.PubkeyToAddress(key.PublicKey)
		genesisBuilder.AddBalance(addr, balance)
		log.Info("Account", "Address", addr.String(), "Private key", hex.EncodeToString(key.D.Bytes()), "Balance", balance)
	}
	log.Info("------------ Fakenet Info -----------")

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
			PrivKey:          key,
			CreationTime:     FakeGenesisTime,
			CreationEpoch:    0,
			DeactivatedTime:  0,
			DeactivatedEpoch: 0,
			Status:           0,
		})
	}

	return validators
}
