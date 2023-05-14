package response

type ImageEdits struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}
