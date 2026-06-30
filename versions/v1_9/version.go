package v1_9

import (
	_ "embed"
	"fmt"

	"github.com/POSIdev-community/aiproj/versions"
)

//go:embed schema.json
var schemaJSON []byte

func Validate(data []byte) error {
	if err := versions.ValidateJSON(schemaJSON, data); err != nil {
		return fmt.Errorf("validate aiproj %s schema: %w", versions.V1_9, err)
	}

	return nil
}

func MigrateToV110(data []byte) ([]byte, error) {
	migrated, err := versions.MigratePayloadToVersion(data, "1.10", removePygrepSettings)
	if err != nil {
		return nil, fmt.Errorf("migrate aiproj %s -> %s: %w", versions.V1_9, versions.V1_10, err)
	}

	return migrated, nil
}

func removePygrepSettings(payload map[string]any) error {
	delete(payload, "PygrepSettings")

	return nil
}
