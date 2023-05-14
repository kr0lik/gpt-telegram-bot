package openAi

import (
	"context"
	"gpt-telegran-bot/internal/domain/service/generator"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
)

var availableImageSizes = []string{
	"256x256",
	"512x512",
	"1024x1024",
}

const defaultImageCount = 3
const defaultImageSize = "1024x1024"
const maxImageCount = 5

type Image struct {
	client *openAi.Client
}

func NewImage(client *openAi.Client) *Image {
	return &Image{
		client: client,
	}
}

func (g Image) GetAvailableImageSizes() []string {
	return availableImageSizes
}

func (g Image) GetMaxImageCount() int {
	return maxImageCount
}

func (g *Image) Generate(prompt string, options generator.ImageOptions, ctx context.Context) ([]string, error) {
	req := request.ImageGenerations{
		Prompt: prompt,
	}

	g.applyOptions(&req, options)

	response, err := g.client.GetImages(req, ctx)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(response.Data))

	for i, u := range response.Data {
		res[i] = u.Url
	}

	return res, nil
}

func (g *Image) applyOptions(req *request.ImageGenerations, options generator.ImageOptions) {
	req.N = defaultImageCount
	if options.Count > 0 {
		req.N = options.Count
	}

	req.Size = defaultImageSize
	if options.Size != "" {
		req.Size = options.Size
	}
}
