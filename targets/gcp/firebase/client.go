package firestore

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
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

func (c *Client) CustomToken(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	token := ""
	if data != nil {
		claims := make(map[string]interface{})
		err := json.Unmarshal(data, &claims)
		if err != nil {
			return nil, err
		}
		if len(claims) == 0 {
			return nil, fmt.Errorf("body was set but data was missing claims")
		}
		token, err = c.client.CustomTokenWithClaims(ctx, meta.tokenID, claims)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		token, err = c.client.CustomToken(ctx, meta.tokenID)
		if err != nil {
			return nil, err
		}
	}
	b, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) VerifyToken(ctx context.Context, meta metadata) (*types.Response, error) {
	token, err := c.client.VerifyIDToken(ctx, meta.tokenID)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
