package firebase

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"firebase.google.com/go/v4/messaging"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/option"
)

type Client struct {
	name            string
	opts            options
	clientAuth      *auth.Client
	dbClient        *db.Client
	messagingClient *messaging.Client
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
	b := []byte(c.opts.credentials)

	config := &firebase.Config{ProjectID: c.opts.projectID, DatabaseURL: c.opts.dbURL}
	app, err := firebase.NewApp(ctx, config, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	if c.opts.authClient {
		client, err := app.Auth(ctx)
		if err != nil {
			return err
		}
		c.clientAuth = client
	}
	if c.opts.dbClient {
		client, err := app.Database(ctx)
		if err != nil {
			return err
		}
		c.dbClient = client
	}

	if c.opts.messagingClient {
		c.messagingClient, err = app.Messaging(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "custom_token":
		return c.customToken(ctx, meta, req.Data)
	case "verify_token":
		return c.verifyToken(ctx, meta)
	case "retrieve_user":
		return c.retrieveUser(ctx, meta)
	case "create_user":
		return c.createUser(ctx, req.Data)
	case "update_user":
		return c.updateUser(ctx, meta, req.Data)
	case "delete_user":
		return c.deleteUser(ctx, meta)
	case "delete_multiple_users":
		return c.deleteMultipleUser(ctx, req.Data)
	case "list_users":
		return c.listAllUsers(ctx)
	case "get_db":
		return c.dbGet(ctx, meta)
	case "update_db":
		return c.dbUpdate(ctx, meta, req.Data)
	case "set_db":
		return c.dbSet(ctx, meta, req.Data)
	case "delete_db":
		return c.dbDelete(ctx, meta)
	case "send_message":
		return c.sendMessage(ctx, req, c.opts)
	case "send_multi":
		return c.sendMessageMulti(ctx, req, c.opts)
	}
	return nil, errors.New("invalid method type")
}
