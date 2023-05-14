package openAi

import (
	"context"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
)

const (
	TextEditModel       = "text-davinci-edit-001"
	TextEditTemperature = float32(0.3)
)

type Text struct {
	client *openAi.Client
}

func NewText(client *openAi.Client) *Text {
	return &Text{
		client: client,
	}
}

func (e *Text) Edit(prompt string, instruction string, ctx context.Context) (string, error) {
	req := request.Edits{
		Model:       TextEditModel,
		Input:       prompt,
		Instruction: instruction,
		Temperature: TextEditTemperature,
	}

	response, err := e.client.EditText(req, ctx)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
