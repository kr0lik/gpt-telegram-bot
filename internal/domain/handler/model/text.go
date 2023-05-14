package model

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/generator"
)

type Text struct {
	messenger service.Messenger
	generator generator.Text
}

func NewText(messenger service.Messenger, generator generator.Text) *Text {
	return &Text{
		messenger: messenger,
		generator: generator,
	}
}

func (h *Text) Model() string {
	return enum.ModelText
}

func (h *Text) Handle(update dto.Income, ctx context.Context) {
	if update.Message == "" {
		h.messenger.Replay("Type your prompt", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	result, err := h.generator.Generate(update.Message, ctx)
	if err != nil {
		errorText := fmt.Sprintf("Failed to generate text: %v", err)
		h.messenger.Replace(messageId, errorText, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if result == "" {
		result = "no answer received"
	}

	h.messenger.Replace(messageId, result, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
}
