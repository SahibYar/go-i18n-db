package i18n

import (
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
