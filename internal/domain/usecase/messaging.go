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
	commandHandlers []command.Handler,
	modelHandlers []model.Handler,
) *Messaging {
	return &Messaging{
		messenger:       messenger,
		cache:           cache,
		commandHandlers: commandHandlers,
		modelHandlers:   modelHandlers,
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
	options := m.cache.Get(update.ChatId)
	if options.Model == "" {
		options.Model = defaultModel

		m.cache.Set(update.ChatId, options)
	}

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
