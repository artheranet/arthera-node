package makegenesis

import (
	"bytes"
	"errors"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/lachesis"
	"github.com/artheranet/arthera-node/arthera/contracts/driver"
	"github.com/artheranet/arthera-node/arthera/contracts/driver/drivercall"
	"github.com/artheranet/arthera-node/arthera/contracts/driverauth"
	"github.com/artheranet/arthera-node/arthera/contracts/evmwriter"
	"github.com/artheranet/arthera-node/arthera/contracts/netinit"
	"github.com/artheranet/arthera-node/arthera/contracts/registry"
	"github.com/artheranet/arthera-node/arthera/contracts/staking"
	"github.com/artheranet/arthera-node/arthera/contracts/subscription"
	"github.com/artheranet/arthera-node/arthera/genesis/gpos"
	"github.com/artheranet/arthera-node/inter/drivertype"
	"io"
	"math/big"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/artheranet/arthera-node/arthera"
	"github.com/artheranet/arthera-node/arthera/genesis"
	"github.com/artheranet/arthera-node/arthera/genesisstore"
	"github.com/artheranet/arthera-node/evmcore"
	"github.com/artheranet/arthera-node/gossip/blockproc"
	"github.com/artheranet/arthera-node/gossip/blockproc/drivermodule"
	"github.com/artheranet/arthera-node/gossip/blockproc/eventmodule"
	"github.com/artheranet/arthera-node/gossip/blockproc/evmmodule"
	"github.com/artheranet/arthera-node/gossip/blockproc/sealmodule"
	"github.com/artheranet/arthera-node/gossip/evmstore"
	"github.com/artheranet/arthera-node/inter"
	"github.com/artheranet/arthera-node/inter/iblockproc"
	"github.com/artheranet/arthera-node/inter/ibr"
	"github.com/artheranet/arthera-node/inter/ier"
	"github.com/artheranet/arthera-node/utils/iodb"
)

type GenesisBuilder struct {
	dbs kvdb.DBProducer

	tmpEvmStore *evmstore.Store
	tmpStateDB  *state.StateDB

	totalSupply *big.Int

	blocks       []ibr.LlrIdxFullBlockRecord
	epochs       []ier.LlrIdxFullEpochRecord
	currentEpoch ier.LlrIdxFullEpochRecord
}

type BlockProc struct {
	SealerModule     blockproc.SealerModule
	TxListenerModule blockproc.TxListenerModule
	PreTxTransactor  blockproc.TxTransactor
	PostTxTransactor blockproc.TxTransactor
	EventsModule     blockproc.ConfirmedEventsModule
	EVMModule        blockproc.EVM
}

func DefaultBlockProc() BlockProc {
	return BlockProc{
		SealerModule:     sealmodule.New(),
		TxListenerModule: drivermodule.NewDriverTxListenerModule(),
		PreTxTransactor:  drivermodule.NewDriverTxPreTransactor(),
		PostTxTransactor: drivermodule.NewDriverTxTransactor(),
		EventsModule:     eventmodule.New(),
		EVMModule:        evmmodule.New(),
	}
}

func (b *GenesisBuilder) GetStateDB() *state.StateDB {
	if b.tmpStateDB == nil {
		tmpEvmStore := evmstore.NewStore(b.dbs, evmstore.LiteStoreConfig())
		b.tmpStateDB, _ = tmpEvmStore.StateDB(hash.Zero)
	}
	return b.tmpStateDB
}

func (b *GenesisBuilder) AddBalance(acc common.Address, balance *big.Int) {
	b.tmpStateDB.AddBalance(acc, balance)
	b.totalSupply.Add(b.totalSupply, balance)
}

func (b *GenesisBuilder) SetCode(acc common.Address, code []byte) {
	b.tmpStateDB.SetCode(acc, code)
}

func (b *GenesisBuilder) SetNonce(acc common.Address, nonce uint64) {
	b.tmpStateDB.SetNonce(acc, nonce)
}

func (b *GenesisBuilder) SetStorage(acc common.Address, key, val common.Hash) {
	b.tmpStateDB.SetState(acc, key, val)
}

func (b *GenesisBuilder) AddBlock(br ibr.LlrIdxFullBlockRecord) {
	b.blocks = append(b.blocks, br)
}

func (b *GenesisBuilder) AddEpoch(er ier.LlrIdxFullEpochRecord) {
	b.epochs = append(b.epochs, er)
}

func (b *GenesisBuilder) SetCurrentEpoch(er ier.LlrIdxFullEpochRecord) {
	b.currentEpoch = er
}

func (b *GenesisBuilder) GetCurrentEpoch() ier.LlrIdxFullEpochRecord {
	return b.currentEpoch
}

func (b *GenesisBuilder) TotalSupply() *big.Int {
	return b.totalSupply
}

func (b *GenesisBuilder) CurrentHash() hash.Hash {
	er := b.epochs[len(b.epochs)-1]
	return er.Hash()
}

func NewGenesisBuilder(dbs kvdb.DBProducer) *GenesisBuilder {
	tmpEvmStore := evmstore.NewStore(dbs, evmstore.LiteStoreConfig())
	statedb, _ := tmpEvmStore.StateDB(hash.Zero)
	return &GenesisBuilder{
		dbs:         dbs,
		tmpEvmStore: tmpEvmStore,
		tmpStateDB:  statedb,
		totalSupply: new(big.Int),
	}
}

type dummyHeaderReturner struct {
}

func (d dummyHeaderReturner) GetHeader(common.Hash, uint64) *evmcore.EvmHeader {
	return &evmcore.EvmHeader{}
}

func (b *GenesisBuilder) ExecuteGenesisTxs(blockProc BlockProc, genesisTxs types.Transactions) error {
	bs, es := b.currentEpoch.BlockState.Copy(), b.currentEpoch.EpochState.Copy()

	blockCtx := iblockproc.BlockCtx{
		Idx:     bs.LastBlock.Idx + 1,
		Time:    bs.LastBlock.Time + 1,
		Atropos: hash.Event{},
	}

	sealer := blockProc.SealerModule.Start(blockCtx, bs, es)
	sealing := true
	txListener := blockProc.TxListenerModule.Start(blockCtx, bs, es, b.tmpStateDB)
	evmProcessor := blockProc.EVMModule.Start(blockCtx, b.tmpStateDB, dummyHeaderReturner{}, func(l *types.Log) {
		txListener.OnNewLog(l)
	}, es.Rules, es.Rules.EvmChainConfig([]opera.UpgradeHeight{
		{
			Upgrades: es.Rules.Upgrades,
			Height:   0,
		},
	}))

	// Execute genesis transactions
	evmProcessor.Execute(genesisTxs)
	bs = txListener.Finalize()

	// Execute pre-internal transactions
	preInternalTxs := blockProc.PreTxTransactor.PopInternalTxs(blockCtx, bs, es, sealing, b.tmpStateDB)
	evmProcessor.Execute(preInternalTxs)
	bs = txListener.Finalize()

	// Seal epoch if requested
	if sealing {
		sealer.Update(bs, es)
		bs, es = sealer.SealEpoch()
		txListener.Update(bs, es)
	}

	// Execute post-internal transactions
	internalTxs := blockProc.PostTxTransactor.PopInternalTxs(blockCtx, bs, es, sealing, b.tmpStateDB)
	evmProcessor.Execute(internalTxs)

	evmBlock, skippedTxs, receipts := evmProcessor.Finalize()
	for _, r := range receipts {
		if r.Status == 0 {
			return errors.New("genesis transaction reverted")
		}
	}
	if len(skippedTxs) != 0 {
		return errors.New("genesis transaction is skipped")
	}
	bs = txListener.Finalize()
	bs.FinalizedStateRoot = hash.Hash(evmBlock.Root)

	bs.LastBlock = blockCtx

	prettyHash := func(root hash.Hash) hash.Event {
		e := inter.MutableEventPayload{}
		// for nice-looking ID
		e.SetEpoch(es.Epoch)
		e.SetLamport(1)
		// actual data hashed
		e.SetExtra(root[:])

		return e.Build().ID()
	}
	receiptsStorage := make([]*types.ReceiptForStorage, len(receipts))
	for i, r := range receipts {
		receiptsStorage[i] = (*types.ReceiptForStorage)(r)
	}
	// add block
	b.blocks = append(b.blocks, ibr.LlrIdxFullBlockRecord{
		LlrFullBlockRecord: ibr.LlrFullBlockRecord{
			Atropos:  prettyHash(bs.FinalizedStateRoot),
			Root:     bs.FinalizedStateRoot,
			Txs:      evmBlock.Transactions,
			Receipts: receiptsStorage,
			Time:     blockCtx.Time,
			GasUsed:  evmBlock.GasUsed,
		},
		Idx: blockCtx.Idx,
	})
	// add epoch
	b.currentEpoch = ier.LlrIdxFullEpochRecord{
		LlrFullEpochRecord: ier.LlrFullEpochRecord{
			BlockState: bs,
			EpochState: es,
		},
		Idx: es.Epoch,
	}
	b.epochs = append(b.epochs, b.currentEpoch)

	return b.tmpEvmStore.Commit(bs.LastBlock.Idx, bs.FinalizedStateRoot, true)
}

type memFile struct {
	*bytes.Buffer
}

func (f *memFile) Close() error {
	*f = memFile{}
	return nil
}

func (b *GenesisBuilder) Build(head genesis.Header) *genesisstore.Store {
	return genesisstore.NewStore(func(name string) (io.Reader, error) {
		buf := &memFile{bytes.NewBuffer(nil)}
		if name == genesisstore.BlocksSection {
			for i := len(b.blocks) - 1; i >= 0; i-- {
				_ = rlp.Encode(buf, b.blocks[i])
			}
			return buf, nil
		}
		if name == genesisstore.EpochsSection {
			for i := len(b.epochs) - 1; i >= 0; i-- {
				_ = rlp.Encode(buf, b.epochs[i])
			}
			return buf, nil
		}
		if name == genesisstore.EvmSection {
			it := b.tmpEvmStore.EvmDb.NewIterator(nil, nil)
			defer it.Release()
			_ = iodb.Write(buf, it)
		}
		if buf.Len() == 0 {
			return nil, errors.New("not found")
		}
		return buf, nil
	}, head, func() error {
		*b = GenesisBuilder{}
		return nil
	})
}

func (b *GenesisBuilder) DeployBaseContracts() {
	// deploy essential contracts
	// pre deploy NetworkInitializer
	b.SetCode(netinit.ContractAddress, netinit.GetContractBin())
	// pre deploy NodeDriver
	b.SetCode(driver.ContractAddress, driver.GetContractBin())
	// pre deploy NodeDriverAuth
	b.SetCode(driverauth.ContractAddress, driverauth.GetContractBin())
	// pre deploy Staking
	b.SetCode(staking.ContractAddress, staking.GetContractBin())
	b.SetCode(staking.ValidatorInfoContractAddress, staking.GetValidatorInfoContractBin())
	// pre deploy registry
	b.SetCode(registry.ContractAddress, registry.GetContractBin())
	// pre deploy subscription
	b.SetCode(subscription.ContractAddress, subscription.GetContractBin())
	// set non-zero code for pre-compiled contracts
	b.SetCode(evmwriter.ContractAddress, []byte{0})
}

func (b *GenesisBuilder) InitializeEpoch(block idx.Block, epoch idx.Epoch, rules opera.Rules, timestamp inter.Timestamp) {
	b.SetCurrentEpoch(ier.LlrIdxFullEpochRecord{
		LlrFullEpochRecord: ier.LlrFullEpochRecord{
			BlockState: iblockproc.BlockState{
				LastBlock: iblockproc.BlockCtx{
					Idx:     block - 1,
					Time:    timestamp,
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
				Epoch:             epoch - 1,
				EpochStart:        timestamp,
				PrevEpochStart:    timestamp - 1,
				EpochStateRoot:    hash.Zero,
				Validators:        pos.NewBuilder().Build(),
				ValidatorStates:   make([]iblockproc.ValidatorEpochState, 0),
				ValidatorProfiles: make(map[idx.ValidatorID]drivertype.Validator),
				Rules:             rules,
			},
		},
		Idx: epoch - 1,
	})
}

func (b *GenesisBuilder) GetGenesisTxs(sealedEpoch idx.Epoch, validators gpos.Validators, totalSupply *big.Int, delegations []drivercall.Delegation, driverOwner common.Address) types.Transactions {
	buildTx := txBuilder()
	internalTxs := make(types.Transactions, 0, 15)
	// initialization
	calldata := netinit.InitializeAll(
		sealedEpoch,
		totalSupply,
		staking.ContractAddress,
		driverauth.ContractAddress,
		driver.ContractAddress,
		evmwriter.ContractAddress,
		staking.ValidatorInfoContractAddress,
		subscription.ContractAddress,
		driverOwner,
	)
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
