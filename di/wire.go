//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gpt-telegran-bot/internal/config"
	"gpt-telegran-bot/internal/domain/handler/command"
	"gpt-telegran-bot/internal/domain/handler/model"
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

var commandSet = wire.NewSet(
	command.NewStart,
	command.NewHelp,
	command.NewStatus,
	command.NewChat,
	command.NewNew,
	command.NewText,
	command.NewTextEdit,
	command.NewCodeEdit,
	command.NewImage,
	command.NewImageEdit,
	command.NewSize,
	command.NewCount,
	command.NewSpeech,
)

var modelSet = wire.NewSet(
	model.NewChat,
	model.NewText,
	model.NewTextEdit,
	model.NewCodeEdit,
	model.NewImage,
	model.NewImageEdit,
	model.NewSpeech,
)

func provideOpenAiCommandList(
	start *command.Start,
	help *command.Help,
	status *command.Status,
	chat *command.Chat,
	new *command.New,
	text *command.Text,
	textEdit *command.TextEdit,
	codeEdit *command.CodeEdit,
	image *command.Image,
	imageEdit *command.ImageEdit,
	size *command.Size,
	count *command.Count,
	speech *command.Speech,
) []command.Handler {
	return []command.Handler{
		start,
		help,
		status,
		chat,
		new,
		text,
		textEdit,
		codeEdit,
		image,
		imageEdit,
		size,
		count,
		speech,
	}
}

func provideOpenAiModelList(
	chat *model.Chat,
	text *model.Text,
	textEdit *model.TextEdit,
	codeEdit *model.CodeEdit,
	image *model.Image,
	imageEdit *model.ImageEdit,
	speech *model.Speech,
) []model.Handler {
	return []model.Handler{
		chat,
		text,
		textEdit,
		codeEdit,
		image,
		imageEdit,
		speech,
	}
}

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
	// Commands
	provideOpenAiCommandList,
	// Models
	provideOpenAiModelList,
)

func InitialiseOpenAiMessaging() (*usecase.Messaging, error) {
	wire.Build(
		cacheSet,
		queueSet,
		messengerSet,
		commandSet,
		modelSet,
		openAiSet,
		usecase.NewMessaging,
	)
	return &usecase.Messaging{}, nil
}
