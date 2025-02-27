package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/feature-vector/harbor/base/hc"
	"github.com/feature-vector/harbor/third_party/llm"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	ModelGPT3dot5    = "gpt-3.5-turbo"
	ModelGPT4        = "gpt-4-0613"
	ModelDavinci     = "text-davinci-003"
	ModelDavinciEdit = "text-davinci-edit-001"

	urlChat        = "/v1/chat/completions"
	urlCompletions = "/v1/completions"
	urlEdits       = "/v1/edits"

	roleSystem    = "system"
	roleAssistant = "assistant"
	roleUser      = "user"
)

type ChatClient struct {
	Temperature         float32
	CompletionsMaxToken int
	AuthKey             string

	Verbose bool
}

func NewChatClient(authKey string) *ChatClient {
	return &ChatClient{
		Temperature:         0.7,
		CompletionsMaxToken: 256,
		AuthKey:             authKey,
	}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *ChatClient) buildGPTBody(req *llm.RunModelRequest) map[string]interface{} {
	var messages []*message
	if req.Prompt != "" {
		messages = append(messages, &message{
			Role:    roleSystem,
			Content: req.Prompt,
		})
	}
	if len(req.History) != 0 {
		for i := range req.History {
			h := req.History[i]
			messages = append(messages, &message{
				Role:    roleUser,
				Content: h[0],
			})
			messages = append(messages, &message{
				Role:    roleAssistant,
				Content: h[1],
			})
		}
	}
	messages = append(messages, &message{
		Role:    roleUser,
		Content: req.Input,
	})
	body := map[string]interface{}{
		"stream":      req.Stream,
		"model":       req.ModelName,
		"messages":    messages,
		"temperature": c.Temperature,
	}
	return body
}

func (c *ChatClient) buildCompletionsBody(req *llm.RunModelRequest) map[string]interface{} {
	return map[string]interface{}{
		"stream":      req.Stream,
		"model":       req.ModelName,
		"prompt":      req.Input,
		"max_tokens":  c.CompletionsMaxToken,
		"temperature": c.Temperature,
	}
}

func (c *ChatClient) buildEditsBody(req *llm.RunModelRequest) map[string]interface{} {
	return map[string]interface{}{
		"model":       req.ModelName,
		"input":       req.Input,
		"instruction": req.Prompt,
		"temperature": c.Temperature,
	}
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

func (c *ChatClient) buildHttpRequest(ctx context.Context, req *llm.RunModelRequest) (*http.Request, error) {
	var url string
	var body map[string]interface{}
	switch req.ModelName {
	case ModelGPT3dot5:
		fallthrough
	case ModelGPT4:
		url = urlChat
		body = c.buildGPTBody(req)
	case ModelDavinci:
		url = urlCompletions
		body = c.buildCompletionsBody(req)
	case ModelDavinciEdit:
		url = urlEdits
		body = c.buildEditsBody(req)
	default:
		return nil, fmt.Errorf("unknown model: %s", req.ModelName)
	}
	if c.Verbose {
		zap.L().Info("ChatClient build body", zap.Any("body", body))
	}
	bodyBytes, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("https://api.openai.com%s", url),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthKey))
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func (c *ChatClient) Execute(ctx context.Context, req *llm.RunModelRequest) (*llm.RunModelResponse, error) {
	httpReq, err := c.buildHttpRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	httpResponse, err := hc.Do(httpReq)
	if isTimeoutError(err) {
		zap.L().Warn("call openai timeout")
		return nil, fmt.Errorf("call openai timeout: %w", err)
	}
	if err != nil {
		zap.L().Error("call openai returns error", zap.Error(err))
		return nil, fmt.Errorf("call openai returns client error: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		defer httpResponse.Body.Close()
		if httpResponse.StatusCode == http.StatusServiceUnavailable || httpResponse.StatusCode == http.StatusBadGateway {
			return nil, fmt.Errorf("openai is down")
		}
		body, err := io.ReadAll(httpResponse.Body)
		zap.L().Error("call openai returns http error",
			zap.Int("status_code", httpResponse.StatusCode),
			zap.ByteString("body", body),
			zap.Error(err),
		)
		return nil, fmt.Errorf("call openai returns http %d error", httpResponse.StatusCode)
	}
	if req.Stream {
		if req.ModelName == ModelDavinciEdit {
			defer httpResponse.Body.Close()
			respBytes, err := io.ReadAll(httpResponse.Body)
			if err != nil {
				return nil, err
			}
			r := &completionsResp{}
			err = json.Unmarshal(respBytes, r)
			result := r.Choices[0].Text
			return &llm.RunModelResponse{
				Stream:       true,
				StreamResult: llm.NewStreamResponseFromString(result),
			}, nil
		} else {
			body := httpResponse.Body
			return &llm.RunModelResponse{
				Stream: true,
				StreamResult: &openaiStreamResponse{
					body:       body,
					reader:     bufio.NewReader(body),
					isChatResp: req.ModelName == ModelGPT3dot5 || req.ModelName == ModelGPT4,
				},
			}, nil
		}
	} else {
		defer httpResponse.Body.Close()
		respBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
		if c.Verbose {
			zap.L().Info("openai chat resp", zap.ByteString("resp", respBytes))
		}
		var result string
		if req.ModelName == ModelGPT3dot5 || req.ModelName == ModelGPT4 {
			r := &chatResp{}
			err = json.Unmarshal(respBytes, r)
			result = r.Choices[0].Message.Content
		} else {
			r := &completionsResp{}
			err = json.Unmarshal(respBytes, r)
			result = r.Choices[0].Text
		}
		return &llm.RunModelResponse{
			Result: result,
		}, nil
	}
}
