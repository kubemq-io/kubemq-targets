package mongodb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	monogOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	id                  = "_id"
	value               = "value"
	connectionURIFormat = "mongodb://%s:%s@%s/%s%s"
)

type Item struct {
	Key   string `bson:"_id"`
	Value string `bson:"value"`
}
type Client struct {
	name       string
	opts       options
	client     *mongo.Client
	collection *mongo.Collection
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
	c.client, err = c.getMongoDBClient(ctx)
	if err != nil {
		return fmt.Errorf("error in creating mongodb client: %s", err)
	}
	wc, err := c.getWriteConcernObject(c.opts.writeConcurrency)
	if err != nil {
		return fmt.Errorf("error in getting write concern object: %s", err)
	}
	rc, err := c.getReadConcernObject(c.opts.readConcurrency)
	if err != nil {
		return fmt.Errorf("error in getting read concern object: %s", err)
	}

	opts := monogOptions.Collection().SetWriteConcern(wc).SetReadConcern(rc)
	collection := c.client.Database(c.opts.database).Collection(c.opts.collection, opts)
	c.collection = collection
	return nil
}

func (c *Client) getWriteConcernObject(cn string) (*writeconcern.WriteConcern, error) {
	var wc *writeconcern.WriteConcern
	if cn != "" {
		if cn == "majority" {
			wc = writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(c.opts.operationTimeout))
		} else {
			w, err := strconv.Atoi(cn)
			wc = writeconcern.New(writeconcern.W(w), writeconcern.J(true), writeconcern.WTimeout(c.opts.operationTimeout))

			return wc, err
		}
	} else {
		wc = writeconcern.New(writeconcern.W(1), writeconcern.J(true), writeconcern.WTimeout(c.opts.operationTimeout))
	}

	return wc, nil
}

func (c *Client) getReadConcernObject(cn string) (*readconcern.ReadConcern, error) {
	switch cn {
	case "local":
		return readconcern.Local(), nil
	case "majority":
		return readconcern.Majority(), nil
	case "available":
		return readconcern.Available(), nil
	case "linearizable":
		return readconcern.Linearizable(), nil
	case "snapshot":
		return readconcern.Snapshot(), nil
	case "":
		return readconcern.Local(), nil
	}

	return nil, fmt.Errorf("readConcern %s not found", cn)
}

func (c *Client) getMongoDBClient(ctx context.Context) (*mongo.Client, error) {
	var uri string

	if c.opts.username != "" && c.opts.password != "" {
		uri = fmt.Sprintf(connectionURIFormat, c.opts.username, c.opts.password, c.opts.host, c.opts.database, c.opts.params)
	}
	clientOptions := monogOptions.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "get":
		return c.Get(ctx, meta)
	case "set":
		return c.Set(ctx, meta, req.Data)
	case "delete":
		return c.Delete(ctx, meta)

	}
	return nil, nil
}

func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	var result Item
	ctx, cancel := context.WithTimeout(ctx, c.opts.operationTimeout)
	defer cancel()

	filter := bson.M{id: meta.key}
	err := c.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("no data found for this key")
	}
	return types.NewResponse().
		SetData([]byte(result.Value)).
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) Set(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.operationTimeout)
	defer cancel()
	filter := bson.M{id: meta.key}
	update := bson.M{"$set": bson.M{id: meta.key, value: string(data)}}
	_, err := c.collection.UpdateOne(ctx, filter, update, monogOptions.Update().SetUpsert(true))
	if err != nil {
		return nil, fmt.Errorf("failed to set key %s: %s", meta.key, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.operationTimeout)
	defer cancel()
	filter := bson.M{id: meta.key}
	_, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete key %s: %s", meta.key, err)
	}

	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}
