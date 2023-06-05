package request

type ImageEdits struct {
	Image  string `json:"image"`
	Mask   string `json:"mask"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}
