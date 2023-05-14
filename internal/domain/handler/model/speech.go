package model

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
)

type Speech struct {
	messenger service.Messenger
	cache     service.Cache
	speecher  service.Speecher
}

func NewSpeech(messenger service.Messenger, cache service.Cache, generator service.Speecher) *Speech {
	return &Speech{
		messenger: messenger,
		cache:     cache,
		speecher:  generator,
	}
}

func (h *Speech) Model() string {
	return enum.ModelSpeech
}

func (h *Speech) Handle(update dto.Income, ctx context.Context) {
	if update.AudioPath == "" {
		h.messenger.Replay("Send your audio to transcript", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	result, err := h.speecher.ToText(update.AudioPath, ctx)
	if err != nil {
		errorText := fmt.Sprintf("Failed to get audio transcription: %v", err)
		h.messenger.Replace(messageId, errorText, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if result == "" {
		result = "no transcription received"
	}

	h.messenger.Replace(messageId, result, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
}
