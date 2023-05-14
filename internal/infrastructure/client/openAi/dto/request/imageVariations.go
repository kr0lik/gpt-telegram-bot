package request

type ImageVariations struct {
	Image []byte `json:"image"`
	N     int    `json:"n"`
	Size  string `json:"size"`
}
