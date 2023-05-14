package model

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/editor"
)

const defaultTextEditInstruction = "Fix the spelling mistakes"

var textEditProcessing = make(map[dto.ChatId]textEditProcess)

type TextEdit struct {
	messenger service.Messenger
	editor    editor.Text
}

type textEditProcess struct {
	prompt             string
	requestInstruction dto.MessageId
}

func NewTextEdit(messenger service.Messenger, editor editor.Text) *TextEdit {
	return &TextEdit{
		messenger: messenger,
		editor:    editor,
	}
}

func (h *TextEdit) Model() string {
	return enum.ModelTextEdit
}

func (h *TextEdit) Handle(update dto.Income, ctx context.Context) {
	process, ok := textEditProcessing[update.ChatId]
	if !ok && update.Message == "" {
		h.messenger.Replay("Type your prompt", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if !ok {
		h.requestInstruction(update)
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	instruction := defaultTextEditInstruction
	if update.Message != "" {
		instruction = update.Message
	}

	result, err := h.editor.Edit(process.prompt, instruction, ctx)

	delete(textEditProcessing, update.ChatId)

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

func (h *TextEdit) requestInstruction(update dto.Income) {
	commands := append([][]dto.Command{
		{
			{
				Id:            defaultTextEditInstruction,
				IsInstruction: true,
			},
		},
	}, helper.GetContextCommands(h.Model())...)

	messageId := h.messenger.Replay(fmt.Sprintf("Type instruction, for example: \"%s\"", defaultTextEditInstruction), update.MessageId, update.ChatId, commands)

	textEditProcessing[update.ChatId] = textEditProcess{
		prompt:             update.Message,
		requestInstruction: messageId,
	}
}
