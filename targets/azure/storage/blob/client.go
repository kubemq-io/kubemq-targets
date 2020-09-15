package blob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"net/url"
)

type Client struct {
	name     string
	opts     options
	pipeLine pipeline.Pipeline
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(c.opts.storageAccount, c.opts.storageAccessKey)
	if err != nil {
		return fmt.Errorf("failed to create shared key credential on error %s , please check storage access key and acccount are correct", err.Error())
	}
	c.pipeLine = azblob.NewPipeline(credential, azblob.PipelineOptions{})

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "upload":
		return c.upload(ctx, meta, req.Data)
	case "get":
		return c.get(ctx, meta)
	case "delete":
		return c.delete(ctx, meta)
	}
	return nil, errors.New("invalid method type")
}

func (c *Client) upload(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {

	if data == nil {
		return nil, errors.New("missing data to upload")
	}
	URL, err := url.Parse(meta.serviceUrl)
	if err != nil {
		return nil, err
	}
	containerURL := azblob.NewContainerURL(*URL, c.pipeLine)
	blobURL := containerURL.NewBlockBlobURL(meta.fileName)
	_, err = azblob.UploadBufferToBlockBlob(ctx, data, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   meta.blockSize,
		Parallelism: meta.parallelism})

	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) get(ctx context.Context, meta metadata) (*types.Response, error) {

	URL, err := url.Parse(meta.serviceUrl)
	if err != nil {
		return nil, err
	}
	containerURL := azblob.NewContainerURL(*URL, c.pipeLine)
	blobURL := containerURL.NewBlobURL(meta.fileName)
	downloadResponse, err := blobURL.Download(ctx, meta.offset, meta.count, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: meta.maxRetryRequests})

	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return nil, err
	}
	b := downloadedData.Bytes()
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {

	URL, err := url.Parse(meta.serviceUrl)
	if err != nil {
		return nil, err
	}
	containerURL := azblob.NewContainerURL(*URL, c.pipeLine)
	blobURL := containerURL.NewBlobURL(meta.fileName)
	_, err = blobURL.Delete(ctx, meta.deleteSnapshotsOptionType, azblob.BlobAccessConditions{})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}
