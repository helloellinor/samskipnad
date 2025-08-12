package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CardPackage represents a standardized pricing package
type CardPackage struct {
	Name          string `yaml:"name"`
	Klipp         int    `yaml:"klipp"`           // Number of klipp/credits
	Price         int    `yaml:"price"`           // Total price in cents
	PricePerKlipp int    `yaml:"price_per_klipp"` // Price per klipp for display
	SavePercent   int    `yaml:"save_percent"`    // Percentage savings
	Badge         string `yaml:"badge"`           // Optional badge text (e.g., "Best Deal")
	Description   string `yaml:"description"`     // Optional description
}

// KlippekortCategory represents a category of classes with klippekort pricing
type KlippekortCategory struct {
	ID          string        `yaml:"id"`
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Icon        string        `yaml:"icon"`
	Color       string        `yaml:"color"`
	Packages    []CardPackage `yaml:"packages"`
	InfoText    string        `yaml:"info_text"` // Trust & action microcopy
}

// Community represents the configuration for a community
type Community struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Tagline     string `yaml:"tagline"`

	Colors struct {
		Primary    string `yaml:"primary"`
		Secondary  string `yaml:"secondary"`
		Accent     string `yaml:"accent"`
		Success    string `yaml:"success"`
		Warning    string `yaml:"warning"`
		Danger     string `yaml:"danger"`
		Background string `yaml:"background"`
		Surface    string `yaml:"surface"`
		Text       string `yaml:"text"`
		Muted      string `yaml:"muted"`
	} `yaml:"colors"`

	Fonts struct {
		Primary   string `yaml:"primary"`
		Secondary string `yaml:"secondary"`
		SizeBase  string `yaml:"size_base"`
	} `yaml:"fonts"`

	Features struct {
		Classes     bool `yaml:"classes"`
		Memberships bool `yaml:"memberships"`
		Community   bool `yaml:"community"`
		Payments    bool `yaml:"payments"`
		Calendar    bool `yaml:"calendar"`
	} `yaml:"features"`

	Content struct {
		Home struct {
			Title       string `yaml:"title"`
			Subtitle    string `yaml:"subtitle"`
			Description string `yaml:"description"`
		} `yaml:"home"`

		Features []struct {
			Title       string `yaml:"title"`
			Description string `yaml:"description"`
		} `yaml:"features"`
	} `yaml:"content"`

	Pricing struct {
		Currency string `yaml:"currency"`
		Monthly  int    `yaml:"monthly"`
		Yearly   int    `yaml:"yearly"`
		DropIn   int    `yaml:"drop_in"`

		// Klippekort pricing by category
		Klippekort struct {
			Categories []KlippekortCategory `yaml:"categories"`
		} `yaml:"klippekort"`
	} `yaml:"pricing"`

	Locale struct {
		Language string `yaml:"language"`
		Country  string `yaml:"country"`
		Timezone string `yaml:"timezone"`
	} `yaml:"locale"`

	Admin struct {
		RegistrationOpen bool   `yaml:"registration_open"`
		RequireApproval  bool   `yaml:"require_approval"`
		DefaultRole      string `yaml:"default_role"`
	} `yaml:"admin"`

	Attribution struct {
		Show bool   `yaml:"show"`
		Text string `yaml:"text"`
		Link string `yaml:"link"`
	} `yaml:"attribution"`
}

var currentCommunity *Community

// Load loads the community configuration from a YAML file
func Load(communityName string) (*Community, error) {
	if communityName == "" {
		communityName = "kjernekraft" // Default community
	}

	configPath := filepath.Join("config", communityName+".yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var community Community
	if err := yaml.Unmarshal(data, &community); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	currentCommunity = &community
	return &community, nil
}

// GetCurrent returns the currently loaded community configuration
func GetCurrent() *Community {
	if currentCommunity == nil {
		// Load default if none is loaded
		community, err := Load("")
		if err != nil {
			panic("Failed to load default community config: " + err.Error())
		}
		return community
	}
	return currentCommunity
}

// FormatPrice formats a price with the community's currency
func (c *Community) FormatPrice(amount int) string {
	return fmt.Sprintf("%d %s", amount, c.Pricing.Currency)
}

// CalculateSavings calculates the savings for a card package
func (c *Community) CalculateSavings(category KlippekortCategory, packageIndex int) int {
	if packageIndex >= len(category.Packages) || packageIndex < 1 {
		return 0
	}

	pkg := category.Packages[packageIndex]
	singlePrice := category.Packages[0] // Assume first package is single klipp

	if singlePrice.Klipp == 1 {
		totalSinglePrice := singlePrice.Price * pkg.Klipp
		savings := totalSinglePrice - pkg.Price
		return savings
	}

	return 0
}

// GetBestValuePackage returns the index of the best value package in a category
func (c *Community) GetBestValuePackage(category KlippekortCategory) int {
	if len(category.Packages) <= 1 {
		return 0
	}

	bestIndex := 0
	bestPricePerKlipp := category.Packages[0].PricePerKlipp

	for i, pkg := range category.Packages {
		if pkg.PricePerKlipp < bestPricePerKlipp {
			bestPricePerKlipp = pkg.PricePerKlipp
			bestIndex = i
		}
	}

	return bestIndex
}
