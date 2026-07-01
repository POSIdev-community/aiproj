package v1_10_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	v110 "github.com/POSIdev-community/aiproj/versions/v1_10"
	v19 "github.com/POSIdev-community/aiproj/versions/v1_9"
)

func TestMigrateToV110_CommonV19EnumValuesRemainValid(t *testing.T) {
	t.Parallel()

	javaVersions := []string{"8", "11", "17", "21"}
	languages := []string{
		"Java", "CSharp (Windows, Linux)", "CSharp (Windows)", "Php", "JavaScript",
		"Python", "ObjectiveC", "Swift", "CAndCPlusPlus", "Go",
		"Kotlin", "Sql", "Ruby", "Solidity", "Scala",
	}
	scanModules := []string{
		"Configuration", "Components", "BlackBox",
		"PatternMatching", "StaticCodeAnalysis", "SoftwareCompositionAnalysis",
	}
	versions := []string{"1.0", "1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9"}

	for _, version := range versions {
		for _, javaVersion := range javaVersions {
			for _, language := range languages {
				for _, module := range scanModules {
					testName := fmt.Sprintf("version_%s_java_%s_lang_%s_module_%s", version, javaVersion, sanitize(language), sanitize(module))

					t.Run(testName, func(t *testing.T) {
						input := fmt.Appendf(nil, `{
							"Version": %q,
							"ProjectName": "demo",
							"JavaSettings": {"Version": %q},
							"ProgrammingLanguages": [%q],
							"ScanModules": [%q]
						}`, version, javaVersion, language, module)

						require.NoError(t, v19.Validate(input))

						migrated, err := v19.MigrateToV110(input)
						require.NoError(t, err)
						require.NoError(t, v110.Validate(migrated))
					})
				}
			}
		}
	}
}

func sanitize(s string) string {
	replacer := strings.NewReplacer(" ", "_", "(", "", ")", "", ",", "", "-", "_")
	return replacer.Replace(s)
}

func TestMigrateToV110_RemovesPygrepSettings(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.9",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Python"],
		"ScanModules": ["StaticCodeAnalysis"],
		"PygrepSettings": {
			"CustomParameters": "foo",
			"RulesDirPath": "/rules"
		}
	}`)

	require.NoError(t, v19.Validate(input))

	migrated, err := v19.MigrateToV110(input)
	require.NoError(t, err)
	require.NoError(t, v110.Validate(migrated))

	var payload map[string]any
	require.NoError(t, json.Unmarshal(migrated, &payload))
	_, exists := payload["PygrepSettings"]
	require.False(t, exists)
}

func TestMigrateToV110_PreservesLargeIntegerAsJSONNumber(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.9",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["JavaScript"],
		"ScanModules": ["StaticCodeAnalysis"],
		"ProxySettings": {"Port": 1234567890123456789}
	}`)

	migrated, err := v19.MigrateToV110(input)
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
