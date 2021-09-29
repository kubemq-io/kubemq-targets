package firebase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

func (c *Client) customToken(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
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
		token, err = c.clientAuth.CustomTokenWithClaims(ctx, meta.tokenID, claims)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		token, err = c.clientAuth.CustomToken(ctx, meta.tokenID)
		if err != nil {
			return nil, err
		}
	}
	b, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) verifyToken(ctx context.Context, meta metadata) (*types.Response, error) {
	token, err := c.clientAuth.VerifyIDToken(ctx, meta.tokenID)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
