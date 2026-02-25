package loglinter

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the linter configuration
type Config struct {
	// Rules configuration
	CheckLowercase      bool     `json:"check_lowercase"`
	CheckEnglishOnly    bool     `json:"check_english_only"`
	CheckSpecialChars   bool     `json:"check_special_chars"`
	CheckSensitiveData  bool     `json:"check_sensitive_data"`
	
	// Custom sensitive keywords
	SensitiveKeywords   []string `json:"sensitive_keywords"`
	
	// Allowed special characters
	AllowedPunctuation  string   `json:"allowed_punctuation"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		CheckLowercase:     true,
		CheckEnglishOnly:   true,
		CheckSpecialChars:  true,
		CheckSensitiveData: true,
		SensitiveKeywords: []string{
			"password", "passwd", "pwd",
			"token", "api_key", "apikey", "api-key",
			"secret", "private_key", "privatekey",
			"credential",
		},
		AllowedPunctuation: ".,;:!?",
	}
}

// LoadConfig loads configuration from .loglinter.json file
func LoadConfig(dir string) (*Config, error) {
	configPath := filepath.Join(dir, ".loglinter.json")
	
	// If config file doesn't exist, use default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	
	return config, nil
}
