package config

import (
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/service/messenger"
)

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
