package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/feature-vector/harbor/base/storage"
	"net/url"
	"time"
)

func NewAzureStorageProvider(accountName string, accountKey string) (storage.StorageProvider, error) {
	azureProvider := &azureStorageProvider{
		accountName: accountName,
		accountKey:  accountKey,
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	if err != nil {
		return nil, err
	}
	svc := azblob.NewServiceURL(*u, p)
	azureProvider.blobService = &svc
	return azureProvider, nil

}

type azureStorageProvider struct {
	accountName string
	accountKey  string
	blobService *azblob.ServiceURL
}

func (az *azureStorageProvider) UploadFile(ctx context.Context, req *storage.UploadFileRequest) (*storage.UploadFileResponse, error) {
	containerURL := az.blobService.NewContainerURL(req.BucketName)

	blobURL := containerURL.NewBlockBlobURL(req.FileName)

	_, err := blobURL.Upload(
		ctx,
		req.Content,
		azblob.BlobHTTPHeaders{ContentType: req.ContentType},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{},
		azblob.ImmutabilityPolicyOptions{},
	)
	return &storage.UploadFileResponse{}, err
}

func (az *azureStorageProvider) GetTempUrl(ctx context.Context, req *storage.GetTempUrlRequest) (*storage.GetTempUrlResponse, error) {
	filePath := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", az.accountName, req.BucketName, req.FileName)
	u, err := url.Parse(filePath)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now().UTC()
	credential, err := azblob.NewSharedKeyCredential(az.accountName, az.accountKey)
	if err != nil {
		return nil, err
	}

	sasqp, err := azblob.BlobSASSignatureValues{
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(req.Timeout),
		Permissions:   "racwd",
		ContainerName: req.BucketName,
		BlobName:      req.FileName,
		Protocol:      azblob.SASProtocolHTTPS,
	}.NewSASQueryParameters(credential)
	if err != nil {
		return nil, err
	}

	p := azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})
	snapParts := azblob.NewBlobURLParts(*u)
	snapParts.SAS = sasqp
	sburl := azblob.NewBlockBlobURL(snapParts.URL(), p)
	return &storage.GetTempUrlResponse{Url: sburl.String()}, nil
}

func (az *azureStorageProvider) DownloadFile(ctx context.Context, req *storage.DownloadFileRequest) (*storage.DownloadFileResponse, error) {
	panic("azureStorageProvider not implement DownloadFile")
}
