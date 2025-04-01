package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigSearchPaths defines the locations where config.yml will be searched
var ConfigSearchPaths = []string{
	".",               // Current directory
	"$HOME/.dbbackup", // User's home directory
	"/etc/dbbackup",   // System-wide configuration
	filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dbbackup"), // XDG config directory
}

// findConfigFile searches for config.yml in predefined locations
func findConfigFile() (string, error) {
	// First, check if XDG_CONFIG_HOME is set
	if os.Getenv("XDG_CONFIG_HOME") == "" {
		os.Setenv("XDG_CONFIG_HOME", filepath.Join(os.Getenv("HOME"), ".config"))
	}

	// Expand all paths
	for i, path := range ConfigSearchPaths {
		expanded := os.ExpandEnv(path)
		ConfigSearchPaths[i] = expanded
	}

	// Search for config.yml in all locations
	for _, path := range ConfigSearchPaths {
		configPath := filepath.Join(path, "config.yml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	return "", fmt.Errorf("config.yml not found in any of the following locations: %v", ConfigSearchPaths)
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type StorageConfig struct {
	Type      string `yaml:"type"` // local, s3, gcs, azure
	Enabled   bool   `yaml:"enabled"`
	Path      string `yaml:"path"` // for local storage
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Region    string `yaml:"region"`
}

type NotificationConfig struct {
	SlackWebhook string `yaml:"slack_webhook"`
	Enabled      bool   `yaml:"enabled"`
}

type Config struct {
	Database     DatabaseConfig     `yaml:"database"`
	Storage      StorageConfig      `yaml:"storage"`
	Notification NotificationConfig `yaml:"notification"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(path string) (*Config, error) {
	// If no path is provided, search for config.yml
	if path == "" {
		var err error
		path, err = findConfigFile()
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
