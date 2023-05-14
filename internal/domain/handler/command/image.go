package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Image struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewImage(messenger service.Messenger, cache service.Cache) *Image {
	return &Image{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *Image) Id() string {
	return enum.CommandImage
}

func (c *Image) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelImage
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Image generation model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
