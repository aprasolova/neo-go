package core

import (
	"testing"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenesisBlockMainNet(t *testing.T) {
	cfg, err := config.Load("../../config", config.ModeMainNet)
	require.NoError(t, err)

	block, err := createGenesisBlock(cfg.ProtocolConfiguration)
	require.NoError(t, err)

	//TODO: After we added Nonce field to transaction.Transaction, goveringTockenTx and UtilityTockenTx hashes
	// have been changed. Consequently, hash of the genesis block has been changed.
	// Update expected genesis block hash for better times.
	// Old hash is "d42561e3d30e15be6400b6df2f328e02d2bf6354c41dce433bc57687c82144bf"
	expect := "8eb36fe47f07a795a1783a9f066603db66c5b76cf878650b1e137c114f46c0cc"
	assert.Equal(t, expect, block.Hash().StringLE())
}

func TestGetConsensusAddressMainNet(t *testing.T) {
	var (
		consensusAddr   = "APtiVEdLi5GEmQ8CL5RcCE7BNcsPsxeXh7"
		consensusScript = "590c459950f1d83e67ee11fcef202a6ebb8b1a77"
	)

	cfg, err := config.Load("../../config", config.ModeMainNet)
	require.NoError(t, err)

	validators, err := getValidators(cfg.ProtocolConfiguration)
	require.NoError(t, err)

	script, err := getNextConsensusAddress(validators)
	require.NoError(t, err)

	assert.Equal(t, consensusScript, script.String())
	assert.Equal(t, consensusAddr, address.Uint160ToString(script))
}

func TestUtilityTokenTX(t *testing.T) {
	//TODO: After we added Nonce field to transaction.Transaction, UtilityTockenTx hash
	// has been changed. Update it for better times.
	// Old hash is "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
	expect := "8ef63ccd6f4ea20a93e7f4e84b2d43f778077612b241d617e42e1750cca4f2c5"
	assert.Equal(t, expect, UtilityTokenID().StringLE())
}

func TestGoverningTokenTX(t *testing.T) {
	//TODO: After we added Nonce field to transaction.Transaction, GoveringTockenTx hash
	// has been changed. Update it for better times.
	// Old hash is "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
	expect := "7dc12f8835b085d468ddbab3f7152c0e2f5679b5815b2028685abb4608e7b7dc"
	assert.Equal(t, expect, GoverningTokenID().StringLE())
}
