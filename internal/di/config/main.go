package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var main *Config

type Config struct {
	TelegramToken            string   `yaml:"telegramToken"`
	TelegramAllowedUsernames []string `yaml:"telegramAllowedUsernames"`
	OpenAIKey                string   `yaml:"openAiApiKey"`
	FileDownloadPath         string   `yaml:"fileDownloadPath"`
}

func ReadConfig(configPath string) error {
	main = new(Config)

	configFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}

	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(main); err != nil {
		return fmt.Errorf("failed to decode config file: %v", err)
	}

	return nil
}
