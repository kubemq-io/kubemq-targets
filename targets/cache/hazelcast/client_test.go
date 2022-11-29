package hazelcast

import (
	"context"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
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
				Name: "hazelcast-target",
				Kind: "hazelcast.target",
				Properties: map[string]string{
					"address": "localhost:5701",
				},
			},
			wantErr: false,
		}, {
			name: "init - incorrect address",
			cfg: config.Spec{
				Name: "hazelcast-target",
				Kind: "hazelcast.target",
				Properties: map[string]string{
					"address": "localhost:5702",
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

			err := c.Init(ctx, tt.cfg, nil)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
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
				Name: "hazelcast-target",
				Kind: "hazelcast.target",
				Properties: map[string]string{
					"address": "localhost:5701",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("map_name", "my_map").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("map_name", "my_map").
				SetMetadataKeyValue("key", "some-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok").
				SetData([]byte("some-data")),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "invalid get - key does not exists",
			cfg: config.Spec{
				Name: "hazelcast-target",
				Kind: "hazelcast.target",
				Properties: map[string]string{
					"address": "localhost:5701",
				},
			},
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "fake_key").
				SetMetadataKeyValue("map_name", "my_map"),
			wantSetErr: true,
			wantGetErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			if tt.setRequest != nil {
				gotSetResponse, err := c.Do(ctx, tt.setRequest)
				if tt.wantSetErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, gotSetResponse)
				require.EqualValues(t, tt.wantSetResponse, gotSetResponse)
			}
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
	tests := []struct {
		name         string
		cfg          config.Spec
		request      *types.Request
		getRequest   *types.Request
		wantResponse *types.Response
		wantErr      bool
	}{
		{
			name: "valid delete request",
			cfg: config.Spec{
				Name: "hazelcast-target",
				Kind: "hazelcast.target",
				Properties: map[string]string{
					"address": "localhost:5701",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("map_name", "my_map"),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantResponse, gotSetResponse)
		})
	}
}
