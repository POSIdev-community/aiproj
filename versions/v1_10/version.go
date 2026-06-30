package v1_10

import (
	_ "embed"
	"fmt"

	"github.com/POSIdev-community/aiproj/versions"
)

//go:embed schema.json
var schemaJSON []byte

func Validate(data []byte) error {
	if err := versions.ValidateJSON(schemaJSON, data); err != nil {
		return fmt.Errorf("validate aiproj %s schema: %w", versions.V1_10, err)
	}

	return nil
}

func MigrateToV111(data []byte) ([]byte, error) {
	migrated, err := versions.MigratePayloadToVersion(data, versions.V1_11)
	if err != nil {
		return nil, fmt.Errorf("migrate aiproj %s -> %s: %w", versions.V1_10, versions.V1_11, err)
	}

	return migrated, nil
}
