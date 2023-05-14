package request

type Edits struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`
	Instruction string  `json:"instruction"`
	Temperature float32 `json:"temperature"`
}
