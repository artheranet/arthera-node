package launcher

import (
	"fmt"
	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
	"github.com/Fantom-foundation/lachesis-base/lachesis"
	"github.com/artheranet/arthera-node/integration/makegenesis"
	"github.com/artheranet/arthera-node/inter"
	"github.com/artheranet/arthera-node/inter/drivertype"
	"github.com/artheranet/arthera-node/inter/iblockproc"
	"github.com/artheranet/arthera-node/inter/ibr"
	"github.com/artheranet/arthera-node/inter/ier"
	"github.com/artheranet/arthera-node/inter/validatorpk"
	"github.com/artheranet/arthera-node/opera"
	"github.com/artheranet/arthera-node/opera/contracts/driver"
	"github.com/artheranet/arthera-node/opera/contracts/driver/drivercall"
	"github.com/artheranet/arthera-node/opera/contracts/driverauth"
	"github.com/artheranet/arthera-node/opera/contracts/evmwriter"
	"github.com/artheranet/arthera-node/opera/contracts/netinit"
	netinitcall "github.com/artheranet/arthera-node/opera/contracts/netinit/netinitcalls"
	"github.com/artheranet/arthera-node/opera/contracts/staking"
	"github.com/artheranet/arthera-node/opera/genesis"
	"github.com/artheranet/arthera-node/opera/genesis/gpos"
	"github.com/artheranet/arthera-node/opera/genesisstore"
	utils2 "github.com/artheranet/arthera-node/utils"
	"github.com/artheranet/arthera-node/utils/iodb"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"gopkg.in/urfave/cli.v1"
	"io"
	"math/big"
	"os"
)

var genesisTypeFlag = cli.StringFlag{
	Name:  "genesis.type",
	Usage: "Type of genesis to generate: mainnet, testnet",
	Value: "testnet",
}

var (
	createGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(createGenesisCmd),
		Name:      "creategen",
		Usage:     "Create genesis",
		ArgsUsage: "",
		Category:  "MISCELLANEOUS COMMANDS",
	}
)

type GenesisValidator struct {
	addr    string
	pubkey  string
	stake   *big.Int
	balance *big.Int
}

type GenesisAccount struct {
	addr    string
	balance *big.Int
}

var (
	TestnetValidators = []GenesisValidator{
		{
			addr:    "0x7a97E50436a074ADDB9A51D50Fbd35ADAFE88442",
			pubkey:  "0xc0041d7405a8bc7dabf1e397e6689ff09482466aea9d3a716bf1dd4fd971c22d035d8d939c88764136a3213106282887f9005b5addf23af781302a0119400706996e",
			stake:   utils2.ToArt(1_000_000), // min stake StakerConstants.sol -> minSelfStake
			balance: utils2.ToArt(0),
		},
		{
			addr:    "0xfE8301b91A8Eb4734ed954f8E2FB84c2F72Cef8a",
			pubkey:  "0xc004a61ec5eb3cf8d6b399ff56682b95277337b601fb31e1a254dd451101b8aafb0218d428fc814faee132aabcc17b3dd39fa35dfce2d5ce29d6bd05615bbd571016",
			stake:   utils2.ToArt(1_000_000), // min stake StakerConstants.sol -> minSelfStake
			balance: utils2.ToArt(0),
		},
		{
			addr:    "0xF51e935061731a129765ff63b3Af0Adb5e4486aC",
			pubkey:  "0xc004c39c38dc49cc4c9b64ea9d817545e713635f808d692f2f500ad801e002c50987e15cf4d9419731adf4cd83edf2207a806685cb2b75c3027d2dcdd78ec126f430",
			stake:   utils2.ToArt(1_000_000),
			balance: utils2.ToArt(0), // min stake StakerConstants.sol -> minSelfStake
		},
	}

	TestnetAccounts = []GenesisAccount{
		{
			addr:    "0x40bd65cfc4D95844704F4b2a2c46a60f6d6CE766",
			balance: utils2.ToArt(10_000_000),
		},
		{
			addr:    "0x35E58946b74fDbD9032aed876FC58629A6e65E79",
			balance: utils2.ToArt(10_000_000),
		},
		{
			addr:    "0x846032c611697818a31cC090D436664b263C6E54",
			balance: utils2.ToArt(10_000_000),
		},
	}

	MainnetValidators = []GenesisValidator{}
	MainnetAccounts   = []GenesisAccount{}

	GenesisTime = inter.FromUnix(1677067996)
)

func createGenesisCmd(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	genesisType := "testnet"
	if ctx.GlobalIsSet(genesisTypeFlag.Name) {
		genesisType = ctx.GlobalString(genesisTypeFlag.Name)
	}

	fileName := ctx.Args().First()

	fmt.Println("Creating " + genesisType + " genesis")
	genesisStore, currentHash := CreateGenesis(genesisType)
	err := WriteGenesisStore(fileName, genesisStore, currentHash)
	if err != nil {
		return err
	}

	return nil
}

func CreateGenesis(genesisType string) (*genesisstore.Store, hash.Hash) {
	builder := makegenesis.NewGenesisBuilder(memorydb.NewProducer(""))

	validators := make(gpos.Validators, 0, 3)
	delegations := make([]drivercall.Delegation, 0, 3)

	var initialValidators = TestnetValidators
	var initialAccounts = TestnetAccounts
	if genesisType == "mainnet" {
		initialValidators = MainnetValidators
		initialAccounts = MainnetAccounts
	}

	// add initial validators, premine and lock their stake to get maximum rewards
	for i, v := range initialValidators {
		validators, delegations = AddValidator(
			uint8(i+1),
			v,
			validators, delegations, builder,
		)
	}

	// premine to genesis accounts
	for _, a := range initialAccounts {
		builder.AddBalance(
			common.HexToAddress(a.addr),
			a.balance,
		)
	}

	// deploy essential contracts
	// pre deploy NetworkInitializer
	builder.SetCode(netinit.ContractAddress, netinit.GetContractBin())
	// pre deploy NodeDriver
	builder.SetCode(driver.ContractAddress, driver.GetContractBin())
	// pre deploy NodeDriverAuth
	builder.SetCode(driverauth.ContractAddress, driverauth.GetContractBin())
	// pre deploy Staking
	builder.SetCode(staking.ContractAddress, staking.GetContractBin())
	builder.SetCode(staking.StakerInfoContractAddress, staking.GetStakerInfoContractBin())
	// set non-zero code for pre-compiled contracts
	builder.SetCode(evmwriter.ContractAddress, []byte{0})

	rules := opera.TestNetRules()
	if genesisType == "mainnet" {
		rules = opera.MainNetRules()
	}

	builder.SetCurrentEpoch(ier.LlrIdxFullEpochRecord{
		LlrFullEpochRecord: ier.LlrFullEpochRecord{
			BlockState: iblockproc.BlockState{
				LastBlock: iblockproc.BlockCtx{
					Idx:     0,
					Time:    GenesisTime,
					Atropos: hash.Event{},
				},
				FinalizedStateRoot:    hash.Hash{},
				EpochGas:              0,
				EpochCheaters:         lachesis.Cheaters{},
				CheatersWritten:       0,
				ValidatorStates:       make([]iblockproc.ValidatorBlockState, 0),
				NextValidatorProfiles: make(map[idx.ValidatorID]drivertype.Validator),
				DirtyRules:            nil,
				AdvanceEpochs:         0,
			},
			EpochState: iblockproc.EpochState{
				Epoch:             1,
				EpochStart:        GenesisTime,
				PrevEpochStart:    GenesisTime - 1,
				EpochStateRoot:    hash.Zero,
				Validators:        pos.NewBuilder().Build(),
				ValidatorStates:   make([]iblockproc.ValidatorEpochState, 0),
				ValidatorProfiles: make(map[idx.ValidatorID]drivertype.Validator),
				Rules:             rules,
			},
		},
		Idx: 1,
	})

	owner := validators[0].Address
	blockProc := makegenesis.DefaultBlockProc()
	genesisTxs := GetGenesisTxs(0, validators, builder.TotalSupply(), delegations, owner)
	err := builder.ExecuteGenesisTxs(blockProc, genesisTxs)
	if err != nil {
		panic(err)
	}

	return builder.Build(genesis.Header{
		GenesisID:   builder.CurrentHash(),
		NetworkID:   rules.NetworkID,
		NetworkName: rules.Name,
	}), builder.CurrentHash()
}

func AddValidator(
	id uint8,
	v GenesisValidator,
	validators gpos.Validators,
	delegations []drivercall.Delegation,
	builder *makegenesis.GenesisBuilder,
) (gpos.Validators, []drivercall.Delegation) {
	validatorId := idx.ValidatorID(id)
	pk, _ := validatorpk.FromString(v.pubkey)
	ecdsaPubkey, _ := crypto.UnmarshalPubkey(pk.Raw)
	addr := crypto.PubkeyToAddress(*ecdsaPubkey)

	validator := gpos.Validator{
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
	builder.AddBalance(validator.Address, v.balance)
	validators = append(validators, validator)

	delegations = append(delegations, drivercall.Delegation{
		Address:            validator.Address,
		ValidatorID:        validator.ID,
		Stake:              v.stake,
		LockedStake:        new(big.Int),
		LockupFromEpoch:    0,
		LockupEndTime:      0,
		LockupDuration:     0,
		EarlyUnlockPenalty: new(big.Int),
		Rewards:            new(big.Int),
	})

	return validators, delegations
}

func GetGenesisTxs(sealedEpoch idx.Epoch, validators gpos.Validators, totalSupply *big.Int, delegations []drivercall.Delegation, driverOwner common.Address) types.Transactions {
	buildTx := txBuilder()
	internalTxs := make(types.Transactions, 0, 15)
	// initialization
	calldata := netinitcall.InitializeAll(sealedEpoch, totalSupply, staking.ContractAddress, driverauth.ContractAddress, driver.ContractAddress, evmwriter.ContractAddress, driverOwner)
	internalTxs = append(internalTxs, buildTx(calldata, netinit.ContractAddress))
	// push genesis validators
	for _, v := range validators {
		calldata := drivercall.SetGenesisValidator(v)
		internalTxs = append(internalTxs, buildTx(calldata, driver.ContractAddress))
	}
	// push genesis delegations
	for _, delegation := range delegations {
		calldata := drivercall.SetGenesisDelegation(delegation)
		internalTxs = append(internalTxs, buildTx(calldata, driver.ContractAddress))
	}
	return internalTxs
}

func txBuilder() func(calldata []byte, addr common.Address) *types.Transaction {
	nonce := uint64(0)
	return func(calldata []byte, addr common.Address) *types.Transaction {
		tx := types.NewTransaction(nonce, addr, common.Big0, 1e10, common.Big0, calldata)
		nonce++
		return tx
	}
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
