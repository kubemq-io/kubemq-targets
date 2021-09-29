package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"firebase.google.com/go/v4/auth"
	"github.com/kubemq-io/kubemq-targets/types"
	"google.golang.org/api/iterator"
)


func (c *Client) createUser(ctx context.Context, data []byte) (*types.Response, error) {
	p, err := getCreateData(data)
	if err != nil {
		return nil, err
	}
	u, err := c.clientAuth.CreateUser(ctx, p)
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

func (c *Client) retrieveUser(ctx context.Context, meta metadata) (*types.Response, error) {
	var b []byte
	switch meta.retrieveBy {
	case "by_uid":
		u, err := c.clientAuth.GetUser(ctx, meta.uid)
		if err != nil {
			return nil, err
		}
		b, err = json.Marshal(u)
		if err != nil {
			return nil, err
		}
	case "by_email":
		u, err := c.clientAuth.GetUserByEmail(ctx, meta.email)
		if err != nil {
			return nil, err
		}
		b, err = json.Marshal(u)
		if err != nil {
			return nil, err
		}
	case "by_phone":
		u, err := c.clientAuth.GetUserByPhoneNumber(ctx, meta.phone)
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

func (c *Client) updateUser(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	p, err := getUpdateData(data)
	if err != nil {
		return nil, err
	}
	u, err := c.clientAuth.UpdateUser(ctx, meta.uid, p)
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
	err := c.clientAuth.DeleteUser(ctx, meta.uid)
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
	r, err := c.clientAuth.DeleteUsers(ctx, l)
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

func (c *Client) listAllUsers(ctx context.Context) (*types.Response, error) {
	var users []*auth.ExportedUserRecord
	iter := c.clientAuth.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	b, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func getUpdateData(data []byte) (*auth.UserToUpdate, error) {
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
			b, err := strconv.ParseBool(fmt.Sprintf("%v", v))
			if err != nil {
				return u, err
			}
			u.Disabled(b)
		case "display_name":
			u.DisplayName(fmt.Sprintf("%v", v))
		case "email":
			u.Email(fmt.Sprintf("%v", v))
		case "email_verified":
			b, err := strconv.ParseBool(fmt.Sprintf("%v", v))
			if err != nil {
				return u, err
			}
			u.EmailVerified(b)
		case "password":
			u.Password(fmt.Sprintf("%v", v))
		case "phone_number":
			u.PhoneNumber(fmt.Sprintf("%v", v))
		case "photo_url":
			u.PhotoURL(fmt.Sprintf("%v", v))
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
			b, err := strconv.ParseBool(fmt.Sprintf("%v", v))
			if err != nil {
				return u, err
			}
			u.Disabled(b)
		case "display_name":
			u.DisplayName(fmt.Sprintf("%v", v))
		case "email":
			u.Email(fmt.Sprintf("%v", v))
		case "email_verified":
			b, err := strconv.ParseBool(fmt.Sprintf("%v", v))
			if err != nil {
				return u, err
			}
			u.EmailVerified(b)
		case "password":
			u.Password(fmt.Sprintf("%s", v))
		case "phone_number":
			u.PhoneNumber(fmt.Sprintf("%v", v))
		case "photo_url":
			u.PhotoURL(fmt.Sprintf("%v", v))
		case "local_id":
			u.UID(fmt.Sprintf("%v", v))
		}
	}
	return u, nil
}
