package elasticsearch

import (
	"context"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	signer "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/kubemq-io/kubemq-targets/config"
)

type Client struct {
	log    *logger.Logger
	opts   options
	signer *signer.Signer
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

func (c *Client) Stop() error {
	return nil
}
