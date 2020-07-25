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
	"strconv"
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
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) retrieveUser(ctx context.Context, meta metadata) (*types.Response, error) {
	var b []byte
	switch meta.retrieveBy {
	case "by_uid":
		u, err := c.client.GetUser(ctx, meta.uid)
		if err != nil {
			return nil, err
		}
		b, err = json.Marshal(u)
		if err != nil {
			return nil, err
		}
	case "by_email":
		u, err := c.client.GetUserByEmail(ctx, meta.email)
		if err != nil {
			return nil, err
		}
		b, err = json.Marshal(u)
		if err != nil {
			return nil, err
		}
	case "by_phone":
		u, err := c.client.GetUserByPhoneNumber(ctx, meta.phone)
		if err != nil {
			return nil, err
		}
		b, err = json.Marshal(u)
		if err != nil {
			return nil, err
		}
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) createUser(ctx context.Context, data []byte) (*types.Response, error) {
	p, err := getCreateData(data)
	if err != nil {
		return nil, err
	}
	u, err := c.client.CreateUser(ctx, p)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) updateUser(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	p, err := getUserUpdateData(data)
	if err != nil {
		return nil, err
	}
	u, err := c.client.UpdateUser(ctx, meta.uid, p)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) deleteUser(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.client.DeleteUser(ctx, meta.uid)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteMultipleUser(ctx context.Context, data []byte) (*types.Response, error) {
	var l []string
	err := json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}
	r, err := c.client.DeleteUsers(ctx, l)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}


func getUserUpdateData(data []byte) (*auth.UserToUpdate, error) {
	u := &auth.UserToUpdate{}
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return u, err
	}
	for k, v := range m {
		switch k {
		case "custom_claims":
			c := make(map[string]interface{})
			err := json.Unmarshal(data, &c)
			if err != nil {
				return u, err
			}
			u.CustomClaims(c)
		case "disabled":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return u, err
			}
			u.Disabled(b)
		case "display_name":
			u.DisplayName(fmt.Sprintf("%s", v))
		case "email":
			u.Email(fmt.Sprintf("%s", v))
		case "email_verified":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return u, err
			}
			u.EmailVerified(b)
		case "password":
			u.Password(fmt.Sprintf("%s", v))
		case "phone_number":
			u.PhoneNumber(fmt.Sprintf("%s", v))
		case "photo_url":
			u.PhotoURL(fmt.Sprintf("%s", v))
		}
	}
	return u, nil
}

func getCreateData(data []byte) (*auth.UserToCreate, error) {
	u := &auth.UserToCreate{}
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return u, err
	}
	for k, v := range m {
		switch k {
		case "disabled":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return u, err
			}
			u.Disabled(b)
		case "display_name":
			u.DisplayName(fmt.Sprintf("%s", v))
		case "email":
			u.Email(fmt.Sprintf("%s", v))
		case "email_verified":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return u, err
			}
			u.EmailVerified(b)
		case "password":
			u.Password(fmt.Sprintf("%s", v))
		case "phone_number":
			u.PhoneNumber(fmt.Sprintf("%s", v))
		case "photo_url":
			u.PhotoURL(fmt.Sprintf("%s", v))
		case "local_id":
			u.UID(fmt.Sprintf("%s", v))
		}
	}
	return u, nil
}
