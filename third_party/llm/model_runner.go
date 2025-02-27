package llm

import "context"

type ModelRunner interface {
	Execute(ctx context.Context, req *RunModelRequest) (*RunModelResponse, error)
}

type RunModelRequest struct {
	ModelName string
	Input     string
	Prompt    string
	History   [][]string
	Stream    bool
}

type RunModelResponse struct {
	Result       string
	Stream       bool
	StreamResult StreamResponse
}
