package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/POSIdev-community/aiproj/versioning"
)

func MigrateBytes(engine *versioning.Engine, data []byte, targetVersion string) ([]byte, error) {
	if engine == nil {
		return nil, fmt.Errorf("version engine is nil")
	}

	detectedVersion, err := engine.Detect(data)
	if err != nil {
		return nil, fmt.Errorf("detect aiproj version: %w", err)
	}

	migrated, err := engine.MigrateTo(data, detectedVersion, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("migrate aiproj %s -> %s: %w", detectedVersion, targetVersion, err)
	}

	return migrated, nil
}

func ReadFromBytes[T any](data []byte) (T, error) {
	r := bytes.NewReader(data)
	var value T

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return value, fmt.Errorf("decode aiproj: %w", err)
	}

	return value, nil
}
