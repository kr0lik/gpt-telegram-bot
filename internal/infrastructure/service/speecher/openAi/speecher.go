package openAi

import (
	"context"
	"fmt"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
	"gpt-telegran-bot/internal/infrastructure/util"
	"os"
)

const (
	transcriptionModel = "whisper-1"
	speechTemperature  = 0
)

type Speecher struct {
	client *openAi.Client
}

func NewSpeecher(client *openAi.Client) *Speecher {
	return &Speecher{
		client: client,
	}
}

func (s *Speecher) ToText(audioPath string, ctx context.Context) (string, error) {
	fileData, err := os.ReadFile(audioPath)
	if err != nil {
		return "", fmt.Errorf("error whil readeing audio: %v", err)
	}

	util.DeleteFile(audioPath)

	// ToDo: Need Converter from ogg (Telegram Voice)

	req := request.AudioTranscriptions{
		File:        fileData,
		Model:       transcriptionModel,
		Temperature: speechTemperature,
	}

	response, err := s.client.GetAudioTranscription(req, ctx)
	if err != nil {
		return "", err
	}

	return response.Text, nil
}
