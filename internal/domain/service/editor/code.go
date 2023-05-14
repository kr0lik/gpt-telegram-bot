package editor

import (
	"context"
)

type Code interface {
	Edit(prompt string, instruction string, ctx context.Context) (string, error)
}
