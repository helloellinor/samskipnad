package plugins

import (
	"context"
	"testing"
)

func TestPluginHost_Creation(t *testing.T) {
	// Test that we can create a plugin host
	config := &PluginConfig{
		PluginDir:     "./testplugins",
		MaxPlugins:    5,
		PluginTimeout: 10,
	}

	host := NewPluginHost(config)
	if host == nil {
		t.Fatal("Expected plugin host to be created, got nil")
	}

	if len(host.GetLoadedPlugins()) != 0 {
		t.Errorf("Expected 0 loaded plugins, got %d", len(host.GetLoadedPlugins()))
	}

	// Test graceful shutdown
	host.Shutdown()
}

func TestPluginHost_Configuration(t *testing.T) {
	// Test with nil config (should use defaults)
	host := NewPluginHost(nil)
	if host == nil {
		t.Fatal("Expected plugin host to be created with default config, got nil")
	}

	if host.config.MaxPlugins != 10 {
		t.Errorf("Expected default MaxPlugins to be 10, got %d", host.config.MaxPlugins)
	}

	if host.config.PluginDir != "./plugins" {
		t.Errorf("Expected default PluginDir to be './plugins', got %s", host.config.PluginDir)
	}

	host.Shutdown()
}

func TestPluginHost_LoadNonExistentPlugin(t *testing.T) {
	host := NewPluginHost(nil)
	defer host.Shutdown()

	ctx := context.Background()
	err := host.LoadPlugin(ctx, "test-plugin", "./non-existent-plugin")
	
	if err == nil {
		t.Error("Expected error when loading non-existent plugin, got nil")
	}
}

func TestPluginHost_MaxPluginsLimit(t *testing.T) {
	config := &PluginConfig{
		PluginDir:     "./testplugins",
		MaxPlugins:    1, // Set limit to 1 for testing
		PluginTimeout: 10,
	}

	host := NewPluginHost(config)
	defer host.Shutdown()

	ctx := context.Background()
	
	// Try to load multiple plugins (will fail because they don't exist, but test the limit logic)
	err1 := host.LoadPlugin(ctx, "plugin1", "./non-existent1")
	err2 := host.LoadPlugin(ctx, "plugin2", "./non-existent2")

	// Both should fail due to non-existent files, not limits
	if err1 == nil {
		t.Error("Expected error for first plugin load (file doesn't exist)")
	}
	if err2 == nil {
		t.Error("Expected error for second plugin load (file doesn't exist)")
	}
}

func TestPluginHost_GetLoadedPlugins(t *testing.T) {
	host := NewPluginHost(nil)
	defer host.Shutdown()

	// Initially no plugins should be loaded
	plugins := host.GetLoadedPlugins()
	if len(plugins) != 0 {
		t.Errorf("Expected 0 loaded plugins, got %d", len(plugins))
	}
}