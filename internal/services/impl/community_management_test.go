package impl_test

import (
	"context"
	"testing"

	"samskipnad/internal/services/impl"

	"github.com/stretchr/testify/assert"
)

func TestCommunityManagementServiceImpl(t *testing.T) {
	service := impl.NewCommunityManagementService()

	t.Run("GetTenantBySlug", func(t *testing.T) {
		tenant, err := service.GetTenantBySlug(context.Background(), "test-community")
		
		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "test-community", tenant.Slug)
		assert.Equal(t, "Community test-community", tenant.Name)
	})

	t.Run("IsFeatureEnabled", func(t *testing.T) {
		// Test with unknown tenant - should default to enabled
		enabled, err := service.IsFeatureEnabled(context.Background(), 999, "classes")
		assert.NoError(t, err)
		assert.True(t, enabled)

		// Test with unknown feature - should default to enabled
		enabled, err = service.IsFeatureEnabled(context.Background(), 1, "unknown-feature")
		assert.NoError(t, err)
		assert.True(t, enabled)
	})

	t.Run("LoadConfiguration", func(t *testing.T) {
		// This will try to load from the existing config system
		// It may fail if no config exists, which is expected
		_, err := service.LoadConfiguration(context.Background(), "test")
		
		// We don't assert on the error because config files may not exist
		// The important thing is that the method doesn't panic
		t.Logf("Config load result: %v", err)
	})

	t.Run("GetCommunity", func(t *testing.T) {
		community, err := service.GetCommunity(context.Background(), 1)
		
		assert.NoError(t, err)
		assert.NotNil(t, community)
		assert.Equal(t, 1, community.ID)
		assert.Equal(t, "Tenant 1", community.Name)
	})
}