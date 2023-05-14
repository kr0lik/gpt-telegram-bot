package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Chat struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewChat(messenger service.Messenger, cache service.Cache) *Chat {
	return &Chat{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *Chat) Id() string {
	return enum.CommandChat
}

func (c *Chat) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelChat
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Chat conversation model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
