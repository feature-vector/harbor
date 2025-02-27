package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/feature-vector/harbor/base/hc"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type ImageSize string

const (
	ImageSize256  ImageSize = "256x256"
	ImageSize512  ImageSize = "512x512"
	ImageSize1024 ImageSize = "1024x1024"
)

type ImageClient struct {
	AuthKey string
}

type GenerateImageRequest struct {
	Prompt string    `json:"prompt"`
	N      int       `json:"n"`
	Size   ImageSize `json:"size"`
}

type generateImageResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func NewImageClient(authKey string) *ImageClient {
	return &ImageClient{
		AuthKey: authKey,
	}
}

func (c *ImageClient) GenerateImage(ctx context.Context, req *GenerateImageRequest) ([]string, error) {
	if req.N == 0 {
		req.N = 1
	}
	if req.Size == "" {
		req.Size = ImageSize256
	}
	bodyBytes, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://api.openai.com/v1/images/generations",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthKey))
	httpReq.Header.Set("Content-Type", "application/json")

	httpResponse, err := hc.Do(httpReq)
	if err == context.DeadlineExceeded || err == context.Canceled {
		zap.L().Warn("call openai timeout")
		return nil, fmt.Errorf("call openai timeout: %w", err)
	}
	if err != nil {
		zap.L().Error("call openai returns error", zap.Error(err))
		return nil, fmt.Errorf("call openai returns client error: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(httpResponse.Body)
		zap.L().Error("call openai returns http error",
			zap.Int("status_code", httpResponse.StatusCode),
			zap.ByteString("body", body),
			zap.Error(err),
		)
		return nil, fmt.Errorf("call openai returns http %d error", httpResponse.StatusCode)
	}
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	resp := &generateImageResponse{}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, d := range resp.Data {
		ret = append(ret, d.Url)
	}
	return ret, nil
}
