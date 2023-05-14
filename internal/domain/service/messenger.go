package service

import (
	"context"
	"gpt-telegran-bot/internal/domain/dto"
)

type UpdatesChannel <-chan dto.Income

type Messenger interface {
	PrintCommands(title string, printCommands [][]dto.Command, chatId dto.ChatId, commands [][]dto.Command)

	Listen(ctx context.Context) (UpdatesChannel, error)
	Send(message string, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId
	Replay(message string, replayId dto.MessageId, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId
	Replace(messageId dto.MessageId, newMessage string, replayId dto.MessageId, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId
	ReplaceWithPhotos(messageId dto.MessageId, urls []string, replayId dto.MessageId, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId

	StartEdit(messageId dto.MessageId, newMessage string, replayId dto.MessageId, chatId dto.ChatId) dto.MessageId
	Edit(messageId dto.MessageId, newMessage string, replayId dto.MessageId, chatId dto.ChatId) dto.MessageId
}
