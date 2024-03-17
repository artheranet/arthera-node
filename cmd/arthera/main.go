package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/artheranet/arthera-node/cmd/arthera/launcher"

	// Force-load the native, to trigger registration
	//_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
)

func main() {
	var majorVer int
	var minorVer int
	var other string
	n, err := fmt.Sscanf(runtime.Version(), "go%d.%d%s", &majorVer, &minorVer, &other)
	if n >= 2 && err == nil {
		if (majorVer*100 + minorVer) > 119 {
			panic(runtime.Version() + " is not supported, please downgrade your go compiler to go 1.19 or older")
		}
	}

	if err := launcher.Launch(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
