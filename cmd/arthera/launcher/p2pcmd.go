package launcher

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"
	"net"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	p2pCommand = cli.Command{
		Name:     "p2p",
		Usage:    "P2P commands",
		Category: "P2P COMMANDS",
		Subcommands: []cli.Command{
			{
				Name:      "genkey",
				Usage:     "Generates node key files",
				ArgsUsage: "keyfile",
				Action:    genkey,
			},
			{
				Name:      "enodeurl",
				Usage:     "Returns an enode URL from a node key file",
				ArgsUsage: "keyfile",
				Action:    keyToURL,
				Flags:     []cli.Flag{hostFlag, tcpPortFlag, udpPortFlag},
			},
			{
				Name:      "ping",
				Usage:     "Sends ping to a node",
				Action:    discv4Ping,
				ArgsUsage: "<node>",
			},
			{
				Name:      "requestenr",
				Usage:     "Requests a node record using EIP-868 enrRequest",
				Action:    discv4RequestRecord,
				ArgsUsage: "<node>",
			},
			{
				Name:      "resolve",
				Usage:     "Finds a node in the DHT",
				Action:    discv4Resolve,
				ArgsUsage: "<node>",
				Flags:     []cli.Flag{bootnodesFlag},
			},
			{
				Name:   "crawl",
				Usage:  "Updates a nodes.json file with random nodes found in the DHT",
				Action: discv4Crawl,
				Flags:  []cli.Flag{bootnodesFlag, crawlTimeoutFlag},
			},
		},
	}
)

var (
	hostFlag = cli.StringFlag{
		Name:  "ip",
		Usage: "IP address of the node",
		Value: "127.0.0.1",
	}
	tcpPortFlag = cli.IntFlag{
		Name:  "tcp",
		Usage: "TCP port of the node",
		Value: 6534,
	}
	udpPortFlag = cli.IntFlag{
		Name:  "udp",
		Usage: "UDP port of the node",
		Value: 6534,
	}
	bootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "Comma separated nodes used for bootstrapping",
	}
	nodekeyFlag = cli.StringFlag{
		Name:  "nodekey",
		Usage: "Hex-encoded node key",
	}
	nodedbFlag = cli.StringFlag{
		Name:  "nodedb",
		Usage: "Nodes database location",
	}
	listenAddrFlag = cli.StringFlag{
		Name:  "addr",
		Usage: "Listening address",
	}
	crawlTimeoutFlag = cli.DurationFlag{
		Name:  "timeout",
		Usage: "Time limit for the crawl.",
		Value: 1 * time.Minute,
	}
)

func genkey(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("error: need key file as argument")
	}
	file := ctx.Args().Get(0)
	key, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("could not generate key: %v", err)
	}
	return crypto.SaveECDSA(file, key)
}

func keyToURL(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("error: need key file as argument")
	}
	var (
		file = ctx.Args().Get(0)
		host = ctx.String(hostFlag.Name)
		tcp  = ctx.Int(tcpPortFlag.Name)
		udp  = ctx.Int(udpPortFlag.Name)
	)
	key, err := crypto.LoadECDSA(file)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("invalid IP address %q", host)
	}
	node := enode.NewV4(&key.PublicKey, ip, tcp, udp)
	fmt.Println(node.URLv4())
	return nil
}

func discv4Ping(ctx *cli.Context) error {
	n := getNodeArg(ctx)
	disc := startV4(ctx)
	defer disc.Close()

	start := time.Now()
	if err := disc.Ping(n); err != nil {
		return fmt.Errorf("node didn't respond: %v", err)
	}
	fmt.Printf("node responded to ping (RTT %v).\n", time.Since(start))
	return nil
}

func discv4RequestRecord(ctx *cli.Context) error {
	n := getNodeArg(ctx)
	disc := startV4(ctx)
	defer disc.Close()

	respN, err := disc.RequestENR(n)
	if err != nil {
		return fmt.Errorf("can't retrieve record: %v", err)
	}
	fmt.Println(respN.String())
	return nil
}

func discv4Resolve(ctx *cli.Context) error {
	n := getNodeArg(ctx)
	disc := startV4(ctx)
	defer disc.Close()

	fmt.Println(disc.Resolve(n).String())
	return nil
}

func discv4Crawl(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("need nodes file as argument")
	}
	nodesFile := ctx.Args().First()
	var inputSet nodeSet
	if common.FileExist(nodesFile) {
		inputSet = loadNodesJSON(nodesFile)
	}

	disc := startV4(ctx)
	defer disc.Close()
	c := newCrawler(inputSet, disc, disc.RandomNodes())
	c.revalidateInterval = 10 * time.Minute
	output := c.run(ctx.Duration(crawlTimeoutFlag.Name))
	writeNodesJSON(nodesFile, output)
	return nil
}

// getNodeArg handles the common case of a single node descriptor argument.
func getNodeArg(ctx *cli.Context) *enode.Node {
	if ctx.NArg() < 1 {
		exit("missing node as command-line argument")
	}
	n, err := parseNode(ctx.Args()[0])
	if err != nil {
		exit(err)
	}
	return n
}

func makeDiscoveryConfig(ctx *cli.Context) (*enode.LocalNode, discover.Config) {
	var cfg discover.Config

	if ctx.IsSet(nodekeyFlag.Name) {
		key, err := crypto.HexToECDSA(ctx.String(nodekeyFlag.Name))
		if err != nil {
			exit(fmt.Errorf("-%s: %v", nodekeyFlag.Name, err))
		}
		cfg.PrivateKey = key
	} else {
		cfg.PrivateKey, _ = crypto.GenerateKey()
	}

	if commandHasFlag(ctx, bootnodesFlag) {
		bn, err := parseBootnodes(ctx)
		if err != nil {
			exit(err)
		}
		cfg.Bootnodes = bn
	}

	dbpath := ctx.String(nodedbFlag.Name)
	db, err := enode.OpenDB(dbpath)
	if err != nil {
		exit(err)
	}
	ln := enode.NewLocalNode(db, cfg.PrivateKey)
	return ln, cfg
}

// startV4 starts an ephemeral discovery V4 node.
func startV4(ctx *cli.Context) *discover.UDPv4 {
	ln, config := makeDiscoveryConfig(ctx)
	socket := listen(ln, ctx.String(listenAddrFlag.Name))
	disc, err := discover.ListenV4(socket, ln, config)
	if err != nil {
		exit(err)
	}
	return disc
}

func listen(ln *enode.LocalNode, addr string) *net.UDPConn {
	if addr == "" {
		addr = "0.0.0.0:0"
	}
	socket, err := net.ListenPacket("udp4", addr)
	if err != nil {
		exit(err)
	}
	usocket := socket.(*net.UDPConn)
	uaddr := socket.LocalAddr().(*net.UDPAddr)
	if uaddr.IP.IsUnspecified() {
		ln.SetFallbackIP(net.IP{127, 0, 0, 1})
	} else {
		ln.SetFallbackIP(uaddr.IP)
	}
	ln.SetFallbackUDP(uaddr.Port)
	return usocket
}

func parseBootnodes(ctx *cli.Context) ([]*enode.Node, error) {
	s := params.RinkebyBootnodes
	if ctx.IsSet(bootnodesFlag.Name) {
		input := ctx.String(bootnodesFlag.Name)
		if input == "" {
			return nil, nil
		}
		s = strings.Split(input, ",")
	}
	nodes := make([]*enode.Node, len(s))
	var err error
	for i, record := range s {
		nodes[i], err = parseNode(record)
		if err != nil {
			return nil, fmt.Errorf("invalid bootstrap node: %v", err)
		}
	}
	return nodes, nil
}

// commandHasFlag returns true if the current command supports the given flag.
func commandHasFlag(ctx *cli.Context, flag cli.Flag) bool {
	flags := ctx.FlagNames()
	sort.Strings(flags)
	i := sort.SearchStrings(flags, flag.GetName())
	return i != len(flags) && flags[i] == flag.GetName()
}

func exit(err interface{}) {
	if err == nil {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
