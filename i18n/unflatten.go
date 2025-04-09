package i18n

import (
	"encoding/json"
	"fmt"
	"strings"
)

// transformMapToJSON takes a map[string]map[string]string and returns a JSON string
// with dot-separated keys converted into a nested structure.
func transformMapToJSON(input map[string]map[string]string) (string, error) {
	// Create the result map to hold the nested structure
	result := make(map[string]interface{})

	// Iterate over the input map
	for key, valueMap := range input {
		// Split the key by the dot separator
		parts := strings.Split(key, ".")

		// Create the nested structure by calling a helper function
		currentMap := result
		for i := 0; i < len(parts)-1; i++ {
			// Create a nested map if it doesn't exist
			if _, exists := currentMap[parts[i]]; !exists {
				currentMap[parts[i]] = make(map[string]interface{})
			}
			currentMap = currentMap[parts[i]].(map[string]interface{})
		}

		// The last part of the key is the field
		field := parts[len(parts)-1]
		// Add the "value" and "tooltip" to the field
		for key, val := range valueMap {
			if key == "value" {
				currentMap[field] = val
			} else if key == "tooltip" {
				currentMap[field+"_tooltip"] = val
			}
		}
	}

	// Marshal the result map to a JSON string
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonData), nil
}
