package openAi

import (
	"bytes"
	"context"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
	"gpt-telegran-bot/internal/infrastructure/util"
	"path/filepath"
	"strings"
)

const (
	speechTranscriptionModel = "whisper-1"
	speechTemperature        = 0
)

type Speech struct {
	client *openAi.Client
}

func NewSpeech(client *openAi.Client) *Speech {
	return &Speech{
		client: client,
	}
}

func (s *Speech) ToText(audioPath string, options service.SpeechOptions, ctx context.Context) (string, error) {
	audioPath, err := s.convertAudio(audioPath)
	if err != nil {
		return "", fmt.Errorf("error while get audio data: %v", err)
	}

	defer util.DeleteFile(audioPath)

	req := request.AudioTranscriptions{
		File:        audioPath,
		Model:       speechTranscriptionModel,
		Temperature: speechTemperature,
	}

	if options.Prompt != "" {
		req.Prompt = options.Prompt
	}

	response, err := s.client.GetAudioTranscription(req, ctx)
	if err != nil {
		return "", err
	}

	return response.Text, nil
}

func (s *Speech) convertAudio(audioPath string) (string, error) {
	defer util.DeleteFile(audioPath)

	newAudioPath := strings.TrimSuffix(audioPath, "."+filepath.Ext(audioPath)) + ".m4a"

	var out bytes.Buffer
	if err := ffmpeg.Input(audioPath).Output(newAudioPath, ffmpeg.KwArgs{"ab": "64k"}).WithOutput(&out).Run(); err != nil {
		return "", fmt.Errorf("%v: %v", err, out)
	}

	return newAudioPath, nil
}
