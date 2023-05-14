package request

type ImageEdits struct {
	Image  []byte `json:"image"`
	Mask   []byte `json:"mask"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}
