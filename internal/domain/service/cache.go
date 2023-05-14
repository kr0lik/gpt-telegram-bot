package service

import "gpt-telegran-bot/internal/domain/dto"

type Cache interface {
	Set(key dto.ChatId, value dto.Options)
	Has(key dto.ChatId) bool
	Get(key dto.ChatId) dto.Options
	Delete(key dto.ChatId)
}
