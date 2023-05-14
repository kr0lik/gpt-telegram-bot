package command

import (
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
	"gpt-telegran-bot/internal/domain/helper"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/domain/service/generator"
	"strconv"
	"strings"
)

type Count struct {
	messenger      service.Messenger
	cache          service.Cache
	imageGenerator generator.Image
}

func NewCount(messenger service.Messenger, cache service.Cache, imageGenerator generator.Image) *Count {
	return &Count{
		messenger:      messenger,
		cache:          cache,
		imageGenerator: imageGenerator,
	}
}

func (c *Count) Id() string {
	return enum.CommandImageCount
}

func (c *Count) Process(update dto.Income) {
	opts := c.cache.Get(update.ChatId)

	maxCount := c.imageGenerator.GetMaxImageCount()

	for i := maxCount; i > 0; i-- {
		if strings.Contains(update.Message, strconv.Itoa(i)) {
			opts.Image.Count = i
			c.cache.Set(update.ChatId, opts)

			c.messenger.Send(fmt.Sprintf("Count %d selected", i), update.ChatId, helper.GetContextCommands(opts.Model))

			return
		}
	}

	c.printSizeCommands(update.ChatId, opts.Model, maxCount)
}

func (c *Count) printSizeCommands(chatId dto.ChatId, model string, maxCount int) {
	var countCmds []dto.Command

	for i := 1; i <= maxCount; i++ {
		countCmd := dto.Command{
			Id:          enum.CommandImageCount + " " + strconv.Itoa(i),
			Description: fmt.Sprintf("Will be %d item in answer", i),
		}

		countCmds = append(countCmds, countCmd)
	}

	toPrintCommands := [][]dto.Command{countCmds}

	allCommands := append(toPrintCommands, helper.GetContextCommands(model)...)

	c.messenger.PrintCommands("Wrong count command. Available count commands:", toPrintCommands, chatId, allCommands)
}
