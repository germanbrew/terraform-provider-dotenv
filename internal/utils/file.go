package utils

import (
	"fmt"
	"os"

	"github.com/direnv/direnv/v2/pkg/dotenv"
)

func ParseDotEnvFile(filename string) (map[string]string, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read dotenv file: %w", err)
	}

	parsedEntries, err := dotenv.Parse(string(contents))
	if err != nil {
		return nil, fmt.Errorf("failed to parse dotenv file: %w", err)
	}

	return parsedEntries, nil
}
