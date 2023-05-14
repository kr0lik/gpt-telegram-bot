package openAi

import (
	"context"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
)

const (
	CodeEditModel       = "code-davinci-edit-001"
	CodeEditTemperature = float32(0.3)
)

type Code struct {
	client *openAi.Client
}

func NewCode(client *openAi.Client) *Code {
	return &Code{
		client: client,
	}
}

func (e *Code) Edit(prompt string, instruction string, ctx context.Context) (string, error) {
	req := request.Edits{
		Model:       CodeEditModel,
		Input:       prompt,
		Instruction: instruction,
		Temperature: CodeEditTemperature,
	}

	response, err := e.client.EditText(req, ctx)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
