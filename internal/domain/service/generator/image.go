package generator

import (
	"context"
)

type ImageOptions struct {
	Count int
	Size  string
}

type Image interface {
	GetAvailableImageSizes() []string
	GetMaxImageCount() int
	Generate(prompt string, options ImageOptions, ctx context.Context) (imageUrls []string, err error)
}
