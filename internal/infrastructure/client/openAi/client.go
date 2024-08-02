package openAi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/request"
	"gpt-telegran-bot/internal/infrastructure/client/openAi/dto/response"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	apiHost         = "https://api.openai.com/v1"
	jsonContentType = "application/json"
	dataContentType = "multipart/form-data"
)

var (
	asyncDataPrefix   = []byte("data: ")
	asyncDoneSequence = []byte("[DONE]")
)

type Client struct {
	apiKey string
	client *http.Client
}

type ClientConfig struct {
	ApiKey   string
	ProxyDsn string
}

func NewClient(config *ClientConfig) (*Client, error) {
	transport := http.Transport{}

	if config.ProxyDsn != "" {
		proxyURL, err := url.Parse(config.ProxyDsn)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy dsn: %v", err)
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &Client{
		apiKey: config.ApiKey,
		client: &http.Client{Transport: &transport},
	}, nil
}

func getSocks5Transport(config *ClientConfig) *http.Transport {
	if config.ProxyDsn == "" {
		return &http.Transport{}
	}

	if !strings.Contains(config.ProxyDsn, "socks5://") {
		fmt.Printf("Invalid format of proxy: %v, expected: socks5://host:port", config.ProxyDsn)
		return &http.Transport{}
	}

	proxyHost := strings.Replace(config.ProxyDsn, "socks5://", "", 1)

	dialer, err := proxy.SOCKS5("tcp", proxyHost, nil, proxy.Direct)
	if err != nil {
		fmt.Printf("Proxy connect error: %v\n", err)
		return &http.Transport{}
	}

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	return &http.Transport{
		DialContext: dialContext,
	}
}

func (c *Client) GetChatCompletions(request request.ChatCompletions, ctx context.Context) (*response.ChatCompletions, error) {
	resp, err := c.request("POST", "/chat/completions", jsonContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call chat/completions: %v", err)
	}

	result := new(response.ChatCompletions)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) GetChatCompletionsStream(request request.ChatCompletionsAsync, ctx context.Context) (<-chan *response.ChatCompletionsAsync, error) {
	request.Stream = true

	resp, err := c.request("POST", "/chat/completions", jsonContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call chat/completions: %v", err)
	}

	resCh := make(chan *response.ChatCompletionsAsync)

	go c.readChatCompletionsStream(resp, resCh, ctx)

	return resCh, nil
}

func (c *Client) GetCompletions(request request.Completions, ctx context.Context) (*response.Completions, error) {
	resp, err := c.request("POST", "/completions", jsonContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call completions: %v", err)
	}

	result := new(response.Completions)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) EditText(request request.Edits, ctx context.Context) (*response.Edits, error) {
	resp, err := c.request("POST", "/edits", jsonContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call edits: %v", err)
	}

	result := new(response.Edits)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) GetImages(request request.ImageGenerations, ctx context.Context) (*response.ImageGenerations, error) {
	resp, err := c.request("POST", "/images/generations", jsonContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call images/generations: %v", err)
	}

	result := new(response.ImageGenerations)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) GetImageEdits(request request.ImageEdits, ctx context.Context) (*response.ImageEdits, error) {
	resp, err := c.request("POST", "/images/edits", dataContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call images/edits: %v", err)
	}

	result := new(response.ImageEdits)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) GetImageVariations(request request.ImageVariations, ctx context.Context) (*response.ImageVariations, error) {
	resp, err := c.request("POST", "/images/variations", dataContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call images/variations: %v", err)
	}

	result := new(response.ImageVariations)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) GetAudioTranscription(request request.AudioTranscriptions, ctx context.Context) (*response.AudioTranscriptions, error) {
	resp, err := c.request("POST", "/audio/transcriptions", dataContentType, request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed call audio/transcriptions: %v", err)
	}

	result := new(response.AudioTranscriptions)

	err = c.getResponseObject(resp, result)

	return result, err
}

func (c *Client) request(method, path, contentType string, request interface{}, ctx context.Context) (*http.Response, error) {
	url := apiHost + path

	payload, err := c.preparePayload(request, &contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request payload: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while do request: %v", err)
	}

	if err := c.checkResponse(resp); err != nil {
		return resp, fmt.Errorf("bad response: %v", err)
	}

	return resp, nil
}

func (c *Client) preparePayload(request interface{}, contentType *string) (payload []byte, err error) {
	switch *contentType {
	case dataContentType:
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		if err := c.multipartFormWrite(writer, request); err != nil {
			return nil, fmt.Errorf("failed write multipart form data: %v", err)
		}

		payload = body.Bytes()
		*contentType = writer.FormDataContentType()
		break
	default:
		body, err := json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("failed encode reques: %v", err)
		}

		payload = body
	}

	return payload, nil
}

func (c *Client) multipartFormWrite(writer *multipart.Writer, request interface{}) error {
	v := reflect.ValueOf(request)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i).Tag.Get("json")
		value := v.Field(i).Interface()
		switch field {
		case "file", "image", "mask":
			filePath := fmt.Sprintf("%v", value)

			part, err := writer.CreateFormFile(field, filepath.Base(filePath))
			if err != nil {
				return err
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}

			if _, err = io.Copy(part, file); err != nil {
				return err
			}
		default:
			if err := writer.WriteField(field, fmt.Sprintf("%v", value)); err != nil {
				return err
			}
		}
	}

	writer.Close()

	return nil
}

func (c *Client) checkResponse(resp *http.Response) error {
	isSuccess := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !isSuccess {
		defer resp.Body.Close()

		err := new(response.Error)
		json.NewDecoder(resp.Body).Decode(err)

		return fmt.Errorf("returned code is %d, `%s`: %s", resp.StatusCode, err.Error.Type, err.Error.Message)
	}

	return nil
}

func (c *Client) getResponseObject(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}

	return nil
}

func (c *Client) readChatCompletionsStream(resp *http.Response, ch chan *response.ChatCompletionsAsync, ctx context.Context) {
	defer resp.Body.Close()
	defer close(ch)

	reader := bufio.NewReader(resp.Body)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("failed to read response: %v", err)
			break
		}

		if !bytes.HasPrefix(line, asyncDataPrefix) {
			continue
		}

		line = bytes.TrimPrefix(line, asyncDataPrefix)

		if bytes.HasPrefix(line, asyncDoneSequence) {
			break
		}

		res := new(response.ChatCompletionsAsync)

		if err := json.Unmarshal(line, res); err != nil {
			log.Printf("failed to decode response: %v", err)
			return
		}

		ch <- res
	}
}
