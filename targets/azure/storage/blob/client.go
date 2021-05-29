package blob

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"net/url"
	"time"
)

type Client struct {
	log      *logger.Logger
	opts     options
	pipeLine pipeline.Pipeline
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}

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
	c.pipeLine = azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			Policy:        c.opts.policy,                           // Use exponential backoff as opposed to linear
			MaxTries:      c.opts.maxTries,                         // Try at most x times to perform the operation (set to 1 to disable retries)
			TryTimeout:    time.Millisecond * c.opts.tryTimeout,    // Maximum time allowed for any single try
			RetryDelay:    time.Millisecond * c.opts.retryDelay,    // Backoff amount for each retry (exponential or linear)
			MaxRetryDelay: time.Millisecond * c.opts.maxRetryDelay, // Max delay between retries
		},
	})

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
	uploadFileOption := azblob.UploadToBlockBlobOptions{
		BlockSize:   meta.blockSize,
		Parallelism: meta.parallelism}
	if len(meta.blobMetadata) > 0 {
		uploadFileOption.Metadata = meta.blobMetadata
	}
	_, err = azblob.UploadBufferToBlockBlob(ctx, data, blobURL, uploadFileOption)

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
	defer func() {
		_ = bodyStream.Close()
	}()
	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return nil, err
	}
	b := downloadedData.Bytes()
	m := downloadResponse.NewMetadata()
	if len(m) > 0 {
		jsonString, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return types.NewResponse().
				SetData(b).
				SetMetadataKeyValue("blob_metadata", fmt.Sprintf("%s", jsonString)).
				SetMetadataKeyValue("result", "ok"),
			nil
	}
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

func (c *Client) Stop() error {
	return nil
}
