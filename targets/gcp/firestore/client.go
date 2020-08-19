package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	name   string
	opts   options
	client *firestore.Client
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
	client, err := firestore.NewClient(ctx, c.opts.projectID, option.WithCredentialsJSON(b))
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
		return c.deleteDocument(ctx, meta)
	}
	return nil, fmt.Errorf("invalid method type")
}

func (c *Client) add(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data as map")
	}
	_, _, err = c.client.Collection(meta.key).Add(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("collection", meta.key).
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
			return nil, err
		}
		retData = append(retData, doc.Data())
	}
	if len(retData) <= 0 {
		return nil, fmt.Errorf("no data found for this key")
	}
	data, err := json.Marshal(retData)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("collection", meta.key), nil
}

func (c *Client) documentKey(ctx context.Context, meta metadata) (*types.Response, error) {
	obj, err := c.client.Collection(meta.key).Doc(meta.item).Get(ctx)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(obj.Data())
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("item", meta.item).
		SetMetadataKeyValue("collection", meta.key), nil
}

func (c *Client) deleteDocument(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.Collection(meta.key).Doc(meta.item).Delete(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
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
	if len(collections) <= 0 {
		return nil, fmt.Errorf("no collections found for this project")
	}
	data, err := json.Marshal(collections)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(data).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) CloseClient() error {
	return c.client.Close()
}
