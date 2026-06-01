package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gerlv/trello-cli/internal/config"
)

func TestDefaults(t *testing.T) {
	cfg := config.Load()

	if cfg.Profile != "default" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "default")
	}
	if cfg.Pretty != false {
		t.Errorf("Pretty = %v, want false", cfg.Pretty)
	}
	if cfg.Verbose != false {
		t.Errorf("Verbose = %v, want false", cfg.Verbose)
	}
	if cfg.Timeout != 15*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 15*time.Second)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
	}
	if cfg.RetryMutations != false {
		t.Errorf("RetryMutations = %v, want false", cfg.RetryMutations)
	}
}

func TestEnvOverridesDefaults(t *testing.T) {
	t.Setenv("TRELLO_PROFILE", "work")
	t.Setenv("TRELLO_PRETTY", "true")
	t.Setenv("TRELLO_TIMEOUT", "30s")
	t.Setenv("TRELLO_MAX_RETRIES", "5")
	t.Setenv("TRELLO_RETRY_MUTATIONS", "true")
	t.Setenv("TRELLO_VERBOSE", "true")

	cfg := config.Load()

	if cfg.Profile != "work" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "work")
	}
	if cfg.Pretty != true {
		t.Errorf("Pretty = %v, want true", cfg.Pretty)
	}
	if cfg.Verbose != true {
		t.Errorf("Verbose = %v, want true", cfg.Verbose)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 30*time.Second)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", cfg.MaxRetries)
	}
	if cfg.RetryMutations != true {
		t.Errorf("RetryMutations = %v, want true", cfg.RetryMutations)
	}
}

func TestMissingConfigFileDoesNotError(t *testing.T) {
	// Ensure no config file exists at the default path
	_ = os.Remove("/tmp/nonexistent-trello-config-test.yaml")
	t.Setenv("TRELLO_CONFIG_PATH", "/tmp/nonexistent-trello-config-test.yaml")

	cfg := config.Load()
	// Should still return defaults without error
	if cfg.Profile != "default" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "default")
	}
}

func TestConfigFileOverridesDefaults(t *testing.T) {
	// Create a temp config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := []byte("profile: personal\npretty: true\ntimeout: 20s\nmax_retries: 1\n")
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	t.Setenv("TRELLO_CONFIG_PATH", configPath)
	for _, key := range []string{
		"TRELLO_PROFILE",
		"TRELLO_PRETTY",
		"TRELLO_TIMEOUT",
		"TRELLO_MAX_RETRIES",
		"TRELLO_RETRY_MUTATIONS",
		"TRELLO_VERBOSE",
	} {
		t.Setenv(key, "")
		_ = os.Unsetenv(key)
	}

	cfg := config.Load()

	if cfg.Profile != "personal" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "personal")
	}
	if cfg.Pretty != true {
		t.Errorf("Pretty = %v, want true", cfg.Pretty)
	}
	if cfg.Timeout != 20*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 20*time.Second)
	}
	if cfg.MaxRetries != 1 {
		t.Errorf("MaxRetries = %d, want 1", cfg.MaxRetries)
	}
}

func TestEnvOverridesConfigFile(t *testing.T) {
	// Create a config file with profile=personal
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := []byte("profile: personal\n")
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	t.Setenv("TRELLO_CONFIG_PATH", configPath)
	for _, key := range []string{
		"TRELLO_PROFILE",
		"TRELLO_PRETTY",
		"TRELLO_TIMEOUT",
		"TRELLO_MAX_RETRIES",
		"TRELLO_RETRY_MUTATIONS",
		"TRELLO_VERBOSE",
	} {
		t.Setenv(key, "")
		_ = os.Unsetenv(key)
	}
	t.Setenv("TRELLO_PROFILE", "work")

	cfg := config.Load()

	// Env should override config file
	if cfg.Profile != "work" {
		t.Errorf("Profile = %q, want %q (env should override config file)", cfg.Profile, "work")
	}
}
