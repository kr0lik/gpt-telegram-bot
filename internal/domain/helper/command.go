package helper

import (
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/enum"
)

func GetAllCommands(model string) [][]dto.Command {
	commands := GetContextCommands(model)

	commands = append(commands, []dto.Command{
		dto.ChatCommand,
		dto.TextCommand,
		dto.TextEditCommand,
		dto.CodeEditCommand,
		dto.ImageCommand,
		dto.ImageEditCommand,
		dto.SpeechCommand,
	})

	commands = append(commands, []dto.Command{
		dto.HelpCommand,
		dto.StatusCommand,
	})

	return commands
}

func GetContextCommands(model string) [][]dto.Command {
	commands := make([][]dto.Command, 0)

	switch model {
	case enum.ModelChat:
		commands = append(commands, []dto.Command{
			dto.ChatNewCommand,
			dto.ChatStopCommand,
		})
	case enum.ModelImage:
		commands = append(commands, []dto.Command{
			dto.ImageCountCommand,
			dto.ImageSizeCommand,
		})
	case enum.ModelImageEdit:
		commands = append(commands, []dto.Command{
			dto.ImageCountCommand,
			dto.ImageSizeCommand,
		})
	}

	return commands
}
