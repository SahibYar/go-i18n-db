package i18n

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformMapToJSON(t *testing.T) {
	// Test 1: Single-level key with value and tooltip
	input := map[string]map[string]string{
		"topbar.profile": {
			"value":   "Profile",
			"tooltip": "Your profile",
		},
	}
	expected := `{
  "topbar": {
    "profile": "Profile",
    "profile_tooltip": "Your profile"
  }
}`
	jsonResult, err := transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 2: Multiple keys with different levels of nesting
	input = map[string]map[string]string{
		"abc.def.topbar.profile": {
			"value":   "Profile",
			"tooltip": "Your profile",
		},
		"mno.xyz.footer.contact": {
			"value":   "Contact",
			"tooltip": "Contact us",
		},
	}
	expected = `{
  "abc": {
    "def": {
      "topbar": {
        "profile": "Profile",
        "profile_tooltip": "Your profile"
      }
    }
  },
  "mno": {
    "xyz": {
      "footer": {
        "contact": "Contact",
        "contact_tooltip": "Contact us"
      }
    }
  }
}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 3: Single key with multiple parts in the key
	input = map[string]map[string]string{
		"xyz.abc.topbar.contact": {
			"value":   "Contact",
			"tooltip": "Contact us",
		},
	}
	expected = `{
  "xyz": {
    "abc": {
      "topbar": {
        "contact": "Contact",
        "contact_tooltip": "Contact us"
      }
    }
  }
}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 4: Key with no value or tooltip
	input = map[string]map[string]string{
		"header.title": {
			"value":   "Title",
			"tooltip": "",
		},
	}
	expected = `{
  "header": {
    "title": "Title",
    "title_tooltip": ""
  }
}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 5: Key with only value, no tooltip
	input = map[string]map[string]string{
		"footer.about": {
			"value": "About Us",
		},
	}
	expected = `{
  "footer": {
    "about": "About Us"
  }
}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 6: Empty input map
	input = map[string]map[string]string{}
	expected = `{}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 7: Key with more than 2 levels
	input = map[string]map[string]string{
		"first.second.third.fourth.value": {
			"value":   "Value at fourth",
			"tooltip": "Tooltip at fourth",
		},
	}
	expected = `{
  "first": {
    "second": {
      "third": {
        "fourth": {
          "value": "Value at fourth",
          "value_tooltip": "Tooltip at fourth"
        }
      }
    }
  }
}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)

	// Test 8: Multiple keys with more than 2 levels
	input = map[string]map[string]string{
		"first.second.third.fourth.value": {
			"value":   "Value at fourth",
			"tooltip": "Tooltip at fourth",
		},
		"first.second.third.fourth.sahib": {
			"value":   "Value at fourth",
			"tooltip": "Tooltip at fourth",
		},
	}
	expected = `{
  "first": {
    "second": {
      "third": {
        "fourth": {
          "value": "Value at fourth",
          "value_tooltip": "Tooltip at fourth",
          "sahib": "Value at fourth",
          "sahib_tooltip": "Tooltip at fourth"
        }
      }
    }
  }
}`
	jsonResult, err = transformMapToJSON(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonResult)
}
