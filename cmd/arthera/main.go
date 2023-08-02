package main

import (
	"fmt"
	"os"

	"github.com/artheranet/arthera-node/cmd/arthera/launcher"
)

func main() {
	if err := launcher.Launch(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
