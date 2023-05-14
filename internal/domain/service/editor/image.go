package editor

import (
	"context"
)

type ImageOptions struct {
	Count    int
	Size     string
	MaskPath string
}

type Image interface {
	GetAvailableImageSizes() []string
	GetMaxImageCount() int
	Edit(imagePath string, instruction string, options ImageOptions, ctx context.Context) (imageUrls []string, err error)
	Variations(imagePath string, options ImageOptions, ctx context.Context) (imageUrls []string, err error)
}
