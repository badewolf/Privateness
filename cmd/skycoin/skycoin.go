/*
skycoin daemon
*/
package main

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

import (
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/ness-network/privateness/src/fiber"
	"github.com/ness-network/privateness/src/readable"
	"github.com/ness-network/privateness/src/skycoin"
	"github.com/ness-network/privateness/src/util/logging"

	// register the supported wallets
	_ "github.com/ness-network/privateness/src/wallet/bip44wallet"
	_ "github.com/ness-network/privateness/src/wallet/collection"
	_ "github.com/ness-network/privateness/src/wallet/deterministic"
	_ "github.com/ness-network/privateness/src/wallet/xpubwallet"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "0.27.1"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""
	// ConfigMode (possible values are "", "STANDALONE_CLIENT").
	// This is used to change the default configuration.
	// Can be set by -ldflags
	ConfigMode = ""

	logger = logging.MustGetLogger("main")

	// CoinName name of coin
	CoinName = "skycoin"

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "05d4045854103f8a8938bb701cc4101c38942a180ba02d328d6f880bf37b387c47e95813d061f94bdf5d894bfebf17f933c5fc92fc9d010480765257c3d19d9b00"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "24GJTLPMoz61sV4J4qg1n14x5qqDwXqyJJy"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "02933015bd2fa1e0a885c05fb08eb7c647bf8c3188ed5120b51d0d09ccaf525036"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1637895025
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 165000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"192.243.100.192:6006",
		"167.114.97.165:6006",
		"198.245.62.172:6006",
		"151.80.37.6:6006",
		"94.23.32.95:6006",
	}

	nodeConfig = skycoin.NewNodeConfig(ConfigMode, fiber.NodeConfig{
		CoinName:            CoinName,
		GenesisSignatureStr: GenesisSignatureStr,
		GenesisAddressStr:   GenesisAddressStr,
		GenesisCoinVolume:   GenesisCoinVolume,
		GenesisTimestamp:    GenesisTimestamp,
		BlockchainPubkeyStr: BlockchainPubkeyStr,
		BlockchainSeckeyStr: BlockchainSeckeyStr,
		DefaultConnections:  DefaultConnections,
		PeerListURL:         "http://privateness.coin/blockchain/peers.txt",
		Port:                6006,
		WebInterfacePort:    6420,
		DataDirectory:       "$HOME/.skycoin",

		UnconfirmedBurnFactor:          10,
		UnconfirmedMaxTransactionSize:  32768,
		UnconfirmedMaxDropletPrecision: 3,
		CreateBlockBurnFactor:          10,
		CreateBlockMaxTransactionSize:  32768,
		CreateBlockMaxDropletPrecision: 3,
		MaxBlockTransactionsSize:       32768,

		DisplayName:           "Privateness",
		Ticker:                "NESS",
		CoinHoursName:         "Coin Hours",
		CoinHoursNameSingular: "Coin Hour",
		CoinHoursTicker:       "NCH",
		QrURIPrefix:           "privateness",
		ExplorerURL:           "https://explorer.privateness.network",
		VersionURL:            "https://version.skycoin.com/skycoin/version.txt",
		Bip44Coin:             8000,
	})

	parseFlags = true
)

func init() {
	nodeConfig.RegisterFlags()
}

func main() {
	if parseFlags {
		flag.Parse()
	}

	// create a new fiber coin instance
	coin := skycoin.NewCoin(skycoin.Config{
		Node: nodeConfig,
		Build: readable.BuildInfo{
			Version: Version,
			Commit:  Commit,
			Branch:  Branch,
		},
	}, logger)

	// parse config values
	if err := coin.ParseConfig(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// run fiber coin node
	if err := coin.Run(); err != nil {
		os.Exit(1)
	}
}
