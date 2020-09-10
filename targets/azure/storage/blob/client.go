package blob

import (
	"context"
	"encoding/json"
	"errors"
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
		return err
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
	}
	return nil, errors.New("invalid method type")
}

func (c *Client) upload(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {

	URL, err := url.Parse(meta.serviceUrl)

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
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

	containerURL := azblob.NewContainerURL(*URL, c.pipeLine)
	blobURL := containerURL.NewBlobURL(meta.fileName)
	r, err := blobURL.Download(ctx, meta.offset, meta.count, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {

	URL, err := url.Parse(meta.serviceUrl)

	containerURL := azblob.NewContainerURL(*URL, c.pipeLine)
	blobURL := containerURL.NewBlobURL(meta.fileName)
	r, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionType(meta.deleteSnapshotsOptionType), azblob.BlobAccessConditions{})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}
