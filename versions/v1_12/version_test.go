package v1_12_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aiproj/versions/v1_11"
	"github.com/POSIdev-community/aiproj/versions/v1_12"
)

func TestMigrateToV112_AllV111VersionEnumValues(t *testing.T) {
	t.Parallel()

	versions := []string{
		"1.0", "1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9", "1.10", "1.11",
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

			require.NoError(t, v1_11.Validate(input))

			migrated, err := v1_11.MigrateToV112(input)
			require.NoError(t, err)
			require.NoError(t, v1_12.Validate(migrated))
		})
	}
}

func TestMigrateToV112_RemovesJavaScriptDependenciesPath(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["JavaScript"],
		"ScanModules": ["StaticCodeAnalysis"],
		"JavaScriptSettings": {
			"DependenciesPath": "./node_modules"
		}
	}`)

	require.NoError(t, v1_11.Validate(input))

	migrated, err := v1_11.MigrateToV112(input)
	require.NoError(t, err)
	require.NoError(t, v1_12.Validate(migrated))
	require.NotContains(t, string(migrated), "DependenciesPath")
}

func TestValidate_AcceptsRaceConditionAndAsintModules(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.12",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Java"],
		"ScanModules": ["RaceConditionDetection", "Asint"]
	}`)

	require.NoError(t, v1_12.Validate(input))
}

func TestValidate_AcceptsAutomatedAiTriageFields(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.12",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Java"],
		"ScanModules": ["StaticCodeAnalysis"],
		"AutomatedAiTriage": true,
		"AutomatedAiTriageMinSeverity": "High"
	}`)

	require.NoError(t, v1_12.Validate(input))
}

func TestValidate_AcceptsDotNetDependenciesPath(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.12",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["CSharp (Windows, Linux)"],
		"ScanModules": ["StaticCodeAnalysis"],
		"DotNetSettings": {
			"DependenciesPath": "./deps"
		}
	}`)

	require.NoError(t, v1_12.Validate(input))
}

func TestValidate_RejectsRaceConditionForV111Schema(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Java"],
		"ScanModules": ["RaceConditionDetection"]
	}`)

	require.Error(t, v1_11.Validate(input))
}
