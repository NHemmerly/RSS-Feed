package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(c)
}

func Read() (*Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("could not build filepath: %w", err)
	}
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	var userConfig Config

	if err = json.Unmarshal(file, &userConfig); err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %w", err)
	}
	return &userConfig, nil
}

func write(cfg *Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("could not build filepath: %w", err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("could not encode struct: %w", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory not found: %w", err)
	}
	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}
