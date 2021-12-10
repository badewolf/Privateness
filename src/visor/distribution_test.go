package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ness-network/privateness/src/cipher"
	"github.com/ness-network/privateness/src/coin"
	"github.com/ness-network/privateness/src/params"
)

func TestTransactionIsLocked(t *testing.T) {
	test := func(addrStr string, expectedIsLocked bool) {
		addr := cipher.MustDecodeBase58Address(addrStr)

		uxOut := coin.UxOut{
			Body: coin.UxBody{
				Address: addr,
			},
		}
		uxArray := coin.UxArray{uxOut}

		isLocked := TransactionIsLocked(uxArray)
		require.Equal(t, expectedIsLocked, isLocked)
	}

	for _, a := range params.GetLockedDistributionAddresses() {
		test(a, true)
	}

	for _, a := range params.GetUnlockedDistributionAddresses() {
		test(a, false)
	}

	// A random address should not be locked
	pubKey, _ := cipher.GenerateKeyPair()
	addr := cipher.AddressFromPubKey(pubKey)
	test(addr.String(), false)
}
