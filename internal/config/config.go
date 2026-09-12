package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	UserName string `json:"user_name"`
	DbUrl    string `json:"db_url"`
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := dir + "/" + configFileName
	return path, err
}

func write(cfg Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.UserName = userName
	err := write(*cfg)
	return err
}
