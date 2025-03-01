package config

import (
	"errors"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Port int

	// Registry configuration
	RegistryURL       string
	RegistryUsername  string
	RegistryPassword  string
	RegistryInsecure  bool
	RegistryAPIPrefix string

	// UI configuration
	DefaultTheme string
}

// NewConfig creates a new configuration instance from environment variables
func NewConfig() (*Config, error) {
	// Default values
	config := &Config{
		Port:              8080,
		RegistryAPIPrefix: "/v2",
		DefaultTheme:      "light",
	}

	// Server configuration
	if port := os.Getenv("PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, errors.New("invalid PORT value")
		}
		config.Port = p
	}

	// Registry configuration
	if url := os.Getenv("REGISTRY_URL"); url != "" {
		config.RegistryURL = url
	} else {
		return nil, errors.New("REGISTRY_URL is required")
	}

	config.RegistryUsername = os.Getenv("REGISTRY_USERNAME")
	config.RegistryPassword = os.Getenv("REGISTRY_PASSWORD")

	if insecure := os.Getenv("REGISTRY_INSECURE"); insecure == "true" {
		config.RegistryInsecure = true
	}

	if apiPrefix := os.Getenv("REGISTRY_API_PREFIX"); apiPrefix != "" {
		config.RegistryAPIPrefix = apiPrefix
	}

	// UI configuration
	if theme := os.Getenv("DEFAULT_THEME"); theme == "dark" || theme == "light" {
		config.DefaultTheme = theme
	}

	return config, nil
}

// UpdateConfig updates the configuration with new values
func (c *Config) UpdateConfig(
	registryURL string,
	registryUsername string,
	registryPassword string,
	registryInsecure bool,
	defaultTheme string,
) error {
	if registryURL == "" {
		return errors.New("registry URL cannot be empty")
	}

	c.RegistryURL = registryURL
	c.RegistryUsername = registryUsername
	c.RegistryPassword = registryPassword
	c.RegistryInsecure = registryInsecure

	if defaultTheme == "dark" || defaultTheme == "light" {
		c.DefaultTheme = defaultTheme
	}

	return nil
}
