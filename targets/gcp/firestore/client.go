package google

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/targets"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"google.golang.org/api/iterator"
)

type Client struct {
	name   string
	opts   options
	client *firestore.Client
	log    *logger.Logger
	target targets.Target
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}

func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	client, err := firestore.NewClient(ctx, c.opts.projectID)
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
	case "documents_all":
		return c.documentAll(ctx, meta)
	case "document_key":
		return c.documentKey(ctx, meta)
	case "add":
		return c.add(ctx, meta, req.Data)
	case "delete_document_key":
		return c.deleteDocument(ctx,meta)
	}
	return nil, nil
}

func (c *Client) add(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", "failed to parse data as map"), nil
	}
	_, _, err = c.client.Collection(meta.key).Add(ctx, m)
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", err.Error()), nil
	}
	return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("error", "false").
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) documentAll(ctx context.Context, meta metadata) (*types.Response, error) {
	iter := c.client.Collection(meta.key).Documents(ctx)
	var retData []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return types.NewResponse().
				SetMetadataKeyValue("collection", meta.key).
				SetMetadataKeyValue("error", "true").
				SetMetadataKeyValue("message", fmt.Sprintf(err.Error())), nil
		}
		retData = append(retData, doc.Data())
	}
	if len(retData) <= 0 {
		return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", "no data found for this key"), nil
	}
	byte, err := json.Marshal(retData)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(byte).
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("collection", meta.key), nil
}

func (c *Client) documentKey(ctx context.Context, meta metadata) (*types.Response, error) {
	obj, err := c.client.Collection(meta.key).Doc(meta.item).Get(ctx)
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("item", meta.item).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", err.Error()), nil
	}
	byte, err := json.Marshal(obj.Data())
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("item", meta.item).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", err.Error()), nil
	}
	return types.NewResponse().
		SetData(byte).
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("item", meta.item).
		SetMetadataKeyValue("collection", meta.key), nil
}

func (c *Client) deleteDocument(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.Collection(meta.key).Doc(meta.item).Delete(ctx)
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
			SetMetadataKeyValue("item", meta.item).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", err.Error()), nil
	}
	return types.NewResponse().
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("item", meta.item).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("collection", meta.key), nil
}


func (c *Client) list(ctx context.Context) (*types.Response, error) {
	var collections []string
	it := c.client.Collections(ctx)
	for {
		collection, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection.ID)
	}
	if len(collections)<=0 {
		return types.NewResponse().
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", "no collections found for this project"), nil
	}
	b,err:=json.Marshal(collections)
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", err.Error()), nil
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}
