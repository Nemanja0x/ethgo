package wallet

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo"
	"pgregory.net/rapid"
)

func TestSigner_SignAndRecover(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		istyped := rapid.Bool().Draw(t, "istyped")

		// fill in common types for a transaction
		txn := &ethgo.Transaction{}

		if rapid.Bool().Draw(t, "to") {
			to := ethgo.BytesToAddress(rapid.SliceOf(rapid.Byte()).Draw(t, "to_addr"))
			txn.To = &to
		}

		// fill in specific fields depending on the type
		// of the transaction.
		if istyped {
			txn.Type = ethgo.TransactionAccessList
		}

		// signer is from a random chain
		chainid := rapid.Uint64().Draw(t, "chainid")
		signer := NewEIP155Signer(chainid)

		key, err := GenerateKey()
		require.NoError(t, err)

		signedTxn, err := signer.SignTx(txn, key)
		require.NoError(t, err)

		// recover the sender
		sender, err := signer.RecoverSender(signedTxn)
		require.NoError(t, err)

		require.Equal(t, sender, key.Address())
	})
}

func TestSigner_EIP1155(t *testing.T) {
	signer1 := NewEIP155Signer(1337)

	addr0 := ethgo.Address{0x1}
	key, err := GenerateKey()
	assert.NoError(t, err)

	txn := &ethgo.Transaction{
		To:       &addr0,
		Value:    big.NewInt(10),
		GasPrice: 0,
	}
	txn, err = signer1.SignTx(txn, key)
	assert.NoError(t, err)

	from, err := signer1.RecoverSender(txn)
	assert.NoError(t, err)
	assert.Equal(t, from, key.addr)
}

func TestTrimBytesZeros(t *testing.T) {
	assert.Equal(t, trimBytesZeros([]byte{0x1, 0x2}), []byte{0x1, 0x2})
	assert.Equal(t, trimBytesZeros([]byte{0x0, 0x1}), []byte{0x1})
	assert.Equal(t, trimBytesZeros([]byte{0x0, 0x0}), []byte{})
}
