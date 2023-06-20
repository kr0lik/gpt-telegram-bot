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
	"time"
)

const chatEditBatchLength = 25

var cancelProgress = make(map[dto.MessageId]func())

type Chat struct {
	messenger service.Messenger
	generator generator.Chat
	speech    service.Speech
	queue     service.Queue
	mu        *sync.Mutex
}

type output struct {
	message string
	err     error
}

func NewChat(messenger service.Messenger, generator generator.Chat, speech service.Speech, queue service.Queue) *Chat {
	return &Chat{
		messenger: messenger,
		generator: generator,
		speech:    speech,
		queue:     queue,
		mu:        &sync.Mutex{},
	}
}

func (h *Chat) Model() string {
	return enum.ModelChat
}

func (h *Chat) Handle(update dto.Income, ctx context.Context) {
	newCtx, cancel := context.WithCancel(ctx)
	cancelProgress[update.MessageId] = cancel

	defer delete(cancelProgress, update.MessageId)

	prompt := update.Message
	replayMessageId := update.MessageId
	callback := dto.Callback{
		Id:          enum.CallbackCancel,
		MessageId:   update.MessageId,
		Description: "cancel generation",
	}

	editMessageId := h.messenger.StartEdit("Processing...", replayMessageId, update.ChatId, [][]dto.Callback{{callback}}, helper.GetContextCommands(h.Model()))

	speechOutputChan := make(chan output)

	go h.getPromptFromSpeech(speechOutputChan, update, newCtx)

	select {
	case <-newCtx.Done():
		return
	case output, ok := <-speechOutputChan:
		if !ok {
			break
		}

		if output.err != nil {
			errorMessage := fmt.Sprintf("Failed to convert speech to text: %v", output.err)
			h.messenger.Replace(editMessageId, errorMessage, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
			return
		}

		replayMessageId = h.messenger.Replace(editMessageId, output.message, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		prompt = output.message

		editMessageId = h.messenger.StartEdit("Processing...", replayMessageId, update.ChatId, [][]dto.Callback{{callback}}, helper.GetContextCommands(h.Model()))
	}

	if prompt == "" {
		h.messenger.Replace(editMessageId, "Type your question", replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	queueOutputChan := make(chan output)

	go h.checkQueue(queueOutputChan, newCtx)

	func() {
		for {
			select {
			case <-newCtx.Done():
				return
			case output, ok := <-queueOutputChan:
				if !ok {
					return
				}

				h.messenger.Edit(editMessageId, output.message, replayMessageId, update.ChatId, [][]dto.Callback{{callback}}, helper.GetContextCommands(h.Model()))
			}
		}
	}()

	resultOutputChan := make(chan output)

	go h.generate(resultOutputChan, prompt, update.ChatId, newCtx)

	result := ""

	for {
		select {
		case <-newCtx.Done():
			return
		case output, ok := <-resultOutputChan:
			if !ok {
				if result == "" {
					result = "no answer received"
				}

				h.messenger.Replace(editMessageId, result, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
				return
			}

			if output.err != nil {
				errorMessage := fmt.Sprintf("Failed to generate text: %v", output.err)
				h.messenger.Replace(editMessageId, errorMessage, replayMessageId, update.ChatId, helper.GetContextCommands(h.Model()))
				return
			}

			editMessageId = h.messenger.Edit(editMessageId, result+"...", replayMessageId, update.ChatId, [][]dto.Callback{{callback}}, helper.GetContextCommands(h.Model()))
		}
	}
}

func (h *Chat) getPromptFromSpeech(outputChan chan<- output, update dto.Income, ctx context.Context) {
	output := output{"", nil}

	if update.Message == "" || update.AudioPath != "" {
		text, err := h.speech.ToText(update.AudioPath, service.SpeechOptions{}, ctx)
		if err != nil {
			output.err = err
			outputChan <- output
			close(outputChan)
			return
		}

		output.message = text
		outputChan <- output
	}

	close(outputChan)
}

func (h *Chat) checkQueue(outputChan chan<- output, ctx context.Context) {
	output := output{"", nil}

	lock, expire := h.queue.IsLocked()
	if lock {
		for {
			output.message = fmt.Sprintf("Queued for %s", expire.String())
			outputChan <- output

			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
			}

			lock, expire = h.queue.IsLocked()
			if !lock {
				output.message = "Processing..."
				outputChan <- output
				close(outputChan)
				return
			}
		}
	}

	close(outputChan)
}

func (h *Chat) generate(outputChan chan<- output, prompt string, chatId dto.ChatId, ctx context.Context) {
	output := output{"", nil}

	generatedTextStream, err := h.generator.Generate(prompt, chatId, ctx)
	if err != nil {
		output.err = err
		outputChan <- output
		close(outputChan)
		return
	}

	result := ""

	for generatedText := range generatedTextStream {
		result += generatedText

		if len(result) == 0 {
			continue
		}

		if len(result)%chatEditBatchLength == 0 {
			output.message = result
			outputChan <- output
		}
	}

	output.message = result
	outputChan <- output
	close(outputChan)
}

func (h *Chat) Callback(update dto.Income, ctx context.Context) {
	switch update.Callback.Command {
	case enum.CallbackCancel:
		cancel, ok := cancelProgress[update.Callback.MessageId]
		if ok {
			delete(cancelProgress, update.Callback.MessageId)
			cancel()
			h.messenger.Callback(update.Callback.Id, "Done")
		} else {
			h.messenger.Callback(update.Callback.Id, "Already done")
		}
	}
}
