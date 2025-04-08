package i18n

import (
	"encoding/json"
	"fmt"
	"os"
)

// FlattenJSON flattens a nested map into a flat map using the given delimiter.
// If delimiter is empty, it defaults to "."
func FlattenJSON(input map[string]interface{}, delimiter string) map[string]string {
	if delimiter == "" {
		delimiter = "."
	}
	flatMap := make(map[string]string)
	flattenRecursive("", input, flatMap, delimiter)
	return flatMap
}

func flattenRecursive(prefix string, m map[string]interface{}, flatMap map[string]string, delimiter string) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + delimiter + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			flattenRecursive(fullKey, child, flatMap, delimiter)
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
	return FlattenJSON(nested, ""), nil
}
