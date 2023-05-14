package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Help struct {
	messenger service.Messenger
	cache     service.Cache
}

func NewHelp(messenger service.Messenger, cache service.Cache) *Help {
	return &Help{
		messenger: messenger,
		cache:     cache,
	}
}

func (c *Help) Id() string {
	return enum.CommandHelp
}

func (c *Help) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)

	c.messenger.PrintCommands("Available commands:", helper.GetAllCommands(""), update.ChatId, helper.GetContextCommands(opts.Model))
}
