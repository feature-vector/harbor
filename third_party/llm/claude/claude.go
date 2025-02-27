package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/feature-vector/harbor/base/hc"
	"github.com/feature-vector/harbor/third_party/llm"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	ModelClaude1dot3 = "claude-1.3"
	ModelClaude2     = "claude-2"
)

var (
	defaultTimeout = 30 * time.Second
)

type ChatClient struct {
	MaxTokensToSample int
	AuthKey           string
	Timeout           time.Duration

	Verbose bool
}

func NewChatClient(authKey string) *ChatClient {
	return &ChatClient{
		MaxTokensToSample: 256,
		AuthKey:           authKey,
		Timeout:           defaultTimeout,
	}
}

func (c *ChatClient) buildBody(req *llm.RunModelRequest) map[string]interface{} {
	var ss strings.Builder
	if req.Prompt != "" {
		ss.WriteString(req.Prompt)
		ss.WriteString("\n\n")
	}
	if len(req.History) != 0 {
		for i := range req.History {
			h := req.History[i]
			ss.WriteString("Human: ")
			ss.WriteString(h[0])
			ss.WriteString("\n\n")
			ss.WriteString("Assistant: ")
			ss.WriteString(h[1])
			ss.WriteString("\n\n")
		}
	}
	ss.WriteString("Human: ")
	ss.WriteString(req.Input)
	ss.WriteString("\n\n")
	ss.WriteString("Assistant: ")
	body := map[string]interface{}{
		"stream":               req.Stream,
		"model":                req.ModelName,
		"prompt":               ss.String(),
		"max_tokens_to_sample": c.MaxTokensToSample,
	}
	return body
}

func (c *ChatClient) buildHttpRequest(ctx context.Context, req *llm.RunModelRequest) (*http.Request, error) {
	body := c.buildBody(req)
	if c.Verbose {
		zap.L().Info("ChatClient build body", zap.Any("body", body))
	}
	bodyBytes, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://api.anthropic.com/v1/complete",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("x-api-key", c.AuthKey)
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func (c *ChatClient) Execute(ctx context.Context, req *llm.RunModelRequest) (*llm.RunModelResponse, error) {
	httpReq, err := c.buildHttpRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	httpResponse, err := hc.Do(httpReq)
	if err == context.DeadlineExceeded || err == context.Canceled {
		zap.L().Warn("call claude timeout")
		return nil, fmt.Errorf("call claude timeout: %w", err)
	}
	if err != nil {
		zap.L().Error("call claude returns error", zap.Error(err))
		return nil, fmt.Errorf("call claude returns client error: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		defer httpResponse.Body.Close()
		body, err := io.ReadAll(httpResponse.Body)
		zap.L().Error("call claude returns http error",
			zap.Int("status_code", httpResponse.StatusCode),
			zap.ByteString("body", body),
			zap.Error(err),
		)
		return nil, fmt.Errorf("call claude returns http %d error", httpResponse.StatusCode)
	}
	if req.Stream {
		body := httpResponse.Body
		return &llm.RunModelResponse{
			Stream: true,
			StreamResult: &claudeStreamResponse{
				body:   body,
				reader: bufio.NewReader(body),
			},
		}, nil
	} else {
		defer httpResponse.Body.Close()
		respBytes, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
		if c.Verbose {
			zap.L().Info("claude chat response", zap.ByteString("resp", respBytes))
		}
		var result string
		r := &chatResp{}
		err = json.Unmarshal(respBytes, r)
		result = r.Completion
		return &llm.RunModelResponse{
			Result: result,
		}, nil
	}
}
