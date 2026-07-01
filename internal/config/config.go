package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return write(*c)
}

func Read() (Config, error) {
	filename, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting filepath: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config file: %w", err)
	}
	return cfg, nil
}

func write(cfg Config) error {
	filename, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting filepath: %w", err)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling config file: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	filename := filepath.Join(home, configFileName)
	return filename, nil
}
