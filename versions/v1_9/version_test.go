package v1_9_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aiproj/versions/v1_8"
	"github.com/POSIdev-community/aiproj/versions/v1_9"
)

func TestMigrateToV19_AllV18VersionEnumValues(t *testing.T) {
	t.Parallel()

	versions := []string{"1.0", "1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8"}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			t.Parallel()

			input := fmt.Appendf(nil, `{
				"Version": %q,
				"ProjectName": "demo",
				"ProgrammingLanguages": ["JavaScript"],
				"ScanModules": ["StaticCodeAnalysis"]
			}`, version)

			require.NoError(t, v1_8.Validate(input))

			migrated, err := v1_8.MigrateToV19(input)
			require.NoError(t, err)
			require.NoError(t, v1_9.Validate(migrated))
		})
	}
}

func TestMigrateToV19_PreservesLargeIntegerAsJSONNumber(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.8",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["JavaScript"],
		"ScanModules": ["StaticCodeAnalysis"],
		"ProxySettings": {"Port": 1234567890123456789}
	}`)

	migrated, err := v1_8.MigrateToV19(input)
	require.NoError(t, err)

	decoder := json.NewDecoder(bytes.NewReader(migrated))
	decoder.UseNumber()

	var payload map[string]any
	require.NoError(t, decoder.Decode(&payload))

	proxy, ok := payload["ProxySettings"].(map[string]any)
	require.True(t, ok)

	port, ok := proxy["Port"].(json.Number)
	require.True(t, ok)
	require.Equal(t, "1234567890123456789", port.String())
}
