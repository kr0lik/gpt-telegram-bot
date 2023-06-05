package openAi

import (
	"bytes"
	"context"
	"fmt"
	"gpt-telegran-bot/internal/domain/service/editor"
	"gpt-telegran-bot/internal/infrastructure/client/openAi"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
	"gpt-telegran-bot/internal/infrastructure/util"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func (e Image) GetAvailableImageSizes() []string {
	return availableImageSizes
}

func (e Image) GetMaxImageCount() int {
	return maxImageCount
}

func (e *Image) Edit(imagePath string, instruction string, options editor.ImageOptions, ctx context.Context) ([]string, error) {
	imagePath, err := e.convertImage(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error while comvert to png: %v", err)
	}

	defer util.DeleteFile(imagePath)

	req := request.ImageEdits{
		Prompt: instruction,
		Image:  imagePath,
	}

	if err := e.applyOptions(&req, options); err != nil {
		return nil, fmt.Errorf("error while apply options: %v", err)
	}

	response, err := e.client.GetImageEdits(req, ctx)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(response.Data))

	for i, u := range response.Data {
		res[i] = u.Url
	}

	return res, nil
}

func (e *Image) Variations(imagePath string, options editor.ImageOptions, ctx context.Context) ([]string, error) {
	imagePath, err := e.convertImage(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error while comvert to png: %v", err)
	}

	defer util.DeleteFile(imagePath)

	req := request.ImageVariations{
		Image: imagePath,
	}

	if err := e.applyOptions(&req, options); err != nil {
		return nil, fmt.Errorf("error while apply options: %v", err)
	}

	response, err := e.client.GetImageVariations(req, ctx)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(response.Data))

	for i, u := range response.Data {
		res[i] = u.Url
	}

	return res, nil
}

func (e *Image) convertImage(imagePath string) (string, error) {
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return imagePath, fmt.Errorf("error while readeing image: %v", err)
	}

	contentType := http.DetectContentType(fileData)

	switch contentType {
	case "image/jpeg":
		defer util.DeleteFile(imagePath)

		img, err := jpeg.Decode(bytes.NewReader(fileData))
		if err != nil {
			return imagePath, fmt.Errorf("error while decode jpeg: %v", err)
		}

		newImagePath := strings.TrimSuffix(imagePath, "."+filepath.Ext(imagePath)) + ".png"

		file, err := os.Create(newImagePath)
		if err != nil {
			return imagePath, fmt.Errorf("error while creating new image file: %v", err)
		}

		if err := png.Encode(file, img); err != nil {
			return newImagePath, fmt.Errorf("error while encoding png: %v", err)
		}

		return newImagePath, nil
	}

	return imagePath, nil
}

func (e *Image) applyOptions(req interface{}, options editor.ImageOptions) error {
	switch dto := req.(type) {
	case *request.ImageEdits:
		if options.MaskPath != "" {
			MaskPath, err := e.convertImage(options.MaskPath)
			if err != nil {
				return fmt.Errorf("error while comvert to png: %v", err)
			}

			defer util.DeleteFile(MaskPath)

			dto.Mask = MaskPath
		}

		dto.N = defaultImageCount
		if options.Count > 0 {
			dto.N = options.Count
		}

		dto.Size = defaultImageSize
		if options.Size != "" {
			dto.Size = options.Size
		}
		break

	case *request.ImageVariations:
		dto.N = defaultImageCount
		if options.Count > 0 {
			dto.N = options.Count
		}

		dto.Size = defaultImageSize
		if options.Size != "" {
			dto.Size = options.Size
		}
		break
	}

	return nil
}
