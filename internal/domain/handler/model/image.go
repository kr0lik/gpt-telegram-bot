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

type Image struct {
	messenger service.Messenger
	cache     service.Cache
	generator generator.Image
}

func NewImage(messenger service.Messenger, cache service.Cache, generator generator.Image) *Image {
	return &Image{
		messenger: messenger,
		cache:     cache,
		generator: generator,
	}
}

func (h *Image) Model() string {
	return enum.ModelImage
}

func (h *Image) Handle(update dto.Income, ctx context.Context) {
	if update.Message == "" {
		h.messenger.Replay("Type your prompt", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	messageId := h.messenger.Replay("Processing...", update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))

	urls, err := h.generator.Generate(update.Message, h.getOpts(update), ctx)
	if err != nil {
		errorText := fmt.Sprintf("Failed to generate image: %v", err)
		h.messenger.Replace(messageId, errorText, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
		return
	}

	h.messenger.ReplaceWithPhotos(messageId, urls, update.MessageId, update.ChatId, helper.GetContextCommands(h.Model()))
}

func (h *Image) getOpts(update dto.Income) generator.ImageOptions {
	imgOpts := new(generator.ImageOptions)

	opts := h.cache.Get(update.ChatId)

	if opts.Image.Count > 0 {
		imgOpts.Count = opts.Image.Count
	}

	if opts.Image.Size != "" {
		imgOpts.Size = opts.Image.Size
	}

	return *imgOpts
}

func (h *Image) Callback(update dto.Income, ctx context.Context) {

}
