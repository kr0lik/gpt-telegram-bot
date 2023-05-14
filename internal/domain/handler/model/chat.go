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

const chatEditBachLength = 15

var cancelProgress = make(map[dto.ChatId]func())

type Chat struct {
	messenger service.Messenger
	generator generator.Chat
}

func NewChat(messenger service.Messenger, generator generator.Chat) *Chat {
	return &Chat{
		messenger: messenger,
		generator: generator,
	}
}

func (h *Chat) Model() string {
	return enum.ModelChat
}

func (h *Chat) Handle(update dto.Income, ctx context.Context) {
	_, ok := cancelProgress[update.ChatId]
	if ok {
		h.messenger.Replay("Your previous request in progress", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if update.Message == "" {
		h.messenger.Replay("Type your question", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	newCtx, cancel := context.WithCancel(ctx)
	cancelProgress[update.ChatId] = cancel

	generatedTextStream, err := h.generator.Generate(update.Message, update.ChatId, newCtx)
	if err != nil {
		errorText := fmt.Sprintf("Failed to generate text: %v", err)
		h.messenger.Replace(messageId, errorText, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		delete(cancelProgress, update.ChatId)
		return
	}

	result := ""
	isFirst := true

	for generatedText := range generatedTextStream {
		result += generatedText

		if len(result) == 0 {
			continue
		}

		if isFirst {
			messageId = h.messenger.StartEdit(messageId, result+"...", update.MessageId, update.ChatId)
			isFirst = false
			continue
		}

		if len(result)%chatEditBachLength == 0 {
			messageId = h.messenger.Edit(messageId, result+"...", update.MessageId, update.ChatId)
		}
	}

	if result == "" {
		result = "no answer received"
	}

	h.messenger.Replace(messageId, result, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	delete(cancelProgress, update.ChatId)
}

func (h Chat) StopProgress(update dto.Income) {
	cancel, ok := cancelProgress[update.ChatId]
	if ok {
		cancel()
	}
}
