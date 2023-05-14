package editor

import (
	"context"
)

type Text interface {
	Edit(prompt string, instruction string, ctx context.Context) (string, error)
}
