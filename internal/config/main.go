package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/service/messenger"
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

func ProvideTelegramBotConfig() *messenger.TelegramConfig {
	return &messenger.TelegramConfig{
		ApiToken:     main.TelegramToken,
		DownloadPath: main.FileDownloadPath,
		AllowedUsers: main.TelegramAllowedUsernames,
	}
}

func ProvideOpenAiClientConfig() *openAi.ClientConfig {
	return &openAi.ClientConfig{
		ApiKey: main.OpenAIKey,
	}
}
