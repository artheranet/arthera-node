package launcher

import (
	"fmt"
	version2 "github.com/artheranet/arthera-node/version"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"

	"github.com/artheranet/arthera-node/gossip"
)

var (
	versionCommand = cli.Command{
		Action:    utils.MigrateFlags(version),
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}
)

func version(_ *cli.Context) error {
	fmt.Println("▄▀█ █▀█ ▀█▀ █░█ █▀▀ █▀█ ▄▀█\n█▀█ █▀▄ ░█░ █▀█ ██▄ █▀▄ █▀█")
	fmt.Println("Version:", version2.VersionWithCommit())
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Protocol Versions:", []uint{gossip.ProtocolVersion})
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}
