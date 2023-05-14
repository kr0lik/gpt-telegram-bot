package request

type ChatCompletions struct {
	Model       string         `json:"model"`
	Messages    []Conversation `json:"messages"`
	MaxTokens   int            `json:"max_tokens"`
	Temperature float32        `json:"temperature"`
}

type Conversation struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
