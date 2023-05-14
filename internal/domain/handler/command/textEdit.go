package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type TextEdit struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewTextEdit(messenger service.Messenger, cache service.Cache) *TextEdit {
	return &TextEdit{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *TextEdit) Id() string {
	return enum.CommandTextEdit
}

func (c *TextEdit) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelTextEdit
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Edit text model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
