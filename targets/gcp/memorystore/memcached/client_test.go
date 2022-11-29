package memcached

import (
	"context"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "gcp-memcached",
				Kind: "gcp.cache.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:31211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - error no connection",
			cfg: config.Spec{
				Name: "gcp-memcached",
				Kind: "gcp.cache.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:3000",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - bad options - invalid hosts",
			cfg: config.Spec{
				Name: "gcp-memcached",
				Kind: "gcp.cache.memcached",
				Properties: map[string]string{
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - bad options - invalid max idle connection",
			cfg: config.Spec{
				Name: "gcp-memcached",
				Kind: "gcp.cache.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "-1",
					"default_timeout_seconds": "10",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - bad options - invalid default timeout seconds",
			cfg: config.Spec{
				Name: "gcp-memcached",
				Kind: "gcp.cache.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "-1",
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

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_Set_Get(t *testing.T) {
	tests := []struct {
		name            string
		cfg             config.Spec
		setRequest      *types.Request
		getRequest      *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Spec{
				Name: "google.target.memcached",
				Kind: "google.target.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
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
			cfg: config.Spec{
				Name: "google.target.memcached",
				Kind: "google.target.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
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
			err := c.Init(ctx, tt.cfg, nil)
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
	err := c.Init(ctx, config.Spec{
		Name: "google.target.memcached",
		Kind: "google.target.memcached",
		Properties: map[string]string{
			"hosts":                   "localhost:11211",
			"max_idle_connections":    "2",
			"default_timeout_seconds": "10",
		},
	}, nil)
	key := uuid.New().String()
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
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Spec{
				Name: "google.target.memcached",
				Kind: "google.target.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
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
			cfg: config.Spec{
				Name: "google.target.memcached",
				Kind: "google.target.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
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
			cfg: config.Spec{
				Name: "google.target.memcached",
				Kind: "google.target.memcached",
				Properties: map[string]string{
					"hosts":                   "localhost:11211",
					"max_idle_connections":    "2",
					"default_timeout_seconds": "10",
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
			err := c.Init(ctx, tt.cfg, nil)
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
