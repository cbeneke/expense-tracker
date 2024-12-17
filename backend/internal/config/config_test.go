package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default values
	config := Load()
	if config.Port != "8080" {
		t.Errorf("Default port should be 8080, got %s", config.Port)
	}

	// Test environment variable override
	os.Setenv("PORT", "3000")
	config = Load()
	if config.Port != "3000" {
		t.Errorf("Port should be 3000, got %s", config.Port)
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		defValue string
		envValue string
		expected string
	}{
		{
			name:     "returns default when env not set",
			key:      "TEST_KEY_1",
			defValue: "default",
			envValue: "",
			expected: "default",
		},
		{
			name:     "returns env value when set",
			key:      "TEST_KEY_2",
			defValue: "default",
			envValue: "custom",
			expected: "custom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				os.Setenv(tc.key, tc.envValue)
				defer os.Unsetenv(tc.key)
			}

			result := getEnvWithDefault(tc.key, tc.defValue)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
