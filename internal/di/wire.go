//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gpt-telegran-bot/internal/di/config"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/editor"
	"gpt-telegran-bot/internal/domain/service/generator"
	"gpt-telegran-bot/internal/domain/usecase"
	openAiClient "gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/service/cache"
	openAiEditor "gpt-telegran-bot/internal/infrastructure/service/editor/openAi"
	openAiGenerator "gpt-telegran-bot/internal/infrastructure/service/generator/openAi"
	"gpt-telegran-bot/internal/infrastructure/service/messenger"
	openAiQueue "gpt-telegran-bot/internal/infrastructure/service/queue"
	openAiSpeech "gpt-telegran-bot/internal/infrastructure/service/speech/openAi"
)

var cacheSet = wire.NewSet(
	cache.NewMemory,
	wire.Bind(new(service.Cache), new(*cache.Memory)),
)

var queueSet = wire.NewSet(
	openAiQueue.NewOpenAi,
	wire.Bind(new(service.Queue), new(*openAiQueue.OpenAi)),
)

var messengerSet = wire.NewSet(
	config.ProvideTelegramBotConfig,
	messenger.NewTelegram,
	wire.Bind(new(service.Messenger), new(*messenger.Telegram)),
)

var openAiSet = wire.NewSet(
	config.ProvideOpenAiClientConfig,
	openAiClient.NewClient,
	// generators
	openAiGenerator.NewChat,
	wire.Bind(new(generator.Chat), new(*openAiGenerator.Chat)),
	openAiGenerator.NewText,
	wire.Bind(new(generator.Text), new(*openAiGenerator.Text)),
	openAiGenerator.NewImage,
	wire.Bind(new(generator.Image), new(*openAiGenerator.Image)),
	// editors
	openAiEditor.NewText,
	wire.Bind(new(editor.Text), new(*openAiEditor.Text)),
	openAiEditor.NewCode,
	wire.Bind(new(editor.Code), new(*openAiEditor.Code)),
	openAiEditor.NewImage,
	wire.Bind(new(editor.Image), new(*openAiEditor.Image)),
	// speech
	openAiSpeech.NewSpeech,
	wire.Bind(new(service.Speech), new(*openAiSpeech.Speech)),
)

func InitialiseMessaging() (*usecase.Messaging, error) {
	wire.Build(
		cacheSet,
		queueSet,
		messengerSet,
		openAiSet,
		usecase.NewMessaging,
	)
	return &usecase.Messaging{}, nil
}
