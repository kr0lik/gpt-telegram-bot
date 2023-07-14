package cache

import (
	"gpt-telegran-bot/internal/domain/dto"
	"sync"
)

type Memory struct {
	data map[dto.ChatId]dto.Options
	mu   *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		data: make(map[dto.ChatId]dto.Options),
		mu:   &sync.Mutex{},
	}
}

func (c *Memory) Set(key dto.ChatId, value dto.Options) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *Memory) Has(key dto.ChatId) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.data[key]

	return ok
}

func (c *Memory) Get(key dto.ChatId) dto.Options {
	c.mu.Lock()
	defer c.mu.Unlock()

	res, ok := c.data[key]

	if ok {
		return res
	}

	return dto.Options{}
}

func (c *Memory) Delete(key dto.ChatId) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
