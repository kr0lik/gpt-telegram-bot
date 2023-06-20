package usecase

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/handler/command"
	"gpt-telegran-bot/internal/domain/handler/model"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/editor"
	"gpt-telegran-bot/internal/domain/service/generator"
	"log"
	"strings"
)

const defaultModel = enum.ModelChat

type Messaging struct {
	messenger       service.Messenger
	cache           service.Cache
	commandHandlers []command.Handler
	modelHandlers   []model.Handler
}

func NewMessaging(
	messenger service.Messenger,
	cache service.Cache,
	queue service.Queue,
	speech service.Speech,
	chatGenerator generator.Chat,
	textGenerator generator.Text,
	imageGenerator generator.Image,
	textEditor editor.Text,
	codeEditor editor.Code,
	imageEditor editor.Image,
) *Messaging {
	return &Messaging{
		messenger: messenger,
		cache:     cache,
		commandHandlers: []command.Handler{
			command.NewStart(messenger, cache),
			command.NewHelp(messenger, cache),
			command.NewStatus(messenger, cache, defaultModel),
			command.NewChat(messenger, cache),
			command.NewNew(messenger, chatGenerator),
			command.NewText(messenger, cache),
			command.NewTextEdit(messenger, cache),
			command.NewCodeEdit(messenger, cache),
			command.NewImage(messenger, cache),
			command.NewImageEdit(messenger, cache),
			command.NewSize(messenger, cache, imageGenerator),
			command.NewCount(messenger, cache, imageGenerator),
			command.NewSpeech(messenger, cache),
		},
		modelHandlers: []model.Handler{
			model.NewChat(messenger, chatGenerator, speech, queue),
			model.NewText(messenger, textGenerator, speech),
			model.NewTextEdit(messenger, textEditor),
			model.NewCodeEdit(messenger, codeEditor),
			model.NewImage(messenger, cache, imageGenerator),
			model.NewImageEdit(messenger, cache, imageEditor),
			model.NewSpeech(messenger, cache, speech),
		},
	}
}

func (m *Messaging) Start(ctx context.Context) error {
	updates, err := m.messenger.Listen(ctx)
	if err != nil {
		return fmt.Errorf("failed to get messanger updates: %v", err)
	}

	log.Print("start listen")

	for update := range updates {
		go m.process(update, ctx)
	}

	return nil
}

func (m *Messaging) process(update dto.Income, ctx context.Context) {
	if update.Command != "" {
		m.handleCommand(update)
		return
	}

	m.handleMessage(update, ctx)
}

func (m *Messaging) handleCommand(update dto.Income) {
	commandName := strings.ToLower(update.Command)

	for _, c := range m.commandHandlers {
		if c.Id() == commandName {
			c.Process(update)
			return
		}
	}

	options := m.cache.Get(update.ChatId)

	m.messenger.Replay("Undefined command", update.MessageId, update.ChatId, helper.GetContextCommands(options.Model))
}

func (m *Messaging) handleMessage(update dto.Income, ctx context.Context) {
	options := m.cache.Get(update.ChatId)
	if options.Model == "" {
		options.Model = defaultModel

		m.cache.Set(update.ChatId, options)
	}

	for _, h := range m.modelHandlers {
		if h.Model() == options.Model {
			if update.Callback.Id != "" {
				h.Callback(update, ctx)
				return
			}

			h.Handle(update, ctx)
			return
		}
	}

	m.messenger.Replay("Undefined model", update.MessageId, update.ChatId, helper.GetContextCommands(options.Model))
}
