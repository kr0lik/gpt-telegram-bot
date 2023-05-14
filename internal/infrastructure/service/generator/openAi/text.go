package openAi

import (
	"context"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
)

const (
	TextModel       = "text-davinci-003"
	TextMaxTokens   = 4000
	TextTemperature = float32(0.3)
)

type Text struct {
	client *openAi.Client
}

func NewText(client *openAi.Client) *Text {
	return &Text{
		client: client,
	}
}

func (g *Text) Generate(prompt string, ctx context.Context) (string, error) {
	req := request.Completions{
		Model:       TextModel,
		Prompt:      prompt,
		MaxTokens:   TextMaxTokens,
		Temperature: TextTemperature,
	}

	response, err := g.client.GetCompletions(req, ctx)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
