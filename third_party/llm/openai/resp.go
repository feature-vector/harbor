package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
)

type chatResp struct {
	Choices []*struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openaiStreamResponse struct {
	body     io.ReadCloser
	reader   *bufio.Reader
	finished bool

	isChatResp bool
}

type completionsResp struct {
	Choices []*struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type chatStreamResp struct {
	Choices []*struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (r *openaiStreamResponse) ReadString() (string, error) {
	prefix := []byte("data: ")
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	line = bytes.TrimSpace(line)
	if !bytes.HasPrefix(line, prefix) {
		return "", nil
	}
	line = bytes.TrimPrefix(line, prefix)
	if string(line) == "[DONE]" {
		return "", io.EOF
	}
	text := ""
	if r.isChatResp {
		r := &chatStreamResp{}
		err = json.Unmarshal(line, r)
		if err != nil {
			return "", err
		}
		if len(r.Choices) == 0 {
			zap.L().Error("openaiStreamResponse ReadString choices is null", zap.ByteString("line", line))
		} else {
			text = r.Choices[0].Delta.Content
		}
	} else {
		r := &completionsResp{}
		err = json.Unmarshal(line, r)
		if err != nil {
			return "", err
		}
		if len(r.Choices) == 0 {
			zap.L().Error("openaiStreamResponse ReadString choices is null", zap.ByteString("line", line))
		} else {
			text = r.Choices[0].Text
		}
	}
	return text, nil
}

func (r *openaiStreamResponse) Close() error {
	return r.body.Close()
}
