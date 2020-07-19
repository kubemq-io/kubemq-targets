package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	name   string
	client *resty.Client
	opts   options
}

func New() *Client {
	return &Client{}
}

func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client = resty.New()
	c.client.SetDoNotParseResponse(true)
	switch c.opts.authType {
	case "basic":
		c.client.SetBasicAuth(c.opts.username, c.opts.password)
	case "auth_token":
		c.client.SetAuthToken(c.opts.token)
	}
	c.client.SetHeaders(c.opts.defaultHeaders)
	if c.opts.proxy != "" {
		c.client.SetProxy(c.opts.proxy)
	}
	if c.opts.rootCertificate != "" {
		c.client.SetRootCertificateFromString(c.opts.rootCertificate)
	}
	if c.opts.clientPrivateKey != "" && c.opts.clientPublicKey != "" {
		cert, err := tls.X509KeyPair([]byte(c.opts.clientPublicKey), []byte(c.opts.clientPrivateKey))
		if err != nil {
			return fmt.Errorf("error loading client certificate: %w", err)
		}
		c.client.SetCertificates(cert)
	}
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	httpReq := c.client.R().
		SetHeaders(meta.headers).
		SetContext(ctx)

	if req.Data != nil {
		httpReq.SetBody(req.Data)
	}
	httpReq.URL = meta.url
	httpReq.Method = strings.ToUpper(meta.method)
	resp, err := httpReq.Send()
	if err != nil {
		return nil, err
	}
	tr, err := newResultFromHttpResponse(resp.RawResponse)
	if err != nil {
		return nil, err
	}
	return tr, resp.RawResponse.Body.Close()
}

func newResultFromHttpResponse(hr *http.Response) (*types.Response, error) {
	resp := types.NewResponse().
		SetMetadataKeyValue("code", fmt.Sprintf("%d", hr.StatusCode)).
		SetMetadataKeyValue("status", hr.Status)

	headers := types.NewMetadata()
	for name, values := range hr.Header {
		headers[name] = strings.Join(values, ",")
	}
	resp.SetMetadataKeyValue("headers", headers.String())
	var err error
	if hr.Body == nil {
		return resp, nil
	}
	resp.Data, err = ioutil.ReadAll(hr.Body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
