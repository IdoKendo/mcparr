package config

import (
	"os"
	"testing"
)

func TestEnvWithDefault(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test-value")
	value := envWithDefault("TEST_ENV_VAR", "default-value")
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", value)
	}

	os.Unsetenv("TEST_ENV_VAR")
	value = envWithDefault("TEST_ENV_VAR", "default-value")
	if value != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", value)
	}
}

func TestEnvIntWithDefault(t *testing.T) {
	os.Setenv("TEST_INT_VAR", "42")
	value := envIntWithDefault("TEST_INT_VAR", 10)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	os.Setenv("TEST_INT_VAR", "not-an-int")
	value = envIntWithDefault("TEST_INT_VAR", 10)
	if value != 10 {
		t.Errorf("Expected 10 for invalid input, got %d", value)
	}

	os.Unsetenv("TEST_INT_VAR")
	value = envIntWithDefault("TEST_INT_VAR", 10)
	if value != 10 {
		t.Errorf("Expected 10 for unset variable, got %d", value)
	}
}
