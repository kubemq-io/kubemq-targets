package firebase

import (
	"context"
	"encoding/json"

	"github.com/kubemq-hub/kubemq-targets/types"
)

func (c *Client) dbGet(ctx context.Context, meta metadata) (*types.Response, error) {
	ref := c.dbClient.NewRef(meta.refPath)
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		return nil, err
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) dbSet(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	ref := c.dbClient.NewRef(meta.refPath)
	var dat map[string]interface{}
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	if meta.childRefPath != "" {
		childRef := ref.Child(meta.childRefPath)
		if err := childRef.Set(ctx, dat); err != nil {
			return nil, err
		}
	} else {
		if err := ref.Set(ctx, dat); err != nil {
			return nil, err
		}
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) dbUpdate(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	ref := c.dbClient.NewRef(meta.refPath)
	var dat map[string]interface{}
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	if meta.childRefPath != "" {
		childRef := ref.Child(meta.childRefPath)
		if err := childRef.Update(ctx, dat); err != nil {
			return nil, err
		}
	} else {
		if err := ref.Update(ctx, dat); err != nil {
			return nil, err
		}
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) dbPush(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	ref := c.dbClient.NewRef(meta.refPath)
	var dat map[string]interface{}
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	var b []byte
	if meta.childRefPath != "" {
		childRef := ref.Child(meta.childRefPath)
		newChildRef, err := childRef.Push(ctx, nil)
		if err != nil {
			return nil, err
		}
		if err := newChildRef.Set(ctx, dat); err != nil {
			return nil, err
		}
		b, err = json.Marshal(newChildRef.Key)
		if err != nil {
			return nil, err
		}
	} else {
		newRef, err := ref.Push(ctx, nil)
		if err != nil {
			return nil, err
		}
		if err := newRef.Set(ctx, dat); err != nil {
			return nil, err
		}
		b, err = json.Marshal(newRef.Key)
		if err != nil {
			return nil, err
		}

	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
