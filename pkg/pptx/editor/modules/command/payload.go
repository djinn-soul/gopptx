package command

import (
	"encoding/json"
	"fmt"
)

func DecodeOptionalPayloadValue(payload map[string]any, key string, target any) error {
	rawValue, ok := payload[key]
	if !ok || rawValue == nil {
		return nil
	}
	raw, err := json.Marshal(rawValue)
	if err != nil {
		return fmt.Errorf("invalid %s structure: %w", key, err)
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("invalid %s payload: %w", key, err)
	}
	return nil
}
