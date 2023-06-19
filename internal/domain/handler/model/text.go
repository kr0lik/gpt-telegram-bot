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
	speech    service.Speech
}

func NewText(messenger service.Messenger, generator generator.Text, speech service.Speech) *Text {
	return &Text{
		messenger: messenger,
		generator: generator,
		speech:    speech,
	}
}

func (h *Text) Model() string {
	return enum.ModelText
}

func (h *Text) Handle(update dto.Income, ctx context.Context) {
	prompt := update.Message
	replayMessageId := update.MessageId

	if update.Message == "" && update.AudioPath != "" {
		messageId := h.messenger.Replay("Processing...⏳", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))

		text, err := h.speech.ToText(update.AudioPath, service.SpeechOptions{}, ctx)
		if err != nil {
			errorText := fmt.Sprintf("Error while transcript audio: %v", err)
			h.messenger.Replace(messageId, errorText, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
			return
		}

		replayMessageId = h.messenger.Replace(messageId, text, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

		prompt = text
	}

	if prompt == "" {
		h.messenger.Replay("Type your prompt", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	result, err := h.generator.Generate(prompt, ctx)
	if err != nil {
		errorText := fmt.Sprintf("Failed to generate text: %v", err)
		h.messenger.Replace(messageId, errorText, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if result == "" {
		result = "no answer received"
	}

	h.messenger.Replace(messageId, result, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
}
