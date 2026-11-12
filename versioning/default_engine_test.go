package versioning_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aiproj/versioning"
	v112 "github.com/POSIdev-community/aiproj/versions/v1_12"
)

func TestMustBuildDefaultEngine_EndToEndMigrationFrom18ToLatest(t *testing.T) {
	t.Parallel()

	engine, _ := versioning.DefaultEngine()
	input := []byte(`{
		"Version":"1.8",
		"ProjectName":"demo",
		"ProgrammingLanguages":["JavaScript"],
		"ScanModules":["StaticCodeAnalysis"]
	}`)

	detected, err := engine.Detect(input)
	require.NoError(t, err)
	require.Equal(t, "1.8", detected)

	migrated, err := engine.MigrateToLatest(input, detected)
	require.NoError(t, err)
	require.NoError(t, v112.Validate(migrated))
}
