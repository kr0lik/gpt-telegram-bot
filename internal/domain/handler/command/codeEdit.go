package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type CodeEdit struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewCodeEdit(messenger service.Messenger, cache service.Cache) *CodeEdit {
	return &CodeEdit{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *CodeEdit) Id() string {
	return enum.CommandCodeEdit
}

func (c *CodeEdit) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelCodeEdit
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Edit code model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
