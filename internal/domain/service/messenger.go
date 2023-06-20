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

	StartEdit(message string, replayId dto.MessageId, chatId dto.ChatId, callbacks [][]dto.Callback, commands [][]dto.Command) dto.MessageId
	Edit(messageId dto.MessageId, newMessage string, replayId dto.MessageId, chatId dto.ChatId, callbacks [][]dto.Callback, commands [][]dto.Command) dto.MessageId

	Callback(callbackId, message string)
}
