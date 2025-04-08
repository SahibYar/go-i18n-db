package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]string
	}{
		{
			name: "Flat JSON",
			input: map[string]interface{}{
				"title": "Dashboard",
				"label": "Welcome",
			},
			expected: map[string]string{
				"title": "Dashboard",
				"label": "Welcome",
			},
		},
		{
			name: "Two-level JSON",
			input: map[string]interface{}{
				"topbar": map[string]interface{}{
					"profile": "My Profile",
					"logout":  "Log Out",
				},
			},
			expected: map[string]string{
				"topbar|profile": "My Profile",
				"topbar|logout":  "Log Out",
			},
		},
		{
			name: "Three-level JSON",
			input: map[string]interface{}{
				"topbar": map[string]interface{}{
					"profile": map[string]interface{}{
						"title": "My Profile",
					},
				},
			},
			expected: map[string]string{
				"topbar|profile|title": "My Profile",
			},
		},
		{
			name: "Mixed Types",
			input: map[string]interface{}{
				"version":      1.0,
				"active":       true,
				"user_message": "Success",
			},
			expected: map[string]string{
				"version":      "1",
				"active":       "true",
				"user_message": "Success",
			},
		},
		{
			name: "Empty Nested Map",
			input: map[string]interface{}{
				"menu": map[string]interface{}{},
			},
			expected: map[string]string{},
		},
		{
			name: "Deep Mixed JSON",
			input: map[string]interface{}{
				"app": map[string]interface{}{
					"home": map[string]interface{}{
						"title":    "Homepage",
						"subtitle": "Welcome back",
					},
					"footer": map[string]interface{}{
						"copyright": "2024",
					},
				},
			},
			expected: map[string]string{
				"app|home|title":       "Homepage",
				"app|home|subtitle":    "Welcome back",
				"app|footer|copyright": "2024",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenJSON(tt.input, "|")
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func createTempJSONFile(t *testing.T, content map[string]interface{}) string {
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.json")

	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	return tmpFile
}

func TestLoadAndFlatten_SimpleNested(t *testing.T) {
	content := map[string]interface{}{
		"greetings": map[string]interface{}{
			"hello": "Hello!",
			"bye":   "Goodbye!",
		},
	}

	filePath := createTempJSONFile(t, content)

	result, err := LoadAndFlatten(filePath)
	if err != nil {
		t.Fatalf("LoadAndFlatten returned error: %v", err)
	}

	expected := map[string]string{
		"greetings.hello": "Hello!",
		"greetings.bye":   "Goodbye!",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestLoadAndFlatten_MultiLevel(t *testing.T) {
	content := map[string]interface{}{
		"menu": map[string]interface{}{
			"file": map[string]interface{}{
				"open":  "Open",
				"close": "Close",
			},
		},
	}

	filePath := createTempJSONFile(t, content)

	result, err := LoadAndFlatten(filePath)
	if err != nil {
		t.Fatalf("LoadAndFlatten returned error: %v", err)
	}

	expected := map[string]string{
		"menu.file.open":  "Open",
		"menu.file.close": "Close",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestLoadAndFlatten_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(invalidPath, []byte(`{invalid json}`), 0644)
	if err != nil {
		t.Fatalf("failed to write invalid JSON file: %v", err)
	}

	_, err = LoadAndFlatten(invalidPath)
	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
	}
}

func TestLoadAndFlatten_FileNotFound(t *testing.T) {
	_, err := LoadAndFlatten("non_existent_file.json")
	if err == nil {
		t.Errorf("expected error for nonexistent file, got nil")
	}
}
