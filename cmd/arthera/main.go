package main

import (
	"fmt"
	"github.com/artheranet/arthera-node/cmd/arthera/launcher"
	"os"
)

func main() {
	fmt.Println("▄▀█ █▀█ ▀█▀ █░█ █▀▀ █▀█ ▄▀█\n█▀█ █▀▄ ░█░ █▀█ ██▄ █▀▄ █▀█")
	if err := launcher.Launch(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
