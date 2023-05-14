package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Start struct {
	messenger   service.Messenger
	helpCommand *Help
}

func NewStart(messenger service.Messenger, cache service.Cache) *Start {
	return &Start{
		messenger:   messenger,
		helpCommand: NewHelp(messenger, cache),
	}
}

func (c *Start) Id() string {
	return enum.CommandStart
}

func (c *Start) Process(update dto.Income) {
	c.messenger.Send("Welcome to the GPT Telegram bot", update.ChatId, helper.GetContextCommands(""))

	c.helpCommand.Process(update)
}
