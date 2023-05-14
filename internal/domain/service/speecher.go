package service

import "context"

type Speecher interface {
	ToText(audioPath string, ctx context.Context) (string, error)
}
