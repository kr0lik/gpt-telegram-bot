package service

import "context"

type SpeechOptions struct {
	Prompt string
}

type Speech interface {
	ToText(audioPath string, options SpeechOptions, ctx context.Context) (string, error)
}
