package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHotReloadConfig(t *testing.T) {
	// Create a temporary directory for test configs
	tempDir := t.TempDir()

	// Create initial test config
	testConfig := `
name: "Test Community"
description: "A test community"
colors:
  primary: "#FF0000"
  secondary: "#00FF00"
  background: "#FFFFFF"
features:
  classes: true
  memberships: true
pricing:
  currency: "USD"
  monthly: 1000
locale:
  language: "en"
  timezone: "UTC"
`

	configPath := filepath.Join(tempDir, "test.yaml")
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Initialize hot-reload config
	hrc, err := NewHotReloadConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to create hot-reload config: %v", err)
	}
	defer hrc.StopWatching()

	// Set up callback to track reloads
	reloadCount := 0
	var lastReloadedConfig *Community
	hrc.SetReloadCallback(func(communityName string, config *Community) {
		reloadCount++
		lastReloadedConfig = config
	})

	// Start watching
	err = hrc.StartWatching()
	if err != nil {
		t.Fatalf("Failed to start watching: %v", err)
	}

	// Load initial config
	config, err := hrc.LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load initial config: %v", err)
	}

	if config.Name != "Test Community" {
		t.Errorf("Expected name 'Test Community', got '%s'", config.Name)
	}

	if config.Colors.Primary != "#FF0000" {
		t.Errorf("Expected primary color '#FF0000', got '%s'", config.Colors.Primary)
	}

	// Update the config file
	updatedConfig := `
name: "Updated Test Community"
description: "An updated test community"
colors:
  primary: "#0000FF"
  secondary: "#00FF00"
  background: "#FFFFFF"
features:
  classes: true
  memberships: true
pricing:
  currency: "USD"
  monthly: 1000
locale:
  language: "en"
  timezone: "UTC"
`

	err = os.WriteFile(configPath, []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	// Wait for the file watcher to detect the change
	// We need to wait a bit longer than the debounce duration
	time.Sleep(1 * time.Second)

	// Check that reload was triggered
	if reloadCount == 0 {
		t.Error("Expected at least one reload, got 0")
	}

	if lastReloadedConfig == nil {
		t.Fatal("No config was reloaded")
	}

	if lastReloadedConfig.Name != "Updated Test Community" {
		t.Errorf("Expected reloaded name 'Updated Test Community', got '%s'", lastReloadedConfig.Name)
	}

	if lastReloadedConfig.Colors.Primary != "#0000FF" {
		t.Errorf("Expected reloaded primary color '#0000FF', got '%s'", lastReloadedConfig.Colors.Primary)
	}

	// Get the updated config from the cache
	cachedConfig := hrc.GetConfig("test")
	if cachedConfig == nil {
		t.Fatal("Cached config should not be nil")
	}

	if cachedConfig.Name != "Updated Test Community" {
		t.Errorf("Expected cached name 'Updated Test Community', got '%s'", cachedConfig.Name)
	}
}

func TestHotReloadMultipleFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create hot-reload config
	hrc, err := NewHotReloadConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to create hot-reload config: %v", err)
	}
	defer hrc.StopWatching()

	reloadCounts := make(map[string]int)
	hrc.SetReloadCallback(func(communityName string, config *Community) {
		reloadCounts[communityName]++
	})

	err = hrc.StartWatching()
	if err != nil {
		t.Fatalf("Failed to start watching: %v", err)
	}

	// Create multiple config files
	configs := map[string]string{
		"community1": `
name: "Community 1"
colors:
  primary: "#FF0000"
`,
		"community2": `
name: "Community 2"
colors:
  primary: "#00FF00"
`,
	}

	for name, content := range configs {
		configPath := filepath.Join(tempDir, name+".yaml")
		err = os.WriteFile(configPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write config for %s: %v", name, err)
		}

		// Load the config to cache it
		_, err = hrc.LoadConfig(name)
		if err != nil {
			t.Fatalf("Failed to load config for %s: %v", name, err)
		}
	}

	// Update one of the files
	updatedContent := `
name: "Updated Community 1"
colors:
  primary: "#0000FF"
`
	configPath := filepath.Join(tempDir, "community1.yaml")
	err = os.WriteFile(configPath, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Wait for reload
	time.Sleep(1 * time.Second)

	// Check that community1 was reloaded
	if reloadCounts["community1"] == 0 {
		t.Error("Expected community1 to be reloaded")
	}

	// Note: Due to file system event timing, community2 might also get triggered
	// during initial file creation. This is acceptable behavior.
	t.Logf("Reload counts: community1=%d, community2=%d", reloadCounts["community1"], reloadCounts["community2"])

	// Verify the updated config
	config := hrc.GetConfig("community1")
	if config == nil {
		t.Fatal("Config should not be nil")
	}

	if config.Name != "Updated Community 1" {
		t.Errorf("Expected name 'Updated Community 1', got '%s'", config.Name)
	}
}

func TestHotReloadInvalidConfig(t *testing.T) {
	tempDir := t.TempDir()

	hrc, err := NewHotReloadConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to create hot-reload config: %v", err)
	}
	defer hrc.StopWatching()

	// Track reload attempts
	var lastReloadError error
	reloadAttempts := 0
	hrc.SetReloadCallback(func(communityName string, config *Community) {
		reloadAttempts++
	})

	err = hrc.StartWatching()
	if err != nil {
		t.Fatalf("Failed to start watching: %v", err)
	}

	// Create valid initial config
	validConfig := `
name: "Valid Community"
colors:
  primary: "#FF0000"
`
	configPath := filepath.Join(tempDir, "test.yaml")
	err = os.WriteFile(configPath, []byte(validConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	// Load initial config
	config, err := hrc.LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load initial config: %v", err)
	}

	if config.Name != "Valid Community" {
		t.Errorf("Expected name 'Valid Community', got '%s'", config.Name)
	}

	// Write invalid YAML
	invalidConfig := `
name: "Invalid Community
colors:
  primary: #FF0000  # Missing quotes
  invalid_yaml_structure
`
	err = os.WriteFile(configPath, []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	// Wait for reload attempt
	time.Sleep(1 * time.Second)

	// The old config should still be available
	cachedConfig := hrc.GetConfig("test")
	if cachedConfig == nil {
		t.Fatal("Cached config should still be available after invalid reload")
	}

	if cachedConfig.Name != "Valid Community" {
		t.Errorf("Expected cached config to remain 'Valid Community', got '%s'", cachedConfig.Name)
	}

	// The reload callback should not have been called for invalid config
	if reloadAttempts > 0 {
		t.Error("Reload callback should not be called for invalid config")
	}

	_ = lastReloadError // Prevent unused variable error
}

func TestGlobalHotReloadFunctions(t *testing.T) {
	tempDir := t.TempDir()

	// Test that global functions work without initialization by falling back to regular loading
	// We need to create the default config directory structure for this test
	originalConfigDir := "config"
	err := os.MkdirAll(originalConfigDir, 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create a minimal default config for fallback testing
	defaultConfig := `
name: "Default Community"
colors:
  primary: "#000000"
`
	defaultConfigPath := filepath.Join(originalConfigDir, "kjernekraft.yaml")
	err = os.WriteFile(defaultConfigPath, []byte(defaultConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write default config: %v", err)
	}
	defer os.Remove(defaultConfigPath) // Clean up after test

	// Test loading nonexistent config (should fail)
	_, err = LoadWithHotReload("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent config")
	}

	// Initialize global hot-reload
	err = InitializeHotReload(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize global hot-reload: %v", err)
	}
	defer ShutdownHotReload()

	// Create a test config in the hot-reload directory
	testConfig := `
name: "Global Test Community"
colors:
  primary: "#123456"
`
	configPath := filepath.Join(tempDir, "global.yaml")
	err = os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test global functions
	config, err := LoadWithHotReload("global")
	if err != nil {
		t.Fatalf("Failed to load config with global hot-reload: %v", err)
	}

	if config.Name != "Global Test Community" {
		t.Errorf("Expected name 'Global Test Community', got '%s'", config.Name)
	}

	// Test callback setting
	callbackCalled := false
	SetGlobalReloadCallback(func(name string, cfg *Community) {
		callbackCalled = true
	})

	// Update config to trigger callback
	updatedConfig := `
name: "Updated Global Test Community"
colors:
  primary: "#654321"
`
	err = os.WriteFile(configPath, []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Wait for callback
	time.Sleep(1 * time.Second)

	if !callbackCalled {
		t.Error("Expected global callback to be called")
	}
}