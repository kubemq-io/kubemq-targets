package consulkv

import (
	"context"
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
			name: "init",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8500",
				},
			},
			wantErr: false,
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

func TestClient_Put(t *testing.T) {
	kvp := `{"Key":"some-key","CreateIndex":0,"ModifyIndex":0,"LockIndex":0,"Flags":0,"Value":"bXkgdmFsdWU=","Session":""}`
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid put key",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8500",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put").
				SetData([]byte(kvp)),
			wantErr: false,
		}, {
			name: "invalid put key - missing data",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8500",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put"),
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


func TestClient_Get(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid get key",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8500",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-key"),
			wantErr: false,
		}, {
			name: "invalid get key invalid address ",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8511",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-key"),
			wantErr: true,
		}, {
			name: "invalid get key key not exists ",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8511",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "fake-key"),
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
			k, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, k.Data)

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
			name: "valid delete key",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8500",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "some-key"),
			wantErr: false,
		}, {
			name: "invalid delete key - missing key ",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8511",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete"),
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
			k, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, k.Data)

		})
	}
}

func TestClient_List(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list",
			cfg: config.Spec{
				Name: "stores-consulkv",
				Kind: "stores.consulkv",
				Properties: map[string]string{
					"address": "localhost:8500",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			k, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, k.Data)

		})
	}
}
