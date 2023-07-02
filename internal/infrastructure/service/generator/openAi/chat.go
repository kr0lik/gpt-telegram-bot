package openAi

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/service/generator"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/response"
	"strings"
)

const (
	ChatModel = "gpt-3.5-turbo"

	ChatMaxTokens   = 2500
	ChatTemperature = float32(0.7)
)

var chatConversationHistory = make(map[dto.ChatId][]request.Conversation)

type Chat struct {
	client *openAi.Client
	model  string
}

func NewChat(client *openAi.Client) *Chat {
	return &Chat{
		client: client,
		model:  ChatModel,
	}
}

func (g *Chat) ClearConversation(chatId dto.ChatId) {
	chatConversationHistory[chatId] = []request.Conversation{
		{
			Role:    "system",
			Content: "You are a helpful AI assistant.",
		},
	}
}

func (g *Chat) Generate(prompt string, chatId dto.ChatId, ctx context.Context) (generator.ChatStreamChannel, error) {
	if _, ok := chatConversationHistory[chatId]; !ok {
		g.ClearConversation(chatId)
	}

	chatConversationHistory[chatId] = append(chatConversationHistory[chatId], request.Conversation{
		Role:    "user",
		Content: prompt,
	})

	req := request.ChatCompletionsAsync{
		Model:       g.model,
		Messages:    chatConversationHistory[chatId],
		MaxTokens:   ChatMaxTokens,
		Temperature: ChatTemperature,
		TopP:        1,
	}

	respCh, err := g.client.GetChatCompletionsStream(req, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Please reduce the length of the messages or completion") && len(chatConversationHistory[chatId]) > 2 {
			h := chatConversationHistory[chatId]
			chatConversationHistory[chatId] = append(h[:1], h[2:len(h)-1]...)
			return g.Generate(prompt, chatId, ctx)
		}

		return nil, err
	}

	resCh := make(chan string)

	go g.stream(chatId, respCh, resCh, ctx)

	return resCh, nil
}

func (g *Chat) stream(chatId dto.ChatId, respCh <-chan *response.ChatCompletionsAsync, resCh chan<- string, ctx context.Context) {
	defer close(resCh)

	fullText := ""

	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-respCh:
			if !ok {
				return
			}

			if len(resp.Choices) == 0 {
				resCh <- fmt.Sprintf("Invalid stream response: %v", resp)
				return
			}

			resCh <- resp.Choices[0].Delta.Content

			fullText += resp.Choices[0].Delta.Content

			if resp.Choices[0].FinishReason != "" {
				chatConversationHistory[chatId] = append(chatConversationHistory[chatId], request.Conversation{
					Role:    "assistant",
					Content: fullText,
				})

				break
			}
		}
	}
}
