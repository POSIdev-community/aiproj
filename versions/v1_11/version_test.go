package v1_11_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aiproj/versions/v1_10"
	"github.com/POSIdev-community/aiproj/versions/v1_11"
)

func TestMigrateToV111_AllV110VersionEnumValues(t *testing.T) {
	t.Parallel()

	versions := []string{
		"1.0", "1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9", "1.10",
	}

	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			t.Parallel()

			input := fmt.Appendf(nil, `{
				"Version": %q,
				"ProjectName": "demo",
				"ProgrammingLanguages": ["JavaScript"],
				"ScanModules": ["StaticCodeAnalysis"]
			}`, version)

			require.NoError(t, v1_10.Validate(input))

			migrated, err := v1_10.MigrateToV111(input)
			require.NoError(t, err)
			require.NoError(t, v1_11.Validate(migrated))
		})
	}
}

func TestMigrateToV111_CommonV110EnumValuesRemainValid(t *testing.T) {
	t.Parallel()

	javaVersions := []string{"8", "11", "17", "21", "25"}
	languages := []string{
		"Java", "CSharp (Windows, Linux)", "CSharp (Windows)", "Php", "JavaScript",
		"Python", "ObjectiveC", "Swift", "CAndCPlusPlus", "Go",
		"Kotlin", "Sql", "Ruby", "Solidity", "Scala", "OneC",
	}
	scanModules := []string{
		"Configuration", "Components", "BlackBox",
		"PatternMatching", "StaticCodeAnalysis", "SoftwareCompositionAnalysis",
	}

	for _, javaVersion := range javaVersions {
		for _, language := range languages {
			for _, module := range scanModules {
				testName := fmt.Sprintf("java_%s_lang_%s_module_%s", javaVersion, sanitize(language), sanitize(module))

				t.Run(testName, func(t *testing.T) {
					input := fmt.Appendf(nil, `{
						"Version": "1.10",
						"ProjectName": "demo",
						"JavaSettings": {"Version": %q},
						"ProgrammingLanguages": [%q],
						"ScanModules": [%q]
					}`, javaVersion, language, module)

					require.NoError(t, v1_10.Validate(input))

					migrated, err := v1_10.MigrateToV111(input)
					require.NoError(t, err)
					require.NoError(t, v1_11.Validate(migrated))
				})
			}
		}
	}
}

func TestValidate_AcceptsDartLanguage(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Dart"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	require.NoError(t, v1_11.Validate(input))
}

func TestValidate_RejectsDartLanguageForV110Schema(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.10",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Dart"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	require.Error(t, v1_10.Validate(input))
}

func TestMigrateToV111_PreservesLargeIntegerAsJSONNumber(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.10",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["JavaScript"],
		"ScanModules": ["StaticCodeAnalysis"],
		"ProxySettings": {"Port": 1234567890123456789}
	}`)

	migrated, err := v1_10.MigrateToV111(input)
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

func sanitize(s string) string {
	replacer := strings.NewReplacer(" ", "_", "(", "", ")", "", ",", "", "-", "_")
	return replacer.Replace(s)
}
