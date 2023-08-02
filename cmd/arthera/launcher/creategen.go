package launcher

import (
	"fmt"
	"github.com/artheranet/arthera-node/contracts/driver"
	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/builder"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/artheranet/arthera-node/internal/inter/ibr"
	"github.com/artheranet/arthera-node/internal/inter/ier"
	"github.com/artheranet/arthera-node/internal/inter/validatorpk"
	"github.com/artheranet/arthera-node/params"
	"github.com/artheranet/arthera-node/utils/iodb"
	"github.com/artheranet/lachesis/hash"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/artheranet/lachesis/kvdb"
	"github.com/artheranet/lachesis/kvdb/memorydb"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"gopkg.in/urfave/cli.v1"
	"io"
	"math/big"
	"os"
	"time"
)

var genesisTypeFlag = cli.StringFlag{
	Name:  "genesis.type",
	Usage: "Type of genesis to generate: mainnet, testnet, devnet",
	Value: "devnet",
}

var (
	createGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(createGenesisCmd),
		Name:      "creategen",
		Usage:     "Create genesis",
		ArgsUsage: "<file>",
		Category:  "MISCELLANEOUS COMMANDS",
		Flags:     []cli.Flag{genesisTypeFlag},
	}
)

var (
	GenesisTime = inter.FromUnix(time.Now().Unix())
)

func createGenesisCmd(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		utils.Fatalf("error: need destination file as argument.")
	}

	genesisType := "devnet"
	if ctx.GlobalIsSet(genesisTypeFlag.Name) {
		genesisType = ctx.GlobalString(genesisTypeFlag.Name)
	}

	if genesisType != "mainnet" && genesisType != "testnet" && genesisType != "devnet" {
		utils.Fatalf("Unknown genesis type %s. Supported values are: devnet, testnet, mainnet", genesisType)
	}

	fileName := ctx.Args().Get(0)

	fmt.Println("Creating " + genesisType + " genesis")
	genesisStore, currentHash := CreateGenesis(genesisType)
	err := WriteGenesisStore(fileName, genesisStore, currentHash)
	if err != nil {
		return err
	}

	return nil
}

func CreateGenesis(genesisType string) (*genesisstore.Store, hash.Hash) {
	genesisBuilder := builder.NewGenesisBuilder(memorydb.NewProducer(""))

	validators := make(genesis.Validators, 0, 3)
	delegations := make([]driver.Delegation, 0, 3)

	var initialValidators = params.DevnetValidators
	var initialAccounts = params.DevnetAccounts
	if genesisType == "mainnet" {
		initialValidators = params.MainnetValidators
		initialAccounts = params.MainnetAccounts
	} else if genesisType == "testnet" {
		initialValidators = params.TestnetValidators
		initialAccounts = params.TestnetAccounts
	}

	// add initial validators, premine and lock their stake to get maximum rewards
	for _, v := range initialValidators {
		validators, delegations = AddValidator(
			v,
			validators, delegations, genesisBuilder,
		)
	}

	// premine to genesis accounts
	for _, a := range initialAccounts {
		genesisBuilder.AddBalance(
			common.HexToAddress(a.Addr),
			a.Balance,
		)
	}

	genesisBuilder.DeployBaseContracts()

	rules := params.DevNetRules()
	if genesisType == "mainnet" {
		rules = params.MainNetRules()
	} else if genesisType == "testnet" {
		rules = params.TestNetRules()
	}

	genesisBuilder.InitializeEpoch(1, 2, rules, GenesisTime)

	owner := validators[0].Address
	blockProc := builder.DefaultBlockProc()
	genesisTxs := genesisBuilder.GetGenesisTxs(0, validators, genesisBuilder.TotalSupply(), delegations, owner)
	err := genesisBuilder.ExecuteGenesisTxs(blockProc, genesisTxs)
	if err != nil {
		panic(err)
	}

	return genesisBuilder.Build(genesis.Header{
		GenesisID:   genesisBuilder.CurrentHash(),
		NetworkID:   rules.NetworkID,
		NetworkName: rules.Name,
	}), genesisBuilder.CurrentHash()
}

func AddValidator(
	v params.GenesisValidator,
	validators genesis.Validators,
	delegations []driver.Delegation,
	builder *builder.GenesisBuilder,
) (genesis.Validators, []driver.Delegation) {
	validatorId := idx.ValidatorID(v.ID)
	pk, _ := validatorpk.FromString(v.Pubkey)
	ecdsaPubkey, _ := crypto.UnmarshalPubkey(pk.Raw)
	addr := crypto.PubkeyToAddress(*ecdsaPubkey)

	validator := genesis.Validator{
		ID:      validatorId,
		Address: addr,
		PubKey: validatorpk.PubKey{
			Raw:  pk.Raw,
			Type: validatorpk.Types.Secp256k1,
		},
		CreationTime:     GenesisTime,
		CreationEpoch:    0,
		DeactivatedTime:  0,
		DeactivatedEpoch: 0,
		Status:           0,
	}
	builder.AddBalance(validator.Address, v.Balance)
	validators = append(validators, validator)

	delegations = append(delegations, driver.Delegation{
		Address:            validator.Address,
		ValidatorID:        validator.ID,
		Stake:              v.Stake,
		LockedStake:        new(big.Int),
		LockupFromEpoch:    0,
		LockupEndTime:      0,
		LockupDuration:     0,
		EarlyUnlockPenalty: new(big.Int),
		Rewards:            new(big.Int),
	})

	return validators, delegations
}

func WriteGenesisStore(fn string, gs *genesisstore.Store, genesisHash hash.Hash) error {
	var plain io.WriteSeeker

	log.Info("GenesisID ", "hash", genesisHash.String())

	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()
	plain = fh

	writer := newUnitWriter(plain)
	err = writer.Start(gs.Header(), genesisstore.EpochsSection, "/tmp/gentmp")
	if err != nil {
		return err
	}

	gs.Epochs().ForEach(func(epochRecord ier.LlrIdxFullEpochRecord) bool {
		b, _ := rlp.EncodeToBytes(epochRecord)
		_, err := writer.Write(b)
		if err != nil {
			panic(err)
		}
		return true
	})

	var epochsHash hash.Hash
	epochsHash, err = writer.Flush()
	if err != nil {
		return err
	}
	log.Info("Exported epochs", "hash", epochsHash.String())

	writer = newUnitWriter(plain)
	err = writer.Start(gs.Header(), genesisstore.BlocksSection, "/tmp/gentmp")
	if err != nil {
		return err
	}
	gs.Blocks().ForEach(func(blockRecord ibr.LlrIdxFullBlockRecord) bool {
		b, _ := rlp.EncodeToBytes(blockRecord)
		_, err := writer.Write(b)
		if err != nil {
			panic(err)
		}
		return true
	})

	var blocksHash hash.Hash
	blocksHash, err = writer.Flush()
	if err != nil {
		return err
	}
	log.Info("Exported blocks", "hash", blocksHash.String())

	writer = newUnitWriter(plain)
	err = writer.Start(gs.Header(), genesisstore.EvmSection, "/tmp/gentmp")
	if err != nil {
		return err
	}

	gs.RawEvmItems().(genesisstore.RawEvmItems).Iterator(func(it kvdb.Iterator) bool {
		defer it.Release()
		err = iodb.Write(writer, it)
		if err != nil {
			panic(err)
		}
		return true
	})

	var evmHash hash.Hash
	evmHash, err = writer.Flush()
	if err != nil {
		return err
	}
	log.Info("Exported EVM data", "hash", evmHash.String())

	fmt.Printf("Exported genesis to file %s\n", fn)
	return nil
}
