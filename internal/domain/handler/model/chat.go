package model

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/generator"
	"sync"
)

const chatEditBatchLength = 25

var cancelProgress = make(map[dto.ChatId]func())

type Chat struct {
	messenger service.Messenger
	generator generator.Chat
	speech    service.Speech
	mu        *sync.Mutex
}

func NewChat(messenger service.Messenger, generator generator.Chat, speech service.Speech) *Chat {
	return &Chat{
		messenger: messenger,
		generator: generator,
		speech:    speech,
		mu:        &sync.Mutex{},
	}
}

func (h *Chat) Model() string {
	return enum.ModelChat
}

func (h *Chat) Handle(update dto.Income, ctx context.Context) {
	prompt := update.Message
	replayMessageId := update.MessageId

	h.mu.Lock()

	_, ok := cancelProgress[update.ChatId]
	if ok {
		h.messenger.Replay("Your previous request in progress", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		h.mu.Unlock()
		return
	}

	newCtx, cancel := context.WithCancel(ctx)
	cancelProgress[update.ChatId] = cancel

	defer delete(cancelProgress, update.ChatId)

	h.mu.Unlock()

	if prompt == "" && update.AudioPath != "" {
		messageId := h.messenger.Replay("Processing...⏳", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))

		text, err := h.speech.ToText(update.AudioPath, service.SpeechOptions{}, ctx)
		if err != nil {
			errorText := fmt.Sprintf("Error while transcript audio: %v", err)
			h.messenger.Replace(messageId, errorText, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
			return
		}

		prompt = text
		replayMessageId = h.messenger.Replace(messageId, text, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
	}

	if prompt == "" {
		h.messenger.Replay("Type your question", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	generatedTextStream, err := h.generator.Generate(prompt, update.ChatId, newCtx)
	if err != nil {
		errorText := fmt.Sprintf("Failed to generate text: %v", err)
		h.messenger.Replace(messageId, errorText, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	result := h.readStream(generatedTextStream, &messageId, replayMessageId, update.ChatId)

	if result == "" {
		result = "no answer received"
	}

	h.messenger.Replace(messageId, result, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
}

func (h *Chat) readStream(generatedTextStream generator.ChatStreamChannel, messageId *dto.MessageId, replayMessageId dto.MessageId, chatId dto.ChatId) string {
	result := ""
	isFirst := true

	for generatedText := range generatedTextStream {
		result += generatedText

		if len(result) == 0 {
			continue
		}

		if isFirst {
			*messageId = h.messenger.StartEdit(*messageId, result+"...", replayMessageId, chatId)
			isFirst = false
			continue
		}

		if len(result)%chatEditBatchLength == 0 {
			*messageId = h.messenger.Edit(*messageId, result+"...", replayMessageId, chatId)
		}
	}

	return result
}

func (h Chat) StopProgress(update dto.Income) {
	cancel, ok := cancelProgress[update.ChatId]
	if ok {
		delete(cancelProgress, update.ChatId)
		cancel()
	}
}
