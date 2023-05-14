package generator

import (
	"context"
)

type Text interface {
	Generate(prompt string, ctx context.Context) (string, error)
}
