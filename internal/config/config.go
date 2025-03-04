package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home directory not found: %w", err)
	}

	file, err := os.ReadFile(homeDir + ".gatorconfig.json")
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	var userConfig Config

	if err = json.Unmarshal(file, &userConfig); err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %w", err)
	}
	return &userConfig, nil
}

func (c *Config) SetUser() {

}
