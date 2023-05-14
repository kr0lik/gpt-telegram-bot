package generator

import (
	"context"
	"gpt-telegran-bot/internal/domain/dto"
)

type ChatStreamChannel <-chan string

type Chat interface {
	Generate(prompt string, chatId dto.ChatId, ctx context.Context) (ChatStreamChannel, error)
	ClearConversation(chatId dto.ChatId)
}
