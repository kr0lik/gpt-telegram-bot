package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/handler/model"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Stop struct {
	messenger service.Messenger
}

func NewStop(messenger service.Messenger) *Stop {
	return &Stop{
		messenger: messenger,
	}
}

func (c *Stop) Id() string {
	return enum.CommandChatStop
}

func (c *Stop) Process(update dto.Income) {
	model.Chat{}.StopProgress(update)

	c.messenger.Send("Type new request", update.ChatId, helper.GetContextCommands(enum.ModelChat))
}
