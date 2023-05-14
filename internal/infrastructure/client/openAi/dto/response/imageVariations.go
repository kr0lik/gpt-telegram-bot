package response

type ImageVariations struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}
