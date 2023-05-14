package request

type ChatCompletionsAsync struct {
	Model       string         `json:"model"`
	Messages    []Conversation `json:"messages"`
	MaxTokens   int            `json:"max_tokens"`
	Temperature float32        `json:"temperature"`
	Stream      bool           `json:"stream"`
	TopP        float32        `json:"top_p"`
}
