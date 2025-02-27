package google

import (
	googleStorage "cloud.google.com/go/storage"
	"context"
	"github.com/feature-vector/harbor/base/storage"
	"io"
	"time"
)

type googleStorageProvider struct {
	storageClient *googleStorage.Client
}

func NewGoogleStorageProvider(client *googleStorage.Client) storage.StorageProvider {
	return &googleStorageProvider{
		storageClient: client,
	}
}

func (gp *googleStorageProvider) UploadFile(ctx context.Context, req *storage.UploadFileRequest) (*storage.UploadFileResponse, error) {
	wc := gp.storageClient.Bucket(req.BucketName).Object(req.FileName).NewWriter(ctx)
	wc.ChunkSize = 0
	wc.ContentType = req.ContentType
	_, err := io.Copy(wc, req.Content)
	if err != nil {
		return nil, err
	}
	err = wc.Close()
	if err != nil {
		return nil, err
	}
	return &storage.UploadFileResponse{}, nil
}

func (gp *googleStorageProvider) GetTempUrl(ctx context.Context, req *storage.GetTempUrlRequest) (*storage.GetTempUrlResponse, error) {
	opts := &googleStorage.SignedURLOptions{
		Scheme:  googleStorage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(req.Timeout),
	}
	u, err := gp.storageClient.Bucket(req.BucketName).SignedURL(req.FileName, opts)
	if err != nil {
		return nil, err
	}
	return &storage.GetTempUrlResponse{Url: u}, nil
}

func (gp *googleStorageProvider) DownloadFile(ctx context.Context, req *storage.DownloadFileRequest) (*storage.DownloadFileResponse, error) {
	reader, err := gp.storageClient.Bucket(req.BucketName).Object(req.FileName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return &storage.DownloadFileResponse{
		Content:     reader,
		ContentType: reader.Attrs.ContentType,
	}, nil
}
