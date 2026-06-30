package versions

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type PayloadTransform func(map[string]any) error

func MigratePayloadToVersion(data []byte, toVersion string, transforms ...PayloadTransform) ([]byte, error) {
	var payload map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	for _, transform := range transforms {
		if transform == nil {
			continue
		}

		if err := transform(payload); err != nil {
			return nil, fmt.Errorf("apply payload transform: %w", err)
		}
	}

	payload["Version"] = toVersion

	migrated, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal migrated payload: %w", err)
	}

	return migrated, nil
}
