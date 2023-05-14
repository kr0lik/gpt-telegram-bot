package command

import "gpt-telegran-bot/internal/domain/dto"

type Handler interface {
	Id() string
	Process(update dto.Income)
}
