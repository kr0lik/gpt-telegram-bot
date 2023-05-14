package command

import (
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/generator"
	"strings"
)

type Size struct {
	messenger      service.Messenger
	cache          service.Cache
	imageGenerator generator.Image
}

func NewSize(messenger service.Messenger, cache service.Cache, imageGenerator generator.Image) *Size {
	return &Size{
		messenger:      messenger,
		cache:          cache,
		imageGenerator: imageGenerator,
	}
}

func (c *Size) Id() string {
	return enum.CommandImageSize
}

func (c *Size) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)

	availableSizes := c.imageGenerator.GetAvailableImageSizes()

	for _, availableSize := range availableSizes {
		if strings.Contains(strings.ToLower(update.Message), availableSize) {
			opts.Image.Size = availableSize
			c.cache.Set(update.ChatId, opts)

			c.messenger.Send(fmt.Sprintf("Size %s selected", availableSize), update.ChatId, helper.GetContextCommands(opts.Model))

			return
		}
	}

	c.printSizeCommands(update.ChatId, opts.Model, availableSizes)
}

func (c *Size) printSizeCommands(chatId dto.ChatId, model string, availableSizes []string) {
	var sizeCommands []dto.Command

	for _, size := range availableSizes {
		sizeCommand := dto.Command{
			Id:          enum.CommandImageSize + " " + size,
			Description: fmt.Sprintf("Will answer with %s images", size),
		}

		sizeCommands = append(sizeCommands, sizeCommand)
	}

	toPrintCommands := [][]dto.Command{sizeCommands}

	allCommands := append(toPrintCommands, helper.GetContextCommands(model)...)

	c.messenger.PrintCommands("Wrong size command. Available size commands:", toPrintCommands, chatId, allCommands)
}
