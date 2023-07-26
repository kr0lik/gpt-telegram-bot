package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"strconv"
)

type Status struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewStatus(messenger service.Messenger, cache service.Cache) *Status {
	return &Status{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *Status) Id() string {
	return enum.CommandStatus
}

func (c *Status) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)

	text := ""

	if opts.Model != "" {
		text = "Current model: " + opts.Model
	}

	if opts.Image.Size != "" {
		text += "\nImage size: " + opts.Image.Size
	}

	if opts.Image.Count > 0 {
		text += "\nImage count: " + strconv.Itoa(opts.Image.Count)
	}

	c.messenger.Send(text, update.ChatId, helper.GetContextCommands(opts.Model))
}
