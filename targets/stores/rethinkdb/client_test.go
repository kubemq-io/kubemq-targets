package rethinkdb

import (
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init - valid",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - invalid host",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28016",
					"username": "root",
					"password": "root",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - invalid user",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "fake",
					"password": "root",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - invalid password",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "fake",
					"password": "root",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
func TestClient_Insert(t *testing.T) {
	args := make(map[string]interface{})
	args["id"] = "test_user"
	args["password"] = 1
	insert, err := json.Marshal(args)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid insert request",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("table", "test").
				SetData(insert),
		},
		{
			name: "invalid insert request - missing db name",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetMetadataKeyValue("table", "test").
				SetData(insert),
			wantErr: true,
		}, {
			name: "invalid insert request - missing table name",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetMetadataKeyValue("db_name", "test").
				SetData(insert),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			response, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
		})
	}
}

func TestClient_Update(t *testing.T) {
	args := make(map[string]interface{})
	args["id"] = "test_user"
	args["password"] = 2
	update, err := json.Marshal(args)
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid update request",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test").
				SetData(update),
		}, {
			name: "invalid update request missing data",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test"),
			wantErr: true,
		}, {
			name: "invalid update request missing table",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user").
				SetData(update),
			wantErr: true,
		}, {
			name: "invalid update request missing db",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test").
				SetData(update),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			response, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
		})
	}
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid get request",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test"),
		}, {
			name: "invalid get request - missing db_name",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test"),
			wantErr: true,
		}, {
			name: "invalid get request - missing table",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user"),
			wantErr: true,
		}, {
			name: "invalid get request - missing key does not exists",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "fake_key").
				SetMetadataKeyValue("table", "test"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			response, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
		})
	}
}

func TestClient_Delete(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete request",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test"),
		}, {
			name: "invalid delete request - missing db",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "test_user").
				SetMetadataKeyValue("table", "test"),
			wantErr: true,
		}, {
			name: "invalid delete request - missing table",
			cfg: config.Spec{
				Name: "stores-rethinkdb",
				Kind: "stores.rethinkdb",
				Properties: map[string]string{
					"host":     "localhost:28015",
					"username": "root",
					"password": "root",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("db_name", "test").
				SetMetadataKeyValue("key", "test_user"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			response, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
		})
	}
}
