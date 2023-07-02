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

const defaultInstruction = "get variations"

var imageEditProcessing = make(map[dto.ChatId]imageEditProcess)

type ImageEdit struct {
	messenger service.Messenger
	cache     service.Cache
	editor    editor.Image
}

type imageEditProcess struct {
	imageToEditPath    string
	requestInstruction dto.MessageId
}

func NewImageEdit(messenger service.Messenger, cache service.Cache, editor editor.Image) *ImageEdit {
	return &ImageEdit{
		messenger: messenger,
		cache:     cache,
		editor:    editor,
	}
}

func (h *ImageEdit) Model() string {
	return enum.ModelImageEdit
}

func (h *ImageEdit) Handle(update dto.Income, ctx context.Context) {
	process, ok := imageEditProcessing[update.ChatId]
	if !ok && update.ImagePath == "" {
		h.messenger.Replay("Send your image to edit", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	if !ok {
		h.requestInstruction(update)
		return
	}

	messageId := h.messenger.Replay("Processing...", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	instruction := update.Message
	if instruction == "" {
		instruction = update.Caption
	}

	var urls []string
	var err error

	if instruction == defaultInstruction || instruction == "" {
		urls, err = h.editor.Variations(process.imageToEditPath, h.getOpts(update), ctx)
	} else {
		urls, err = h.editor.Edit(process.imageToEditPath, instruction, h.getOpts(update), ctx)
	}

	if err != nil {
		errorText := fmt.Sprintf("Failed to eedit image: %v", err)
		h.messenger.Replace(messageId, errorText, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	h.messenger.ReplaceWithPhotos(messageId, urls, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
}

func (h *ImageEdit) getOpts(update dto.Income) editor.ImageOptions {
	imgOpts := new(editor.ImageOptions)

	opts := h.cache.Get(update.ChatId)

	if opts.Image.Count > 0 {
		imgOpts.Count = opts.Image.Count
	}

	if opts.Image.Size != "" {
		imgOpts.Size = opts.Image.Size
	}

	if len(update.ImagePath) == 2 {
		imgOpts.MaskPath = update.ImagePath
	}

	return *imgOpts
}

func (h *ImageEdit) requestInstruction(update dto.Income) {
	commands := append([][]dto.Command{
		{
			{
				Id:            defaultInstruction,
				IsInstruction: true,
			},
		},
	}, helper.GetContextCommands(h.Model())...)

	messageId := h.messenger.Replay(fmt.Sprintf("You can send mask and/or type instruction, or type %s", defaultInstruction), update.MessageId, update.ChatId, commands)

	imageEditProcessing[update.ChatId] = imageEditProcess{
		imageToEditPath:    update.ImagePath,
		requestInstruction: messageId,
	}
}

func (h *ImageEdit) Callback(update dto.Income, ctx context.Context) {

}
