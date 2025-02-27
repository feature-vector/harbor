package storage

import (
	"context"
	"io"
	"time"
)

type StorageProvider interface {
	UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error)
	GetTempUrl(ctx context.Context, req *GetTempUrlRequest) (*GetTempUrlResponse, error)
	DownloadFile(ctx context.Context, req *DownloadFileRequest) (*DownloadFileResponse, error)
}

type UploadFileRequest struct {
	BucketName  string
	FileName    string
	Content     io.ReadSeeker
	ContentType string
}

type UploadFileResponse struct {
}

type GetTempUrlRequest struct {
	BucketName string
	FileName   string
	Timeout    time.Duration
}

type GetTempUrlResponse struct {
	Url string
}

type DownloadFileRequest struct {
	BucketName string
	FileName   string
}

type DownloadFileResponse struct {
	Content     io.Reader
	ContentType string
}
