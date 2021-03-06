package payload

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nspcc-dev/neo-go/pkg/core/block"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/hash"
	"github.com/nspcc-dev/neo-go/pkg/internal/testserdes"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

func newDumbBlock() *block.Base {
	return &block.Base{
		Version:       0,
		PrevHash:      hash.Sha256([]byte("a")),
		MerkleRoot:    hash.Sha256([]byte("b")),
		Timestamp:     100500,
		Index:         1,
		NextConsensus: hash.Hash160([]byte("a")),
		Script: transaction.Witness{
			VerificationScript: []byte{0x51}, // PUSH1
			InvocationScript:   []byte{0x61}, // NOP
		},
	}
}

func TestMerkleBlock_EncodeDecodeBinary(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		b := newDumbBlock()
		_ = b.Hash()
		expected := &MerkleBlock{
			Base:    b,
			TxCount: 0,
			Hashes:  []util.Uint256{},
			Flags:   []byte{},
		}
		testserdes.EncodeDecodeBinary(t, expected, new(MerkleBlock))
	})

	t.Run("bad contents count", func(t *testing.T) {
		b := newDumbBlock()
		_ = b.Hash()
		expected := &MerkleBlock{
			Base:    b,
			TxCount: block.MaxContentsPerBlock + 1,
			Hashes:  make([]util.Uint256, block.MaxContentsPerBlock),
			Flags:   []byte{},
		}
		data, err := testserdes.EncodeBinary(expected)
		require.NoError(t, err)
		require.True(t, errors.Is(block.ErrMaxContentsPerBlock, testserdes.DecodeBinary(data, new(MerkleBlock))))
	})

	t.Run("bad flags size", func(t *testing.T) {
		b := newDumbBlock()
		_ = b.Hash()
		expected := &MerkleBlock{
			Base:    b,
			TxCount: 0,
			Hashes:  []util.Uint256{},
			Flags:   []byte{1, 2, 3, 4, 5},
		}
		data, err := testserdes.EncodeBinary(expected)
		require.NoError(t, err)
		require.Error(t, testserdes.DecodeBinary(data, new(MerkleBlock)))
	})
}
