package versions

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJSON(schemaData, jsonData []byte) error {
	result, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(schemaData),
		gojsonschema.NewBytesLoader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("validate json against schema: %w", err)
	}

	if result.Valid() {
		return nil
	}

	errorsText := make([]string, 0, len(result.Errors()))
	for _, validationErr := range result.Errors() {
		errorsText = append(errorsText, validationErr.String())
	}

	return fmt.Errorf("schema validation failed: %s", strings.Join(errorsText, "; "))
}
