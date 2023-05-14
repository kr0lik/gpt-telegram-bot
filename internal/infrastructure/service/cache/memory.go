package cache

import "gpt-telegran-bot/internal/domain/dto"

type Memory struct {
	data map[dto.ChatId]dto.Options
}

func NewMemory() *Memory {
	return &Memory{
		data: make(map[dto.ChatId]dto.Options),
	}
}

func (c *Memory) Set(key dto.ChatId, value dto.Options) {
	c.data[key] = value
}

func (c *Memory) Has(key dto.ChatId) bool {
	_, ok := c.data[key]

	return ok
}

func (c *Memory) Get(key dto.ChatId) dto.Options {
	res, ok := c.data[key]

	if ok {
		return res
	}

	return dto.Options{}
}

func (c *Memory) Delete(key dto.ChatId) {
	delete(c.data, key)
}
