package model

import (
	"context"
	"gpt-telegran-bot/internal/domain/dto"
)

type Handler interface {
	Model() string
	Handle(update dto.Income, ctx context.Context)
	Callback(update dto.Income, ctx context.Context)
}
