package command

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/generator"
)

type New struct {
	messenger service.Messenger
	chat      generator.Chat
}

func NewNew(messenger service.Messenger, chat generator.Chat) *New {
	return &New{
		messenger: messenger,
		chat:      chat,
	}
}

func (c *New) Id() string {
	return enum.CommandChatNew
}

func (c *New) Process(update dto.Income) {
	c.chat.ClearConversation(update.ChatId)

	c.messenger.Send("Conversion history cleared", update.ChatId, helper.GetContextCommands(enum.ModelChat))
}
