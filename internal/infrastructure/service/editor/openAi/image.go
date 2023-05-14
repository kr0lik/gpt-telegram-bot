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
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error while readeing image: %v", err)
	}

	util.DeleteFile(imagePath)

	image, err := e.toPng(fileData)
	if err != nil {
		return nil, fmt.Errorf("error while comvert to png: %v", err)
	}

	req := request.ImageEdits{
		Prompt: instruction,
		Image:  image,
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
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error whil readeing image: %v", err)
	}

	util.DeleteFile(imagePath)

	image, err := e.toPng(fileData)
	if err != nil {
		return nil, fmt.Errorf("error while comvert to png: %v", err)
	}

	req := request.ImageVariations{
		Image: image,
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

func (e *Image) toPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return imageBytes, fmt.Errorf("error while decode jpeg: %v", err)
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return imageBytes, fmt.Errorf("error while encode png: %v", err)
		}

		return buf.Bytes(), nil
	}

	return imageBytes, nil
}

func (e *Image) applyOptions(req interface{}, options editor.ImageOptions) error {
	switch dto := req.(type) {
	case *request.ImageEdits:
		if options.MaskPath != "" {
			MaskData, err := os.ReadFile(options.MaskPath)
			if err != nil {
				return fmt.Errorf("error while readeing image mask: %v", err)
			}

			util.DeleteFile(options.MaskPath)

			mask, err := e.toPng(MaskData)
			if err != nil {
				return fmt.Errorf("error while comvert to png: %v", err)
			}

			dto.Mask = mask
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
