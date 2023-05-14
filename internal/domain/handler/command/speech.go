package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Speech struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewSpeech(messenger service.Messenger, cache service.Cache) *Speech {
	return &Speech{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *Speech) Id() string {
	return enum.CommandSpeech
}

func (c *Speech) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)
	opts.Model = enum.ModelSpeech
	c.cache.Set(update.ChatId, opts)

	c.messenger.Send("Audio transcription model selected", update.ChatId, helper.GetContextCommands(opts.Model))
}
