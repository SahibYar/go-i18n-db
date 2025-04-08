package i18n

import (
	"encoding/json"
	"fmt"
	"os"
)

// FlattenJSON flattens a nested map into a flat map with pipe-separated keys.
func FlattenJSON(input map[string]interface{}) map[string]string {
	flatMap := make(map[string]string)
	flattenRecursive("", input, flatMap)
	return flatMap
}

func flattenRecursive(prefix string, m map[string]interface{}, flatMap map[string]string) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "|" + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			flattenRecursive(fullKey, child, flatMap)
		case string:
			flatMap[fullKey] = child
		default:
			flatMap[fullKey] = fmt.Sprintf("%v", child)
		}
	}
}

// LoadAndFlatten reads a JSON file and flattens its contents.
func LoadAndFlatten(filePath string) (map[string]string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var nested map[string]interface{}
	if err := json.Unmarshal(bytes, &nested); err != nil {
		return nil, err
	}
	return FlattenJSON(nested), nil
}
