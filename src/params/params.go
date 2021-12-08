package params

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

var (
	// MainNetDistribution Skycoin mainnet coin distribution parameters
	MainNetDistribution = Distribution{
		MaxCoinSupply:        165000000,
		InitialUnlockedCount: 100,
		UnlockAddressRate:    5,
		UnlockTimeInterval:   31536000,
		Addresses: []string{
			"yDctMREirofdnxEgZLVJwaZdEb9Em1bJyk",
			"L7JU5g8zkfg3q24yqJSojSzRAZc1dbUJn5",
			"hY5By3kHpwWp3a5VnDCyrJGtPFAnnPMrqS",
			"XTRghTPUUfiz9P3LUWaBRpT7LFoTR9jSAH",
			"a5NctjK7wkpXw17adFFSYRAYcJps2FEEhQ",
			"2MSeuyaReBbyePioGhpWoQStX8W9MPAZgTu",
			"7gbFAfL6XA5Wkpa3VjiwjhZNv4HeDzgrQ3",
			"29377Ntb9C2AgumgABigzoq1RdGJDXauBEP",
			"5o7zPhehRerBP2JexqwfUAyQiqzQKsbcod",
			"2fknwEypYS8WafAd3AJgv6vAm2QtiahU5Vx",
		},
	}

	// UserVerifyTxn transaction verification parameters for user-created transactions
	UserVerifyTxn = VerifyTxn{
		// BurnFactor can be overriden with `USER_BURN_FACTOR` env var
		BurnFactor: 10,
		// MaxTransactionSize can be overriden with `USER_MAX_TXN_SIZE` env var
		MaxTransactionSize: 32768, // in bytes
		// MaxDropletPrecision can be overriden with `USER_MAX_DECIMALS` env var
		MaxDropletPrecision: 3,
	}
)
