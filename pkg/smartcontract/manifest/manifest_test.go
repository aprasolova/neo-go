package manifest

import (
	"encoding/json"
	"testing"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/stretchr/testify/require"
)

// Test vectors are taken from the main NEO repo
// https://github.com/neo-project/neo/blob/master/tests/neo.UnitTests/SmartContract/Manifest/UT_ContractManifest.cs#L10
func TestManifest_MarshalJSON(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		s := `{"groups":[],"features":{"storage":false,"payable":false},"supportedstandards":[],"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"*","methods":"*"}],"trusts":[],"safemethods":[],"extra":null}`
		m := testUnmarshalMarshalManifest(t, s)
		require.Equal(t, DefaultManifest(util.Uint160{}), m)
	})

	// this vector is missing from original repo
	t.Run("features", func(t *testing.T) {
		s := `{"groups":[],"features":{"storage":true,"payable":true},"supportedstandards":[],"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"*","methods":"*"}],"trusts":[],"safemethods":[],"extra":null}`
		testUnmarshalMarshalManifest(t, s)
	})

	t.Run("permissions", func(t *testing.T) {
		s := `{"groups":[],"features":{"storage":false,"payable":false},"supportedstandards":[],"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"0x0000000000000000000000000000000000000000","methods":["method1","method2"]}],"trusts":[],"safemethods":[],"extra":null}`
		testUnmarshalMarshalManifest(t, s)
	})

	t.Run("safe methods", func(t *testing.T) {
		s := `{"groups":[],"features":{"storage":false,"payable":false},"supportedstandards":[],"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"*","methods":"*"}],"trusts":[],"safemethods":["balanceOf"],"extra":null}`
		testUnmarshalMarshalManifest(t, s)
	})

	t.Run("trust", func(t *testing.T) {
		s := `{"groups":[],"features":{"storage":false,"payable":false},"supportedstandards":[],"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"*","methods":"*"}],"trusts":["0x0000000000000000000000000000000000000001"],"safemethods":[],"extra":null}`
		testUnmarshalMarshalManifest(t, s)
	})

	t.Run("groups", func(t *testing.T) {
		s := `{"groups":[{"pubkey":"03b209fd4f53a7170ea4444e0cb0a6bb6a53c2bd016926989cf85f9b0fba17a70c","signature":"QUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQQ=="}],"supportedstandards":[],"features":{"storage":false,"payable":false},"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"*","methods":"*"}],"trusts":[],"safemethods":[],"extra":null}`
		testUnmarshalMarshalManifest(t, s)
	})

	t.Run("extra", func(t *testing.T) {
		s := `{"groups":[],"features":{"storage":false,"payable":false},"supportedstandards":[],"abi":{"hash":"0x0000000000000000000000000000000000000000","methods":[],"events":[]},"permissions":[{"contract":"*","methods":"*"}],"trusts":[],"safemethods":[],"extra":{"key":"value"}}`
		testUnmarshalMarshalManifest(t, s)
	})
}

func testUnmarshalMarshalManifest(t *testing.T, s string) *Manifest {
	js := []byte(s)
	c := NewManifest(util.Uint160{})
	require.NoError(t, json.Unmarshal(js, c))

	data, err := json.Marshal(c)
	require.NoError(t, err)
	require.JSONEq(t, s, string(data))

	return c
}

func TestManifest_CanCall(t *testing.T) {
	t.Run("safe methods", func(t *testing.T) {
		man1 := NewManifest(util.Uint160{})
		man2 := DefaultManifest(util.Uint160{})
		require.False(t, man1.CanCall(man2, "method1"))
		man2.SafeMethods.Add("method1")
		require.True(t, man1.CanCall(man2, "method1"))
	})

	t.Run("wildcard permission", func(t *testing.T) {
		man1 := DefaultManifest(util.Uint160{})
		man2 := DefaultManifest(util.Uint160{})
		require.True(t, man1.CanCall(man2, "method1"))
	})
}

func TestPermission_IsAllowed(t *testing.T) {
	manifest := DefaultManifest(util.Uint160{})

	t.Run("wildcard", func(t *testing.T) {
		perm := NewPermission(PermissionWildcard)
		require.True(t, perm.IsAllowed(manifest, "AAA"))
	})

	t.Run("hash", func(t *testing.T) {
		perm := NewPermission(PermissionHash, util.Uint160{})
		require.True(t, perm.IsAllowed(manifest, "AAA"))

		t.Run("restrict methods", func(t *testing.T) {
			perm.Methods.Restrict()
			require.False(t, perm.IsAllowed(manifest, "AAA"))
			perm.Methods.Add("AAA")
			require.True(t, perm.IsAllowed(manifest, "AAA"))
		})
	})

	t.Run("invalid hash", func(t *testing.T) {
		perm := NewPermission(PermissionHash, util.Uint160{1})
		require.False(t, perm.IsAllowed(manifest, "AAA"))
	})

	priv, err := keys.NewPrivateKey()
	require.NoError(t, err)
	manifest.Groups = []Group{{PublicKey: priv.PublicKey()}}

	t.Run("group", func(t *testing.T) {
		perm := NewPermission(PermissionGroup, priv.PublicKey())
		require.True(t, perm.IsAllowed(manifest, "AAA"))
	})

	t.Run("invalid group", func(t *testing.T) {
		priv2, err := keys.NewPrivateKey()
		require.NoError(t, err)
		perm := NewPermission(PermissionGroup, priv2.PublicKey())
		require.False(t, perm.IsAllowed(manifest, "AAA"))
	})
}

func TestIsValid(t *testing.T) {
	contractHash := util.Uint160{1, 2, 3}
	m := NewManifest(contractHash)

	t.Run("valid, no groups", func(t *testing.T) {
		require.True(t, m.IsValid(contractHash))
	})

	t.Run("invalid, no groups", func(t *testing.T) {
		require.False(t, m.IsValid(util.Uint160{9, 8, 7}))
	})

	t.Run("with groups", func(t *testing.T) {
		m.Groups = make([]Group, 3)
		pks := make([]*keys.PrivateKey, 3)
		for i := range pks {
			pk, err := keys.NewPrivateKey()
			require.NoError(t, err)
			pks[i] = pk
			m.Groups[i] = Group{
				PublicKey: pk.PublicKey(),
				Signature: pk.Sign(contractHash.BytesBE()),
			}
		}

		t.Run("valid", func(t *testing.T) {
			require.True(t, m.IsValid(contractHash))
		})

		t.Run("invalid, wrong contract hash", func(t *testing.T) {
			require.False(t, m.IsValid(util.Uint160{4, 5, 6}))
		})

		t.Run("invalid, wrong group signature", func(t *testing.T) {
			pk, err := keys.NewPrivateKey()
			require.NoError(t, err)
			m.Groups = append(m.Groups, Group{
				PublicKey: pk.PublicKey(),
				// actually, there shouldn't be such situation, as Signature is always the signature
				// of the contract hash.
				Signature: pk.Sign([]byte{1, 2, 3}),
			})
			require.False(t, m.IsValid(contractHash))
		})
	})
}
