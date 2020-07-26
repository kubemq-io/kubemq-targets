package messaging

import (
	"context"
	"fmt"
	"strconv"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"firebase.google.com/go/messaging"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name             string
	opts             options
	client           *messaging.Client
	defult           *messaging.Message
	defaultMulticast *messaging.MulticastMessage
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
	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}
	c.client = client

	return nil

}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	m, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	if m.sendMulticast {
		b, err := c.client.SendMulticast(ctx, m.multicastMessage)
		if err != nil {
			return nil, err
		}
		r := types.NewResponse().
			SetMetadataKeyValue("SuccessCount", strconv.Itoa(b.SuccessCount)).
			SetMetadataKeyValue("FailureCount", strconv.Itoa(b.FailureCount))

		for _, res := range b.Responses {
			r.SetMetadataKeyValue(fmt.Sprintf("mesage_s%", res.MessageID), string(res))
		}
		return r
	}

	r, err := c.client.SendMulticast()(ctx, m.multicastMessage)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().SetMetadataKeyValue("result", r), nil

}
