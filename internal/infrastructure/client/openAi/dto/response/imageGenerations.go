package response

type ImageGenerations struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}
