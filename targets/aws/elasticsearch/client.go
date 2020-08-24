package elasticsearch

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	signer "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type Client struct {
	name   string
	opts   options
	signer *signer.Signer
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

	signer := signer.NewSigner(credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token))
	c.signer = signer

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	reader := strings.NewReader(meta.json)
	request, err := http.NewRequestWithContext(ctx, meta.method, meta.endpoint, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	_, err = c.signer.Sign(request, reader, meta.service, meta.region, time.Now())
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
