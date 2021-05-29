package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getUserStruct() map[string]interface{} {
	m := make(map[string]interface{})
	m["disabled"] = false
	m["display_name"] = "te3st"
	m["email"] = "test315413a@test.com"
	m["email_verified"] = true
	m["password"] = "testPassword"
	m["phone_number"] = "+12343678912"
	m["photo_url"] = "https://kubemq.io/wp-content/uploads/2018/11/24350KubeMQ_clean.png"

	return m
}

func TestClient_createUser(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	u := getUserStruct()
	b, err := json.Marshal(&u)
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "gcp-firebase",
		Kind: "gcp.firebase",
		Properties: map[string]string{
			"project_id":       dat.projectID,
			"credentials":      dat.cred,
			"auth_client":      "true",
			"messaging_client": "false",
			"db_client":        "false",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		{
			name: "create user-valid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_user").
				SetData(b),
			wantErr: false,
		},
		{
			name: "create user-invalid - user already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_user").
				SetData(b),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Data)

		})
	}
}

func TestClient_retrieveUser(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	u := getUserStruct()
	cfg := config.Spec{
		Name: "gcp-firebase",
		Kind: "gcp.firebase",
		Properties: map[string]string{
			"project_id":       dat.projectID,
			"credentials":      dat.cred,
			"auth_client":      "true",
			"messaging_client": "false",
			"db_client":        "false",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		{
			name: "retrieve user- by email",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "retrieve_user").
				SetMetadataKeyValue("email", fmt.Sprintf("%s", u["email"])).
				SetMetadataKeyValue("retrieve_by", "by_email"),
			wantErr: false,
		},
		{
			name: "retrieve user-  by_uid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "retrieve_user").
				SetMetadataKeyValue("uid", dat.uid).
				SetMetadataKeyValue("retrieve_by", "by_uid"),
			wantErr: false,
		},
		{
			name: "retrieve user -  by_phone",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "retrieve_user").
				SetMetadataKeyValue("phone", fmt.Sprintf("%s", u["phone_number"])).
				SetMetadataKeyValue("retrieve_by", "by_phone"),
			wantErr: false,
		},
		{
			name: "retrieve user- by email",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "retrieve_user").
				SetMetadataKeyValue("email", fmt.Sprintf("%s", u["email"])).
				SetMetadataKeyValue("retrieve_by", "by_email"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Data)

			t.Logf("received response: %s for test: %s", r.Data, tt.name)
		})
	}
}

func TestClient_listUsers(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "gcp-firebase",
		Kind: "gcp.firebase",
		Properties: map[string]string{
			"project_id":       dat.projectID,
			"credentials":      dat.cred,
			"auth_client":      "true",
			"messaging_client": "false",
			"db_client":        "false",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		{
			name: "list users",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_users"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Data)

			t.Logf("received response: %s for test: %s", r.Data, tt.name)
		})
	}
}

func TestClient_updateUser(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	u := make(map[string]interface{})
	u["email"] = "test315413a@test.com"
	b, err := json.Marshal(&u)
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "gcp-firebase",
		Kind: "gcp.firebase",
		Properties: map[string]string{
			"project_id":       dat.projectID,
			"credentials":      dat.cred,
			"auth_client":      "true",
			"messaging_client": "false",
			"db_client":        "false",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		{
			name: "update user-valid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update_user").
				SetMetadataKeyValue("uid", dat.uid).
				SetData(b),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Data)

		})
	}
}

func TestClient_deleteUser(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	u := []string{dat.uid, dat.uid2}
	b, err := json.Marshal(&u)
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "gcp-firebase",
		Kind: "gcp.firebase",
		Properties: map[string]string{
			"project_id":       dat.projectID,
			"credentials":      dat.cred,
			"auth_client":      "true",
			"messaging_client": "false",
			"db_client":        "false",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		{
			name: "delete user-valid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_user").
				SetMetadataKeyValue("uid", dat.uid),
			wantErr: false,
		},
		{
			name: "delete delete_multiple_users-valid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_multiple_users").
				SetData(b),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Metadata)
		})
	}
}
