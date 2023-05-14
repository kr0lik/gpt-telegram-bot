package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type ImageEdit struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewImageEdit(messenger service.Messenger, cache service.Cache) *ImageEdit {
	return &ImageEdit{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *ImageEdit) Id() string {
	return enum.CommandImageEdit
}

func (c *ImageEdit) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelImageEdit
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Edit image model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
