package firestore

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/option"
)

type Client struct {
	name   string
	opts   options
	client *auth.Client
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
	b := []byte(c.opts.credentials)

	config := &firebase.Config{ProjectID: c.opts.projectID}
	app, err := firebase.NewApp(ctx, config, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}
	c.client = client

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "custom_token":
		return c.CustomToken(ctx, meta, req.Data)
	case "verify_token":
		return c.VerifyToken(ctx, meta)
	}
	return nil, nil
}