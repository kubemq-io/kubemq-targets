package files

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-file-go/azfile"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"net/url"
	"time"
)

type Client struct {
	name     string
	opts     options
	pipeLine pipeline.Pipeline
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azfile.NewSharedKeyCredential(c.opts.storageAccount, c.opts.storageAccessKey)
	if err != nil {
		return fmt.Errorf("failed to create shared key credential on error %s , please check storage access key and acccount are correct", err.Error())
	}
	c.pipeLine = azfile.NewPipeline(credential, azfile.PipelineOptions{
		Retry: azfile.RetryOptions{
			Policy:        c.opts.policy,                           // Use exponential or linear
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
	case "create":
		return c.create(ctx, meta)
	case "upload":
		return c.upload(ctx, meta, req.Data)
	case "get":
		return c.get(ctx, meta)
	case "delete":
		return c.delete(ctx, meta)
	}
	return nil, errors.New("invalid method type")
}

func (c *Client) create(ctx context.Context, meta metadata) (*types.Response, error) {

	URL, err := url.Parse(meta.serviceUrl)
	if err != nil {
		return nil, err
	}
	fileURL := azfile.NewFileURL(*URL, c.pipeLine)
	if len(meta.fileMetadata) > 0 {
		_, err = fileURL.Create(ctx, meta.size, azfile.FileHTTPHeaders{}, meta.fileMetadata)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = fileURL.Create(ctx, meta.size, azfile.FileHTTPHeaders{}, nil)
		if err != nil {
			return nil, err
		}
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) upload(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {

	if data == nil {
		return nil, errors.New("missing data to upload")
	}
	URL, err := url.Parse(meta.serviceUrl)
	if err != nil {
		return nil, err
	}
	fileURL := azfile.NewFileURL(*URL, c.pipeLine)
	uploadFileOption := azfile.UploadToAzureFileOptions{
		RangeSize:   meta.rangeSize,
		Parallelism: meta.parallelism}
	if len(meta.fileMetadata) > 0 {
		uploadFileOption.Metadata = meta.fileMetadata
	}
	err = azfile.UploadBufferToAzureFile(ctx, data, fileURL, uploadFileOption)

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
	fileURL := azfile.NewFileURL(*URL, c.pipeLine)
	downloadResponse, err := fileURL.Download(ctx, meta.offset, meta.count, false)
	if err != nil {
		return nil, err
	}
	bodyStream := downloadResponse.Body(azfile.RetryReaderOptions{MaxRetryRequests: meta.maxRetryRequests})
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
				SetMetadataKeyValue("file_metadata", fmt.Sprintf("%s", jsonString)).
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
	fileURL := azfile.NewFileURL(*URL, c.pipeLine)
	_, err = fileURL.Delete(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}
