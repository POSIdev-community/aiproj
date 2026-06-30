package versioning_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aiproj/versioning"
)

func TestDetect_PicksFirstValidFromNewestToOldest(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: func([]byte) error { return nil }},
			{Version: "1.9", Validate: func([]byte) error { return nil }},
			{Version: "1.8", Validate: func([]byte) error { return nil }},
		},
		nil,
	)
	require.NoError(t, err)

	version, err := engine.Detect([]byte(`{"Version":""}`))
	require.NoError(t, err)
	require.Equal(t, "1.10", version)
}

func TestDetect_UsesExplicitVersionWhenSupported(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: func([]byte) error { return nil }},
			{Version: "1.9", Validate: func([]byte) error { return nil }},
		},
		nil,
	)
	require.NoError(t, err)

	version, err := engine.Detect([]byte(`{"Version":"1.9"}`))
	require.NoError(t, err)
	require.Equal(t, "1.9", version)
}

func TestDetect_ReturnsErrorForInvalidExplicitVersionPayload(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: func([]byte) error { return nil }},
			{Version: "1.9", Validate: func([]byte) error { return errors.New("invalid for 1.9") }},
		},
		nil,
	)
	require.NoError(t, err)

	_, err = engine.Detect([]byte(`{"Version":"1.9"}`))
	require.Error(t, err)
	require.ErrorContains(t, err, "validate explicit version 1.9")
}

func TestDetect_ReturnsErrorWhenNoSchemaMatches(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: func([]byte) error { return errors.New("invalid for 1.10") }},
			{Version: "1.9", Validate: func([]byte) error { return errors.New("invalid for 1.9") }},
		},
		nil,
	)
	require.NoError(t, err)

	_, err = engine.Detect([]byte(`{}`))
	require.Error(t, err)
	require.ErrorIs(t, err, versioning.ErrNoValidVersion)
}

func TestDetect_ReturnsErrorForUnsupportedExplicitVersion(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: func([]byte) error { return nil }},
			{Version: "1.9", Validate: func([]byte) error { return nil }},
		},
		nil,
	)
	require.NoError(t, err)

	_, err = engine.Detect([]byte(`{"Version":"2.0"}`))
	require.Error(t, err)
	require.ErrorIs(t, err, versioning.ErrUnsupportedVersion)
}

func TestMigrateTo_MigratesToSpecificVersion(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: hasVersion(t, "1.10")},
			{Version: "1.9", Validate: hasVersion(t, "1.9")},
			{Version: "1.8", Validate: hasVersion(t, "1.8")},
		},
		[]versioning.Migrator{
			versioning.NewFunctionMigrator("1.8", "1.9", migrateTo("1.9")),
			versioning.NewFunctionMigrator("1.9", "1.10", migrateTo("1.10")),
		},
	)
	require.NoError(t, err)

	migrated, err := engine.MigrateTo([]byte(`{"Version":"1.8","ProjectName":"demo"}`), "1.8", "1.9")
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(migrated, &payload))
	require.Equal(t, "1.9", payload["Version"])
}

func TestMigrateTo_StopsAtTargetWithoutGoingToLatest(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.11", Validate: hasVersion(t, "1.11")},
			{Version: "1.10", Validate: hasVersion(t, "1.10")},
			{Version: "1.9", Validate: hasVersion(t, "1.9")},
			{Version: "1.8", Validate: hasVersion(t, "1.8")},
		},
		[]versioning.Migrator{
			versioning.NewFunctionMigrator("1.8", "1.9", migrateTo("1.9")),
			versioning.NewFunctionMigrator("1.9", "1.10", migrateTo("1.10")),
			versioning.NewFunctionMigrator("1.10", "1.11", migrateTo("1.11")),
		},
	)
	require.NoError(t, err)

	migrated, err := engine.MigrateTo([]byte(`{"Version":"1.8","ProjectName":"demo"}`), "1.8", "1.10")
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(migrated, &payload))
	require.Equal(t, "1.10", payload["Version"])
}

func TestMigrateTo_ReturnsErrorOnDowngrade(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: hasVersion(t, "1.10")},
			{Version: "1.9", Validate: hasVersion(t, "1.9")},
		},
		nil,
	)
	require.NoError(t, err)

	_, err = engine.MigrateTo([]byte(`{"Version":"1.10"}`), "1.10", "1.9")
	require.Error(t, err)
	require.ErrorContains(t, err, "cannot downgrade")
}

func TestMigrateToLatest_ChainsAllSteps(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: hasVersion(t, "1.10")},
			{Version: "1.9", Validate: hasVersion(t, "1.9")},
			{Version: "1.8", Validate: hasVersion(t, "1.8")},
		},
		[]versioning.Migrator{
			versioning.NewFunctionMigrator("1.8", "1.9", migrateTo("1.9")),
			versioning.NewFunctionMigrator("1.9", "1.10", migrateTo("1.10")),
		},
	)
	require.NoError(t, err)

	migrated, err := engine.MigrateToLatest([]byte(`{"Version":"1.8","ProjectName":"demo"}`), "1.8")
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(migrated, &payload))
	require.Equal(t, "1.10", payload["Version"])
}

func TestMigrateToLatest_ValidatesStartVersionBeforeMigration(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: hasVersion(t, "1.10")},
			{Version: "1.9", Validate: hasVersion(t, "1.9")},
		},
		[]versioning.Migrator{
			versioning.NewFunctionMigrator("1.9", "1.10", migrateTo("1.10")),
		},
	)
	require.NoError(t, err)

	_, err = engine.MigrateToLatest([]byte(`{"Version":"1.8"}`), "1.9")
	require.Error(t, err)
	require.ErrorContains(t, err, "validate data as 1.9")
}

func TestMigrateToLatest_ValidatesEachStep(t *testing.T) {
	t.Parallel()

	engine, err := versioning.NewEngine(
		[]versioning.VersionSpec{
			{Version: "1.10", Validate: hasVersion(t, "1.10")},
			{Version: "1.9", Validate: func([]byte) error { return errors.New("invalid intermediate") }},
			{Version: "1.8", Validate: hasVersion(t, "1.8")},
		},
		[]versioning.Migrator{
			versioning.NewFunctionMigrator("1.8", "1.9", migrateTo("1.9")),
			versioning.NewFunctionMigrator("1.9", "1.10", migrateTo("1.10")),
		},
	)
	require.NoError(t, err)

	_, err = engine.MigrateToLatest([]byte(`{"Version":"1.8"}`), "1.8")
	require.Error(t, err)
	require.ErrorContains(t, err, "validate migrated data as 1.9")
}

func hasVersion(t *testing.T, expected string) versioning.Validator {
	t.Helper()

	return func(data []byte) error {
		var payload map[string]any
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}

		if payload["Version"] != expected {
			return errors.New("unexpected version")
		}

		return nil
	}
}

func migrateTo(nextVersion string) func(data []byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		var payload map[string]any
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}

		payload["Version"] = nextVersion

		return json.Marshal(payload)
	}
}
