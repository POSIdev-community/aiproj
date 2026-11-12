package v1_11

import (
	_ "embed"
	"fmt"

	"github.com/POSIdev-community/aiproj/versions"
)

//go:embed schema.json
var schemaJSON []byte

func Validate(data []byte) error {
	if err := versions.ValidateJSON(schemaJSON, data); err != nil {
		return fmt.Errorf("validate aiproj %s schema: %w", versions.V1_11, err)
	}

	return nil
}

func MigrateToV112(data []byte) ([]byte, error) {
	migrated, err := versions.MigratePayloadToVersion(data, versions.V1_12, removeJavaScriptDependenciesPath)
	if err != nil {
		return nil, fmt.Errorf("migrate aiproj %s -> %s: %w", versions.V1_11, versions.V1_12, err)
	}

	return migrated, nil
}

// removeJavaScriptDependenciesPath drops the key 1.12 no longer allows: the JavaScript core never used it.
func removeJavaScriptDependenciesPath(payload map[string]any) error {
	if jsSettings, ok := payload["JavaScriptSettings"].(map[string]any); ok {
		delete(jsSettings, "DependenciesPath")
	}

	return nil
}
