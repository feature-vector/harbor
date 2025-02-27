package claude

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

type chatResp struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Model      string `json:"model"`
}

type claudeStreamResponse struct {
	body     io.ReadCloser
	reader   *bufio.Reader
	finished bool
}

type chatStreamResp struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Model      string `json:"model"`
}

func (r *claudeStreamResponse) ReadString() (string, error) {
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
	resp := &chatStreamResp{}
	err = json.Unmarshal(line, resp)
	if err != nil {
		return "", err
	}
	if resp.StopReason != "" {
		return "", io.EOF
	}
	return resp.Completion, nil
}

func (r *claudeStreamResponse) Close() error {
	return r.body.Close()
}
