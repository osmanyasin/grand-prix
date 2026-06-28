package io

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadLevel[T any](filepath string) (*T, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	var config T
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &config, nil
}

func WriteStrategy(filepath string, strategy any) error {
	data, err := json.MarshalIndent(strategy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal strategy: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write submission to %s: %w", filepath, err)
	}

	return nil
}
