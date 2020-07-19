package mongodb

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "majority",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error connection",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "bad-host:32017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "majority",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad host",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad database",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad collection",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad write concurrency",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "bad-concurrency",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad read concurrency",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "bad-concurrency",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad operation timeout",
			cfg: config.Metadata{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "-2",
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
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
func TestClient_Set_Get(t *testing.T) {
	tests := []struct {
		name            string
		cfg             config.Metadata
		setRequest      *types.Request
		getRequest      *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Metadata{
				Name: "target.mongodb",
				Kind: "target.mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "valid set , no key get request",
			cfg: config.Metadata{
				Name: "target.mongodb",
				Kind: "target.mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "bad-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: nil,
			wantSetErr:      false,
			wantGetErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.setRequest)
			if tt.wantSetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantSetResponse, gotSetResponse)
			gotGetResponse, err := c.Do(ctx, tt.getRequest)
			if tt.wantGetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.EqualValues(t, tt.wantGetResponse, gotGetResponse)
		})
	}
}
func TestClient_Delete(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := New()
	err := c.Init(ctx, config.Metadata{
		Name: "target.mongodb",
		Kind: "target.mongodb",
		Properties: map[string]string{
			"host":                      "localhost:27017",
			"username":                  "admin",
			"password":                  "password",
			"database":                  "admin",
			"collection":                "test",
			"write_concurrency":         "",
			"read_concurrency":          "",
			"params":                    "",
			"operation_timeout_seconds": "2",
		},
	})
	key := nuid.Next()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("key", key).
		SetData([]byte("some-data"))

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("key", key)
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, []byte("some-data"), gotGetResponse.Data)

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("key", key)
	_, err = c.Do(ctx, delRequest)
	require.NoError(t, err)
	gotGetResponse, err = c.Do(ctx, getRequest)
	require.Error(t, err)
	require.Nil(t, gotGetResponse)
}
func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Metadata{
				Name: "target.mongodb",
				Kind: "target.mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Metadata{
				Name: "target.mongodb",
				Kind: "target.mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - no key",
			cfg: config.Metadata{
				Name: "target.mongodb",
				Kind: "target.mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetData([]byte("some-data")),
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
			_, err = c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
