package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// HotReloadConfig manages hot-reloading of community configurations
type HotReloadConfig struct {
	mu             sync.RWMutex
	watcher        *fsnotify.Watcher
	configDir      string
	currentConfig  map[string]*Community // community name -> config
	reloadCallback func(communityName string, config *Community)
	isWatching     bool
}

// NewHotReloadConfig creates a new hot-reload configuration manager
func NewHotReloadConfig(configDir string) (*HotReloadConfig, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	hrc := &HotReloadConfig{
		watcher:       watcher,
		configDir:     configDir,
		currentConfig: make(map[string]*Community),
	}

	return hrc, nil
}

// SetReloadCallback sets the callback function to be called when a config is reloaded
func (hrc *HotReloadConfig) SetReloadCallback(callback func(communityName string, config *Community)) {
	hrc.mu.Lock()
	defer hrc.mu.Unlock()
	hrc.reloadCallback = callback
}

// StartWatching begins watching the config directory for changes
func (hrc *HotReloadConfig) StartWatching() error {
	hrc.mu.Lock()
	defer hrc.mu.Unlock()

	if hrc.isWatching {
		return nil // Already watching
	}

	// Add config directory to watcher
	err := hrc.watcher.Add(hrc.configDir)
	if err != nil {
		return fmt.Errorf("failed to watch config directory: %w", err)
	}

	hrc.isWatching = true

	// Start processing file system events
	go hrc.processEvents()

	log.Printf("Hot-reload: Started watching config directory: %s", hrc.configDir)
	return nil
}

// StopWatching stops watching for file changes
func (hrc *HotReloadConfig) StopWatching() error {
	hrc.mu.Lock()
	defer hrc.mu.Unlock()

	if !hrc.isWatching {
		return nil
	}

	hrc.isWatching = false
	return hrc.watcher.Close()
}

// LoadConfig loads a community configuration with hot-reload support
func (hrc *HotReloadConfig) LoadConfig(communityName string) (*Community, error) {
	hrc.mu.RLock()
	if config, exists := hrc.currentConfig[communityName]; exists {
		hrc.mu.RUnlock()
		return config, nil
	}
	hrc.mu.RUnlock()

	// Load the config for the first time using the hot-reload config directory
	config, err := hrc.loadConfigFromFile(communityName)
	if err != nil {
		return nil, err
	}

	hrc.mu.Lock()
	hrc.currentConfig[communityName] = config
	hrc.mu.Unlock()

	return config, nil
}

// loadConfigFromFile loads a config file from the hot-reload config directory
func (hrc *HotReloadConfig) loadConfigFromFile(communityName string) (*Community, error) {
	if communityName == "" {
		communityName = "kjernekraft" // Default community
	}

	configPath := filepath.Join(hrc.configDir, communityName+".yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var community Community
	if err := yaml.Unmarshal(data, &community); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	return &community, nil
}

// GetConfig returns the current loaded configuration for a community
func (hrc *HotReloadConfig) GetConfig(communityName string) *Community {
	hrc.mu.RLock()
	defer hrc.mu.RUnlock()
	return hrc.currentConfig[communityName]
}

// processEvents handles file system events and reloads configurations
func (hrc *HotReloadConfig) processEvents() {
	debounceTimer := make(map[string]*time.Timer)
	debounceDuration := 500 * time.Millisecond // Wait 500ms for multiple rapid changes

	for {
		select {
		case event, ok := <-hrc.watcher.Events:
			if !ok {
				return
			}

			// Only process .yaml files
			if filepath.Ext(event.Name) != ".yaml" {
				continue
			}

			// Only process write events (and create for new files)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				filename := filepath.Base(event.Name)
				communityName := filename[:len(filename)-5] // Remove .yaml extension

				log.Printf("Hot-reload: Detected change in %s", filename)

				// Debounce rapid changes
				if timer, exists := debounceTimer[communityName]; exists {
					timer.Stop()
				}

				debounceTimer[communityName] = time.AfterFunc(debounceDuration, func() {
					hrc.reloadConfig(communityName)
					delete(debounceTimer, communityName)
				})
			}

		case err, ok := <-hrc.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Hot-reload: File watcher error: %v", err)
		}
	}
}

// reloadConfig reloads a specific community configuration
func (hrc *HotReloadConfig) reloadConfig(communityName string) {
	log.Printf("Hot-reload: Reloading configuration for community: %s", communityName)

	newConfig, err := hrc.loadConfigFromFile(communityName)
	if err != nil {
		log.Printf("Hot-reload: Failed to reload config for %s: %v", communityName, err)
		return
	}

	hrc.mu.Lock()
	hrc.currentConfig[communityName] = newConfig
	hrc.mu.Unlock()

	// Call the reload callback if set
	hrc.mu.RLock()
	callback := hrc.reloadCallback
	hrc.mu.RUnlock()

	if callback != nil {
		callback(communityName, newConfig)
	}

	log.Printf("Hot-reload: Successfully reloaded configuration for community: %s", communityName)
}

// Global hot-reload manager instance
var globalHotReload *HotReloadConfig

// InitializeHotReload initializes the global hot-reload configuration manager
func InitializeHotReload(configDir string) error {
	var err error
	globalHotReload, err = NewHotReloadConfig(configDir)
	if err != nil {
		return err
	}

	return globalHotReload.StartWatching()
}

// LoadWithHotReload loads a configuration with hot-reload support
func LoadWithHotReload(communityName string) (*Community, error) {
	if globalHotReload == nil {
		// Fallback to regular loading if hot-reload is not initialized
		return Load(communityName)
	}

	return globalHotReload.LoadConfig(communityName)
}

// SetGlobalReloadCallback sets the global reload callback
func SetGlobalReloadCallback(callback func(communityName string, config *Community)) {
	if globalHotReload != nil {
		globalHotReload.SetReloadCallback(callback)
	}
}

// GetCurrentWithHotReload returns the current configuration with hot-reload support
func GetCurrentWithHotReload(communityName string) *Community {
	if globalHotReload == nil {
		return GetCurrent()
	}

	return globalHotReload.GetConfig(communityName)
}

// ShutdownHotReload stops the hot-reload system
func ShutdownHotReload() error {
	if globalHotReload != nil {
		return globalHotReload.StopWatching()
	}
	return nil
}