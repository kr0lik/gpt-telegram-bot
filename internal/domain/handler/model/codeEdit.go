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

const defaultCodeEditInstruction = "Fix the syntax and logic mistakes"

var codeEditProcessing = make(map[dto.ChatId]codeEditProcess)

type CodeEdit struct {
	messenger service.Messenger
	editor    editor.Code
}

type codeEditProcess struct {
	prompt             string
	requestInstruction dto.MessageId
}

func NewCodeEdit(messenger service.Messenger, editor editor.Code) *CodeEdit {
	return &CodeEdit{
		messenger: messenger,
		editor:    editor,
	}
}

func (h *CodeEdit) Model() string {
	return enum.ModelCodeEdit
}

func (h *CodeEdit) Handle(update dto.Income, ctx context.Context) {
	process, ok := codeEditProcessing[update.ChatId]
	if !ok && update.Message == "" {
		h.messenger.Replay("Type your prompt", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if !ok {
		h.requestInstruction(update)
		return
	}

	messageId := h.messenger.Replay("Processing...⏳", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	instruction := defaultCodeEditInstruction
	if update.Message != "" {
		instruction = update.Message
	}

	result, err := h.editor.Edit(process.prompt, instruction, ctx)

	delete(codeEditProcessing, update.ChatId)

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

func (h *CodeEdit) requestInstruction(update dto.Income) {
	commands := append([][]dto.Command{
		{
			{
				Id:            defaultCodeEditInstruction,
				IsInstruction: true,
			},
		},
	}, helper.GetContextCommands(h.Model())...)

	messageId := h.messenger.Replay(fmt.Sprintf("Type instruction, for example: \"%s\"", defaultCodeEditInstruction), update.MessageId, update.ChatId, commands)

	codeEditProcessing[update.ChatId] = codeEditProcess{
		prompt:             update.Message,
		requestInstruction: messageId,
	}
}
