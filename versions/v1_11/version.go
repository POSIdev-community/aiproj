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
