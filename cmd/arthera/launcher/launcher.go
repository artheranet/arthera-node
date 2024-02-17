package launcher

import (
	"fmt"
	"github.com/artheranet/arthera-node/gossip/emitter"
	"github.com/artheranet/arthera-node/params"
	version2 "github.com/artheranet/arthera-node/version"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/scwallet"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/crypto"
	evmetrics "github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/p2p/discover/discfilter"
	"gopkg.in/urfave/cli.v1"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/artheranet/arthera-node/cmd/arthera/launcher/flags"
	"github.com/artheranet/arthera-node/cmd/arthera/launcher/metrics"
	"github.com/artheranet/arthera-node/cmd/arthera/launcher/tracing"
	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
	"github.com/artheranet/arthera-node/gossip"
	"github.com/artheranet/arthera-node/internal/dbconfig"
	"github.com/artheranet/arthera-node/internal/evmcore"
	"github.com/artheranet/arthera-node/internal/valkeystore"
	"github.com/artheranet/arthera-node/utils/debug"
	"github.com/artheranet/arthera-node/utils/errlock"
	_ "github.com/artheranet/arthera-node/version"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console/prompt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
)

const (
	// clientIdentifier to advertise over the network.
	clientIdentifier = "arthera-node"
)

var (
	// The app that holds all commands and flags.
	app = flags.NewApp("arthera-node command line interface")

	nodeFlags        []cli.Flag
	testFlags        []cli.Flag
	gpoFlags         []cli.Flag
	accountFlags     []cli.Flag
	performanceFlags []cli.Flag
	networkingFlags  []cli.Flag
	txpoolFlags      []cli.Flag
	artheraFlags     []cli.Flag
	legacyRpcFlags   []cli.Flag
	rpcFlags         []cli.Flag
	metricsFlags     []cli.Flag
	vmFlags          []cli.Flag
)

func initFlags() {
	// Flags for testing purpose.
	testFlags = []cli.Flag{
		FakeNetFlag,
	}

	// Flags that configure the node.
	gpoFlags = []cli.Flag{}
	accountFlags = []cli.Flag{
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.ExternalSignerFlag,
		utils.InsecureUnlockAllowedFlag,
	}
	performanceFlags = []cli.Flag{
		CacheFlag,
	}
	networkingFlags = []cli.Flag{
		utils.BootnodesFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.IPrestrictFlag,
		utils.PrivateNodeFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
	}
	txpoolFlags = []cli.Flag{
		utils.TxPoolLocalsFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
	}
	artheraFlags = []cli.Flag{
		GenesisFlag,
		ExperimentalGenesisFlag,
		utils.IdentityFlag,
		DataDirFlag,
		utils.MinFreeDiskSpaceFlag,
		utils.KeyStoreDirFlag,
		utils.USBFlag,
		utils.SmartCardDaemonPathFlag,
		ExitWhenAgeFlag,
		ExitWhenEpochFlag,
		utils.LightKDFFlag,
		configFileFlag,
		validatorIDFlag,
		validatorPubkeyFlag,
		validatorPasswordFlag,
		SyncModeFlag,
		GCModeFlag,
		genesisTypeFlag,
		TestnetFlag,
		DevnetFlag,
	}
	legacyRpcFlags = []cli.Flag{
		utils.NoUSBFlag,
	}

	rpcFlags = []cli.Flag{
		utils.HTTPEnabledFlag,
		utils.HTTPListenAddrFlag,
		utils.HTTPPortFlag,
		utils.HTTPCORSDomainFlag,
		utils.HTTPVirtualHostsFlag,
		utils.GraphQLEnabledFlag,
		utils.GraphQLCORSDomainFlag,
		utils.GraphQLVirtualHostsFlag,
		utils.HTTPApiFlag,
		utils.HTTPPathPrefixFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.WSPathPrefixFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
		RPCGlobalGasCapFlag,
		RPCGlobalTxFeeCapFlag,
		RPCGlobalEVMTimeoutFlag,
		RPCGlobalTimeoutFlag,
		SubDummyBalanceFlag,
	}

	metricsFlags = []cli.Flag{
		utils.MetricsEnabledFlag,
		utils.MetricsEnabledExpensiveFlag,
		utils.MetricsHTTPFlag,
		utils.MetricsPortFlag,
		utils.MetricsEnableInfluxDBFlag,
		utils.MetricsInfluxDBEndpointFlag,
		utils.MetricsInfluxDBDatabaseFlag,
		utils.MetricsInfluxDBUsernameFlag,
		utils.MetricsInfluxDBPasswordFlag,
		utils.MetricsInfluxDBTagsFlag,
		utils.MetricsEnableInfluxDBV2Flag,
		utils.MetricsInfluxDBTokenFlag,
		utils.MetricsInfluxDBBucketFlag,
		utils.MetricsInfluxDBOrganizationFlag,
		tracing.EnableFlag,
	}

	vmFlags = []cli.Flag{
		utils.EVMInterpreterFlag,
		utils.EWASMInterpreterFlag,
	}

	nodeFlags = []cli.Flag{}
	nodeFlags = append(nodeFlags, gpoFlags...)
	nodeFlags = append(nodeFlags, accountFlags...)
	nodeFlags = append(nodeFlags, performanceFlags...)
	nodeFlags = append(nodeFlags, networkingFlags...)
	nodeFlags = append(nodeFlags, txpoolFlags...)
	nodeFlags = append(nodeFlags, artheraFlags...)
	nodeFlags = append(nodeFlags, legacyRpcFlags...)
	nodeFlags = append(nodeFlags, vmFlags...)
}

// init the CLI app.
func init() {
	discfilter.Enable()
	overrideFlags()
	overrideParams()

	initFlags()

	// App.

	app.Action = artheraMain
	app.Version = version2.VersionWithCommit()
	app.HideVersion = true // we have a command to print the version
	app.Commands = []cli.Command{
		// see p2pcmd.go:
		p2pCommand,
		// See accountcmd.go:
		accountCommand,
		walletCommand,
		// see validatorcmd.go:
		validatorCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See config.go:
		dumpConfigCommand,
		checkConfigCommand,
		// See misccmd.go:
		versionCommand,
		// See chaincmd.go
		importCommand,
		exportCommand,
		checkCommand,
		deleteCommand,
		// See snapshot.go
		snapshotCommand,
		// See dbcmd.go
		dbCommand,
		createGenesisCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, testFlags...)
	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = append(app.Flags, metricsFlags...)

	app.Before = func(ctx *cli.Context) error {
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.

		return nil
	}
}

func Launch(args []string) error {
	return app.Run(args)
}

// artheraMain is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func artheraMain(ctx *cli.Context) error {
	fmt.Println("Arthera node version: " + app.Version)
	fmt.Println("▄▀█ █▀█ ▀█▀ █░█ █▀▀ █▀█ ▄▀█\n█▀█ █▀▄ ░█░ █▀█ ██▄ █▀▄ █▀█")

	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}

	cfg := makeAllConfigs(ctx)

	genesisStore := mayGetGenesisStore(ctx)
	node, _, nodeClose := makeNode(ctx, cfg, genesisStore)

	defer nodeClose()
	startNode(ctx, node)
	node.Wait()
	return nil
}

func makeNode(ctx *cli.Context, cfg *config, genesisStore *genesisstore.Store) (*node.Node, *gossip.Service, func()) {
	// check errlock file
	errlock.SetDefaultDatadir(cfg.Node.DataDir)
	errlock.Check()

	var g *genesis.Genesis
	if genesisStore != nil {
		gv := genesisStore.Genesis()
		g = &gv
	}

	engine, dagIndex, gdb, cdb, blockProc, closeDBs := dbconfig.MakeEngine(path.Join(cfg.Node.DataDir, "chaindata"), g, cfg.AppConfigs())
	if genesisStore != nil {
		_ = genesisStore.Close()
	}
	if evmetrics.Enabled {
		metrics.SetDataDir(cfg.Node.DataDir)
	}

	// substitute default bootnodes if requested
	networkName := ""
	if gdb.HasBlockEpochState() {
		networkName = gdb.GetRules().Name
	}
	if len(networkName) == 0 && genesisStore != nil {
		networkName = genesisStore.Header().NetworkName
	}
	if needDefaultBootnodes(cfg.Node.P2P.BootstrapNodes) {
		bootnodes := params.Bootnodes[networkName]
		if bootnodes == nil {
			bootnodes = []string{}
		}
		setBootnodes(ctx, bootnodes, &cfg.Node)
	}

	for _, bn := range cfg.Node.P2P.BootstrapNodes {
		log.Info("Bootnode", "url", bn)
	}

	stack := makeConfigNode(ctx, &cfg.Node)

	valKeystore := valkeystore.NewDefaultFileKeystore(path.Join(getValKeystoreDir(cfg.Node), "validator"))
	valPubkey := cfg.Emitter.Validator.PubKey
	if key := getFakeValidatorKey(ctx); key != nil && cfg.Emitter.Validator.ID != 0 {
		addFakeValidatorKey(ctx, key, valPubkey, valKeystore)
		coinbase := dbconfig.SetAccountKey(stack.AccountManager(), key, "fakepassword")
		log.Info("Unlocked fake validator account", "address", coinbase.Address.Hex())
	}

	// unlock validator key
	if !valPubkey.Empty() {
		err := unlockValidatorKey(ctx, valPubkey, valKeystore)
		if err != nil {
			utils.Fatalf("Failed to unlock validator key: %v", err)
		}
	}
	signer := valkeystore.NewSigner(valKeystore)

	// Create and register a gossip network service.
	newTxPool := func(reader evmcore.StateReader) gossip.TxPool {
		if cfg.TxPool.Journal != "" {
			cfg.TxPool.Journal = stack.ResolvePath(cfg.TxPool.Journal)
		}
		return evmcore.NewTxPool(cfg.TxPool, reader.Config(), reader)
	}
	haltCheck := func(oldEpoch, newEpoch idx.Epoch, age time.Time) bool {
		stop := ctx.GlobalIsSet(ExitWhenAgeFlag.Name) && ctx.GlobalDuration(ExitWhenAgeFlag.Name) >= time.Since(age)
		stop = stop || ctx.GlobalIsSet(ExitWhenEpochFlag.Name) && idx.Epoch(ctx.GlobalUint64(ExitWhenEpochFlag.Name)) <= newEpoch
		if stop {
			go func() {
				// do it in a separate thread to avoid deadlock
				_ = stack.Close()
			}()
			return true
		}
		return false
	}
	svc, err := gossip.NewService(stack, cfg.Arthera, gdb, blockProc, engine, dagIndex, newTxPool, haltCheck)
	if err != nil {
		utils.Fatalf("Failed to create the service: %v", err)
	}
	err = engine.StartFrom(svc.GetConsensusCallbacks(), gdb.GetEpoch(), gdb.GetValidators())
	if err != nil {
		utils.Fatalf("Failed to start the engine: %v", err)
	}
	svc.ReprocessEpochEvents()
	if cfg.Emitter.Validator.ID != 0 {
		svc.RegisterEmitter(emitter.NewEmitter(cfg.Emitter, svc.EmitterWorld(signer)))
	}

	stack.RegisterAPIs(svc.APIs())
	stack.RegisterProtocols(svc.Protocols())
	stack.RegisterLifecycle(svc)

	return stack, svc, func() {
		_ = stack.Close()
		gdb.Close()
		_ = cdb.Close()
		if closeDBs != nil {
			_ = closeDBs()
		}
	}
}

func makeConfigNode(ctx *cli.Context, cfg *node.Config) *node.Node {
	stack, err := node.New(cfg)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}

	// Node doesn't by default populate account manager backends
	if err := setAccountManagerBackends(stack); err != nil {
		utils.Fatalf("Failed to set account manager backends: %v", err)
	}

	return stack
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces.
func startNode(ctx *cli.Context, stack *node.Node) {
	debug.Memsize.Add("node", stack)

	// Start up the node itself
	utils.StartNode(ctx, stack)

	if !ctx.GlobalIsSet(utils.MetricsInfluxDBBucketFlag.Name) {
		nodeid := fmt.Sprintf("%x", crypto.FromECDSAPub(stack.Server().LocalNode().Node().Pubkey())[1:])
		bucket := fmt.Sprintf("%s%s%s", nodeid[0:4], nodeid[len(nodeid)/2:len(nodeid)/2+4], nodeid[len(nodeid)-4:])
		err := ctx.Set(utils.MetricsInfluxDBBucketFlag.Name, bucket)
		if err != nil {
			log.Warn("Could not set the influxdb bucket. Please set it manually by adding --influxdb.metrics.bucket "+bucket, "error", err.Error())
		}
	}

	// Start metrics export if enabled
	utils.SetupMetrics(ctx)
	// Start system runtime metrics collection
	go evmetrics.CollectProcessMetrics(3 * time.Second)

	// Unlock any account specifically requested
	unlockAccounts(ctx, stack)

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	// Create a client to interact with local arthera node.
	rpcClient, err := stack.Attach()
	if err != nil {
		utils.Fatalf("Failed to attach to self: %v", err)
	}
	ethClient := ethclient.NewClient(rpcClient)

	go func() {
		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, ethClient)

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

}

// unlockAccounts unlocks any account specifically requested.
func unlockAccounts(ctx *cli.Context, stack *node.Node) {
	var unlocks []string
	inputs := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for _, input := range inputs {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	// Short circuit if there is no account to unlock.
	if len(unlocks) == 0 {
		return
	}
	// If insecure account unlocking is not allowed if node's APIs are exposed to external.
	// Print warning log to user and skip unlocking.
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		utils.Fatalf("Account unlock with HTTP access is forbidden!")
	}
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords := utils.MakePasswordList(ctx)
	for i, account := range unlocks {
		unlockAccount(ks, account, i, passwords)
	}
}

func setAccountManagerBackends(stack *node.Node) error {
	conf := stack.Config()
	am := stack.AccountManager()
	keydir := stack.KeyStoreDir()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if conf.UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	// Assemble the supported backends
	if len(conf.ExternalSigner) > 0 {
		log.Info("Using external signer", "url", conf.ExternalSigner)
		if extapi, err := external.NewExternalBackend(conf.ExternalSigner); err == nil {
			am.AddBackend(extapi)
			return nil
		} else {
			return fmt.Errorf("error connecting to external signer: %v", err)
		}
	}

	// For now, we're using EITHER external signer OR local signers.
	// If/when we implement some form of lockfile for USB and keystore wallets,
	// we can have both, but it's very confusing for the user to see the same
	// accounts in both externally and locally, plus very racey.
	am.AddBackend(keystore.NewKeyStore(keydir, scryptN, scryptP))
	if conf.USB {
		// Start a USB hub for Ledger hardware wallets
		if ledgerhub, err := usbwallet.NewLedgerHub(); err != nil {
			log.Warn(fmt.Sprintf("Failed to start Ledger hub, disabling: %v", err))
		} else {
			am.AddBackend(ledgerhub)
		}
		// Start a USB hub for Trezor hardware wallets (HID version)
		if trezorhub, err := usbwallet.NewTrezorHubWithHID(); err != nil {
			log.Warn(fmt.Sprintf("Failed to start HID Trezor hub, disabling: %v", err))
		} else {
			am.AddBackend(trezorhub)
		}
		// Start a USB hub for Trezor hardware wallets (WebUSB version)
		if trezorhub, err := usbwallet.NewTrezorHubWithWebUSB(); err != nil {
			log.Warn(fmt.Sprintf("Failed to start WebUSB Trezor hub, disabling: %v", err))
		} else {
			am.AddBackend(trezorhub)
		}
	}
	if len(conf.SmartCardDaemonPath) > 0 {
		// Start a smart card hub
		if schub, err := scwallet.NewHub(conf.SmartCardDaemonPath, scwallet.Scheme, keydir); err != nil {
			log.Warn(fmt.Sprintf("Failed to start smart card hub, disabling: %v", err))
		} else {
			am.AddBackend(schub)
		}
	}

	return nil
}
