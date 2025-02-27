package llm

import (
	"io"
)

type StreamResponse interface {
	io.Closer
	ReadString() (string, error)
}

func NewStreamResponseFromString(content string) StreamResponse {
	return &stringStreamResponse{content: content}
}

type stringStreamResponse struct {
	content string

	finished bool
}

func (r *stringStreamResponse) ReadString() (string, error) {
	if r.finished {
		return "", io.EOF
	}
	r.finished = true
	return r.content, nil
}

func (r *stringStreamResponse) Close() error {
	return nil
}
