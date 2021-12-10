/*
cli is a command line client for interacting with a skycoin node and offline wallet management
*/
package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ness-network/privateness/src/cli"
	"github.com/ness-network/privateness/src/util/logging"

	// register the supported wallets
	_ "github.com/ness-network/privateness/src/wallet/bip44wallet"
	_ "github.com/ness-network/privateness/src/wallet/collection"
	_ "github.com/ness-network/privateness/src/wallet/deterministic"
	_ "github.com/ness-network/privateness/src/wallet/xpubwallet"
)

func main() {
	logging.SetLevel(logrus.WarnLevel)

	cfg, err := cli.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	skyCLI, err := cli.NewCLI(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := skyCLI.Execute(); err != nil {
		os.Exit(1)
	}
}
