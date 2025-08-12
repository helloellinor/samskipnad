package impl

import (
	"context"
	"testing"
)

func TestPluginHostService_Creation(t *testing.T) {
	service := NewPluginHostService()
	if service == nil {
		t.Fatal("Expected plugin host service to be created, got nil")
	}
}

func TestPluginHostService_Initialize(t *testing.T) {
	service := NewPluginHostService()
	ctx := context.Background()

	err := service.Initialize(ctx)
	if err != nil {
		t.Errorf("Expected no error during initialization, got: %v", err)
	}

	// Clean shutdown
	err = service.Shutdown(ctx)
	if err != nil {
		t.Errorf("Expected no error during shutdown, got: %v", err)
	}
}

func TestPluginHostService_GetLoadedPlugins(t *testing.T) {
	service := NewPluginHostService()
	ctx := context.Background()

	// Initially no plugins should be loaded
	plugins := service.GetLoadedPlugins(ctx)
	if len(plugins) != 0 {
		t.Errorf("Expected 0 loaded plugins, got %d", len(plugins))
	}

	// Clean shutdown
	service.Shutdown(ctx)
}

func TestPluginHostService_LoadNonExistentPlugin(t *testing.T) {
	service := NewPluginHostService()
	ctx := context.Background()

	defer service.Shutdown(ctx)

	err := service.LoadPlugin(ctx, "test-plugin", "./non-existent-plugin")
	if err == nil {
		t.Error("Expected error when loading non-existent plugin, got nil")
	}
}

func TestPluginHostService_ExecuteNonExistentPlugin(t *testing.T) {
	service := NewPluginHostService()
	ctx := context.Background()

	defer service.Shutdown(ctx)

	result, err := service.ExecutePlugin(ctx, "non-existent-plugin", map[string]interface{}{
		"test": "value",
	})

	if err == nil {
		t.Error("Expected error when executing non-existent plugin, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result for non-existent plugin, got: %v", result)
	}
}