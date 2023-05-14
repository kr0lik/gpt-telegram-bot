package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Text struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewText(messenger service.Messenger, cache service.Cache) *Text {
	return &Text{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *Text) Id() string {
	return enum.CommandText
}

func (c *Text) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelText
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Text generation model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
