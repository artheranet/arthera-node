package launcher

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/fake"
	"github.com/artheranet/arthera-node/params"
	utils2 "github.com/artheranet/arthera-node/utils"
	"github.com/artheranet/arthera-node/utils/memory"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Fantom-foundation/lachesis-base/abft"
	"github.com/Fantom-foundation/lachesis-base/utils/cachescale"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/naoina/toml"
	"gopkg.in/urfave/cli.v1"

	"github.com/artheranet/arthera-node/genesis/genesisstore"
	"github.com/artheranet/arthera-node/gossip"
	"github.com/artheranet/arthera-node/gossip/emitter"
	"github.com/artheranet/arthera-node/gossip/gasprice"
	"github.com/artheranet/arthera-node/internal/dbconfig"
	"github.com/artheranet/arthera-node/internal/evmcore"
	"github.com/artheranet/arthera-node/internal/vecmt"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, testFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}
	checkConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(checkConfig),
		Name:        "checkconfig",
		Usage:       "Checks configuration file",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, testFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The checkconfig checks configuration file.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}

	// DataDirFlag defines directory to store Lachesis state and user's wallets
	DataDirFlag = utils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: utils.DirectoryString(DefaultDataDir()),
	}

	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal caching",
		Value: DefaultCacheSize,
	}
	// GenesisFlag specifies network genesis configuration
	GenesisFlag = cli.StringFlag{
		Name:  "genesis",
		Usage: "'path to genesis file' - sets the network genesis configuration.",
	}
	ExperimentalGenesisFlag = cli.BoolFlag{
		Name:  "genesis.allowExperimental",
		Usage: "Allow to use experimental genesis file.",
	}

	RPCGlobalGasCapFlag = cli.Uint64Flag{
		Name:  "rpc.gascap",
		Usage: "Sets a cap on gas that can be used in art_call/estimateGas (0=infinite)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCGasCap,
	}
	RPCGlobalTxFeeCapFlag = cli.Float64Flag{
		Name:  "rpc.txfeecap",
		Usage: "Sets a cap on transaction fee (in ART) that can be sent via the RPC APIs (0 = no cap)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCTxFeeCap,
	}

	SyncModeFlag = cli.StringFlag{
		Name:  "syncmode",
		Usage: `Blockchain sync mode ("full" or "snap")`,
		Value: "full",
	}

	GCModeFlag = cli.StringFlag{
		Name:  "gcmode",
		Usage: `Blockchain garbage collection mode ("light", "full", "archive")`,
		Value: "archive",
	}

	ExitWhenAgeFlag = cli.DurationFlag{
		Name:  "exitwhensynced.age",
		Usage: "Exits after synchronisation reaches the required age",
	}
	ExitWhenEpochFlag = cli.Uint64Flag{
		Name:  "exitwhensynced.epoch",
		Usage: "Exits after synchronisation reaches the required epoch",
	}
	// TraceNodeFlag enables transaction tracing recording
	TraceNodeFlag = cli.BoolFlag{
		Name:  "tracenode",
		Usage: "If present, than this node records inner transaction traces",
	}
	TestnetFlag = cli.BoolFlag{
		Name:  "testnet",
		Usage: "Run a testnet node",
	}
)

type GenesisTemplate struct {
	Name   string
	Header genesis.Header
	Hashes genesis.Hashes
}

const (
	// DefaultCacheSize is calculated as memory consumption in a worst case scenario with default configuration
	// Average memory consumption might be 3-5 times lower than the maximum
	DefaultCacheSize  = 3600
	ConstantCacheSize = 400
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		return fmt.Errorf("field '%s' is not defined in %s", field, rt.String())
	},
}

type config struct {
	Node          node.Config
	Opera         gossip.Config
	Emitter       emitter.Config
	TxPool        evmcore.TxPoolConfig
	ArtheraStore  gossip.StoreConfig
	Lachesis      abft.Config
	LachesisStore abft.StoreConfig
	VectorClock   vecmt.IndexConfig
	DBs           dbconfig.DBsConfig
}

func (c *config) AppConfigs() dbconfig.Configs {
	return dbconfig.Configs{
		Opera:         c.Opera,
		OperaStore:    c.ArtheraStore,
		Lachesis:      c.Lachesis,
		LachesisStore: c.LachesisStore,
		VectorClock:   c.VectorClock,
		DBs:           c.DBs,
	}
}

func loadAllConfigs(file string, cfg *config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	if err != nil {
		return errors.New(fmt.Sprintf("TOML config file error: %v.\n"+
			"Use 'dumpconfig' command to get an example config file.\n"+
			"If node was recently upgraded and a previous network config file is used, then check updates for the config file.", err))
	}
	return err
}

func mayGetGenesisStore(ctx *cli.Context) *genesisstore.Store {
	switch {
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, num, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		return fake.FakeGenesisStore(num, utils2.ToArt(1_000_000_000), utils2.ToArt(5_000_000))
	case ctx.GlobalIsSet(TestnetFlag.Name):
		gen, _ := CreateGenesis("testnet")
		return gen
	case ctx.GlobalIsSet(GenesisFlag.Name):
		genesisPath := ctx.GlobalString(GenesisFlag.Name)

		f, err := os.Open(genesisPath)
		if err != nil {
			utils.Fatalf("Failed to open genesis file: %v", err)
		}
		genesisStore, genesisHashes, err := genesisstore.OpenGenesisStore(f)
		if err != nil {
			utils.Fatalf("Failed to read genesis file: %v", err)
		}

		// check if it's a trusted preset
		{
			g := genesisStore.Genesis()
			gHeader := genesis.Header{
				GenesisID:   g.GenesisID,
				NetworkID:   g.NetworkID,
				NetworkName: g.NetworkName,
			}
			for _, allowed := range AllowedArtheraGenesis {
				if allowed.Hashes.Equal(genesisHashes) && allowed.Header.Equal(gHeader) {
					log.Info("Genesis file is a known preset", "name", allowed.Name)
					goto notExperimental
				}
			}
			if ctx.GlobalBool(ExperimentalGenesisFlag.Name) {
				log.Warn("Genesis file doesn't refer to any trusted preset")
			} else {
				utils.Fatalf("Genesis file doesn't refer to any trusted preset. Enable experimental genesis with --genesis.allowExperimental")
			}
		notExperimental:
		}
		return genesisStore
	}
	return nil
}

func setBootnodes(ctx *cli.Context, urls []string, cfg *node.Config) {
	cfg.P2P.BootstrapNodesV5 = []*enode.Node{}
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				log.Error("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.P2P.BootstrapNodesV5 = append(cfg.P2P.BootstrapNodesV5, node)
		}
	}
	cfg.P2P.BootstrapNodes = cfg.P2P.BootstrapNodesV5
}

func setDataDir(ctx *cli.Context, cfg *node.Config) {
	defaultDataDir := DefaultDataDir()

	switch {
	case ctx.GlobalIsSet(DataDirFlag.Name):
		cfg.DataDir = ctx.GlobalString(DataDirFlag.Name)
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, num, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		cfg.DataDir = filepath.Join(defaultDataDir, fmt.Sprintf("fakenet-%d", num))
	}
}

func setGPO(ctx *cli.Context, cfg *gasprice.Config) {}

func setTxPool(ctx *cli.Context, cfg *evmcore.TxPoolConfig) {
	if ctx.GlobalIsSet(utils.TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(utils.TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsHexAddress(trimmed) {
				utils.Fatalf("Invalid account in --txpool.locals: %s", trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.HexToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(utils.TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(utils.TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(utils.TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(utils.TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(utils.TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(utils.TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(utils.TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(utils.TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(utils.TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(utils.TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(utils.TxPoolLifetimeFlag.Name)
	}
}

func gossipConfigWithFlags(ctx *cli.Context, src gossip.Config) (gossip.Config, error) {
	cfg := src

	setGPO(ctx, &cfg.GPO)

	if ctx.GlobalIsSet(RPCGlobalGasCapFlag.Name) {
		cfg.RPCGasCap = ctx.GlobalUint64(RPCGlobalGasCapFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(RPCGlobalTxFeeCapFlag.Name)
	}
	if ctx.GlobalIsSet(SyncModeFlag.Name) {
		if syncmode := ctx.GlobalString(SyncModeFlag.Name); syncmode != "full" && syncmode != "snap" {
			utils.Fatalf("--%s must be either 'full' or 'snap'", SyncModeFlag.Name)
		}
		cfg.AllowSnapsync = ctx.GlobalString(SyncModeFlag.Name) == "snap"
	}

	return cfg, nil
}

func gossipStoreConfigWithFlags(ctx *cli.Context, src gossip.StoreConfig) (gossip.StoreConfig, error) {
	cfg := src
	if ctx.GlobalIsSet(utils.GCModeFlag.Name) {
		if gcmode := ctx.GlobalString(utils.GCModeFlag.Name); gcmode != "light" && gcmode != "full" && gcmode != "archive" {
			utils.Fatalf("--%s must be 'light', 'full' or 'archive'", GCModeFlag.Name)
		}
		cfg.EVM.Cache.TrieDirtyDisabled = ctx.GlobalString(utils.GCModeFlag.Name) == "archive"
		cfg.EVM.Cache.GreedyGC = ctx.GlobalString(utils.GCModeFlag.Name) == "full"
	}
	return cfg, nil
}

func nodeConfigWithFlags(ctx *cli.Context, cfg node.Config) node.Config {
	utils.SetNodeConfig(ctx, &cfg)

	setDataDir(ctx, &cfg)
	return cfg
}

func cacheScaler(ctx *cli.Context) cachescale.Func {
	targetCache := ctx.GlobalInt(CacheFlag.Name)
	baseSize := DefaultCacheSize
	totalMemory := int(memory.TotalMemory() / opt.MiB)
	maxCache := totalMemory * 3 / 5
	if maxCache < baseSize {
		maxCache = baseSize
	}
	if !ctx.GlobalIsSet(CacheFlag.Name) {
		recommendedCache := totalMemory / 2
		if recommendedCache > baseSize {
			log.Warn(fmt.Sprintf("Please add '--%s %d' flag to allocate more cache for Arthera. Total memory is %d MB.", CacheFlag.Name, recommendedCache, totalMemory))
		}
		return cachescale.Identity
	}
	if targetCache < baseSize {
		log.Crit("Invalid flag", "flag", CacheFlag.Name, "err", fmt.Sprintf("minimum cache size is %d MB", baseSize))
	}
	if totalMemory != 0 && targetCache > maxCache {
		log.Warn(fmt.Sprintf("Requested cache size exceeds 60%% of availible memory. Reducing cache size to %d MB.", maxCache))
		targetCache = maxCache
	}

	return cachescale.Ratio{
		Base:   uint64(baseSize - ConstantCacheSize),
		Target: uint64(targetCache - ConstantCacheSize),
	}
}

func mayMakeAllConfigs(ctx *cli.Context) (*config, error) {
	// Defaults (low priority)
	cacheRatio := cacheScaler(ctx)
	cfg := config{
		Node:          defaultNodeConfig(),
		Opera:         gossip.DefaultConfig(cacheRatio),
		Emitter:       emitter.DefaultConfig(),
		TxPool:        evmcore.DefaultTxPoolConfig,
		ArtheraStore:  gossip.DefaultStoreConfig(cacheRatio),
		Lachesis:      abft.DefaultConfig(),
		LachesisStore: abft.DefaultStoreConfig(cacheRatio),
		VectorClock:   vecmt.DefaultConfig(cacheRatio),
	}

	if ctx.IsSet(utils.EWASMInterpreterFlag.Name) {
		params.DefaultVMConfig.EWASMInterpreter = ctx.String(utils.EWASMInterpreterFlag.Name)
		vm.InitEVMCEwasm(params.DefaultVMConfig.EWASMInterpreter)
	}

	if ctx.IsSet(utils.EVMInterpreterFlag.Name) {
		params.DefaultVMConfig.EVMInterpreter = ctx.String(utils.EVMInterpreterFlag.Name)
		vm.InitEVMCEVM(params.DefaultVMConfig.EVMInterpreter)
	}

	if ctx.GlobalIsSet(FakeNetFlag.Name) {
		_, num, _ := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		cfg.Emitter = emitter.FakeConfig(num)
		setBootnodes(ctx, []string{}, &cfg.Node)
	} else if ctx.GlobalIsSet(TestnetFlag.Name) {
		cfg.Emitter = emitter.TestnetConfig()
		cfg.Node.P2P.BootstrapNodes = asDefault
		cfg.Node.P2P.BootstrapNodesV5 = asDefault
	} else {
		// "asDefault" means set network defaults
		cfg.Node.P2P.BootstrapNodes = asDefault
		cfg.Node.P2P.BootstrapNodesV5 = asDefault
	}

	// Load config file (medium priority)
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadAllConfigs(file, &cfg); err != nil {
			return &cfg, err
		}
	}
	// apply default for DB config if it wasn't touched by config file
	dbDefault := dbconfig.DefaultDBsConfig(cacheRatio.U64, uint64(utils.MakeDatabaseHandles()))

	if len(cfg.DBs.Routing.Table) == 0 {
		cfg.DBs.Routing = dbDefault.Routing
	}
	if len(cfg.DBs.GenesisCache.Table) == 0 {
		cfg.DBs.GenesisCache = dbDefault.GenesisCache
	}
	if len(cfg.DBs.RuntimeCache.Table) == 0 {
		cfg.DBs.RuntimeCache = dbDefault.RuntimeCache
	}

	// Apply flags (high priority)
	var err error
	cfg.Opera, err = gossipConfigWithFlags(ctx, cfg.Opera)
	if err != nil {
		return nil, err
	}
	cfg.ArtheraStore, err = gossipStoreConfigWithFlags(ctx, cfg.ArtheraStore)
	if err != nil {
		return nil, err
	}
	cfg.Node = nodeConfigWithFlags(ctx, cfg.Node)

	err = setValidator(ctx, &cfg.Emitter)
	if err != nil {
		return nil, err
	}
	if cfg.Emitter.Validator.ID != 0 && len(cfg.Emitter.PrevEmittedEventFile.Path) == 0 {
		cfg.Emitter.PrevEmittedEventFile.Path = cfg.Node.ResolvePath(path.Join("emitter", fmt.Sprintf("last-%d", cfg.Emitter.Validator.ID)))
	}
	setTxPool(ctx, &cfg.TxPool)

	if err := cfg.Opera.Validate(); err != nil {
		return nil, err
	}

	if ctx.GlobalIsSet(TraceNodeFlag.Name) {
		cfg.ArtheraStore.TraceTransactions = true
	}

	return &cfg, nil
}

func makeAllConfigs(ctx *cli.Context) *config {
	cfg, err := mayMakeAllConfigs(ctx)
	if err != nil {
		utils.Fatalf("%v", err)
	}
	return cfg
}

func defaultNodeConfig() node.Config {
	cfg := NodeDefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = ethparams.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "art", "dag", "abft", "web3")
	cfg.WSModules = append(cfg.WSModules, "eth", "art", "dag", "abft", "web3")
	cfg.IPCPath = "arthera.ipc"
	cfg.DataDir = DefaultDataDir()
	return cfg
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)
	comment := ""

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	return nil
}

func checkConfig(ctx *cli.Context) error {
	_, err := mayMakeAllConfigs(ctx)
	return err
}
