package mongodb

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	monogOptions "go.mongodb.org/mongo-driver/mongo/options"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	id    = "_id"
	value = "value"
)

type Item struct {
	Key   string `bson:"_id"`
	Value string `bson:"value"`
}
type Client struct {
	log        *logger.Logger
	opts       options
	client     *mongo.Client
	collection *mongo.Collection
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
	c.client, err = c.getMongoDBClient(ctx)
	if err != nil {
		return fmt.Errorf("error in creating mongodb client: %s", err)
	}

	collection := c.client.Database(c.opts.database).Collection(c.opts.collection)
	c.collection = collection
	return nil
}

func (c *Client) getMongoDBClient(ctx context.Context) (*mongo.Client, error) {
	opts := monogOptions.Client().ApplyURI(c.opts.url)
	client, err := mongo.Connect(ctx, opts)
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
	case "get_by_key":
		return c.Get(ctx, meta)
	case "set_by_key":
		return c.Set(ctx, meta, req.Data)
	case "delete_by_key":
		return c.Delete(ctx, meta)
	case "find":
		return c.FindOne(ctx, meta)
	case "find_many":
		return c.Find(ctx, meta)
	case "insert":
		return c.Insert(ctx, req.Data)
	case "insert_many":
		return c.InsertMany(ctx, req.Data)
	case "update":
		return c.UpdateOne(ctx, meta, req.Data)
	case "update_many":
		return c.UpdateMany(ctx, meta, req.Data)
	case "delete":
		return c.DeleteOne(ctx, meta)
	case "delete_many":
		return c.DeleteMany(ctx, meta)
	case "aggregate":
		return c.Aggregate(ctx, req.Data)
	case "distinct":
		return c.Distinct(ctx, meta)
	}
	return nil, nil
}

func (c *Client) FindOne(ctx context.Context, meta metadata) (*types.Response, error) {
	if len(meta.filter) == 0 {
		return nil, fmt.Errorf("find one document filter is invalid")
	}
	result := map[string]interface{}{}
	err := c.collection.FindOne(ctx, meta.filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("find one error, %s", err.Error())
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("find one json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Find(ctx context.Context, meta metadata) (*types.Response, error) {
	if len(meta.filter) == 0 {
		return nil, fmt.Errorf("find documents filter is invalid")
	}
	results := []map[string]interface{}{}
	cursor, err := c.collection.Find(ctx, meta.filter)
	if err != nil {
		return nil, fmt.Errorf("find error, %s", err.Error())
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("find results parsing error, %s", err.Error())
	}
	data, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("find json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Insert(ctx context.Context, reqData []byte) (*types.Response, error) {
	var doc interface{}

	err := json.Unmarshal(reqData, &doc)
	if err != nil {
		return nil, fmt.Errorf("insert document json parsing error, %s", err.Error())
	}
	result, err := c.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("insert error, %s", err.Error())
	}
	data, err := json.Marshal(result.InsertedID)
	if err != nil {
		return nil, fmt.Errorf("insert result json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) InsertMany(ctx context.Context, reqData []byte) (*types.Response, error) {
	var docs []interface{}
	err := json.Unmarshal(reqData, &docs)
	if err != nil {
		return nil, fmt.Errorf("insert many documents json parsing error, %s", err.Error())
	}

	results, err := c.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("insert many error, %s", err.Error())
	}
	data, err := json.Marshal(results.InsertedIDs)
	if err != nil {
		return nil, fmt.Errorf("insert many results json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) UpdateOne(ctx context.Context, meta metadata, reqData []byte) (*types.Response, error) {
	var doc interface{}
	err := json.Unmarshal(reqData, &doc)
	if err != nil {
		return nil, fmt.Errorf("update one document json parsing error, %s", err.Error())
	}
	update := bson.M{"$set": &doc}
	result, err := c.collection.UpdateOne(ctx, meta.filter, update, monogOptions.Update().SetUpsert(meta.setUpsert))
	if err != nil {
		return nil, fmt.Errorf("update one document error, %s", err.Error())
	}
	data, err := json.Marshal(result.UpsertedID)
	if err != nil {
		return nil, fmt.Errorf("update one document result json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) UpdateMany(ctx context.Context, meta metadata, reqData []byte) (*types.Response, error) {
	var doc interface{}
	err := json.Unmarshal(reqData, &doc)
	if err != nil {
		return nil, fmt.Errorf("update many document json parsing error, %s", err.Error())
	}

	update := bson.M{"$set": &doc}
	result, err := c.collection.UpdateMany(ctx, meta.filter, update, monogOptions.Update().SetUpsert(meta.setUpsert))
	if err != nil {
		return nil, fmt.Errorf("update many documents error, %s", err.Error())
	}
	data, err := json.Marshal(result.UpsertedID)
	if err != nil {
		return nil, fmt.Errorf("update many documents result json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) DeleteOne(ctx context.Context, meta metadata) (*types.Response, error) {
	if len(meta.filter) == 0 {
		return nil, fmt.Errorf("delete one document filter is invalid")
	}
	result, err := c.collection.DeleteOne(ctx, meta.filter)
	if err != nil {
		return nil, fmt.Errorf("delete one document error, %s", err.Error())
	}
	return types.NewResponse().
		SetMetadataKeyValue("deleted_count", fmt.Sprintf("%d", result.DeletedCount)).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) DeleteMany(ctx context.Context, meta metadata) (*types.Response, error) {
	if len(meta.filter) == 0 {
		return nil, fmt.Errorf("delete many documents filter is invalid")
	}
	result, err := c.collection.DeleteMany(ctx, meta.filter)
	if err != nil {
		return nil, fmt.Errorf("delete many documents error, %s", err.Error())
	}
	return types.NewResponse().
		SetMetadataKeyValue("deleted_count", fmt.Sprintf("%d", result.DeletedCount)).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Aggregate(ctx context.Context, reqData []byte) (*types.Response, error) {
	var pipeline interface{}
	err := json.Unmarshal(reqData, &pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate pipeline json parsing error, %s", err.Error())
	}
	results := []map[string]interface{}{}
	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error, %s", err.Error())
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("aggregate results parsing error, %s", err.Error())
	}
	data, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("aggregate json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Distinct(ctx context.Context, meta metadata) (*types.Response, error) {
	if meta.fieldName == "" {
		return nil, fmt.Errorf("distinct field name missing")
	}
	if len(meta.filter) == 0 {
		return nil, fmt.Errorf("distinct filter is invalid")
	}

	results, err := c.collection.Distinct(ctx, meta.fieldName, meta.filter)
	if err != nil {
		return nil, fmt.Errorf("distinct error, %s", err.Error())
	}
	data, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("distinct json parsing error, %s", err.Error())
	}
	return types.NewResponse().
		SetData(data).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	if meta.key == "" {
		return nil, fmt.Errorf("get by key error, invalid key")
	}
	var result Item

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
	if meta.key == "" {
		return nil, fmt.Errorf("set by key error, invalid key")
	}
	if data == nil {
		return nil, fmt.Errorf("set by key error, invalid document")
	}
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
	if meta.key == "" {
		return nil, fmt.Errorf("delete by key error, invalid key")
	}

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

func (c *Client) Stop() error {
	return c.client.Disconnect(context.Background())
}
