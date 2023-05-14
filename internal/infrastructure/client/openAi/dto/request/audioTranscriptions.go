package request

type AudioTranscriptions struct {
	File        []byte  `json:"file"`
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
}
