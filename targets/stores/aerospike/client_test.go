package aerospike

import (
	"context"
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go"
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
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "127.0.0.1",
					"port": "3000",
				},
			},
			wantErr: false,
		},
		{
			name: "init - invalid port",
			cfg: config.Spec{
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "127.0.0.1",
					"port": "3001",
				},
			},
			wantErr: true,
		},
		{
			name: "init - invalid host",
			cfg: config.Spec{
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "00.0.0.1",
					"port": "3001",
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
func TestClient_Set_Get(t *testing.T) {
	k := PutRequest{
		UserKey:   "user_key1",
		KeyName:   "some-key",
		Namespace: "test",
		BinMap: // define some bins with data
		aero.BinMap{
			"bin1": 42,
			"bin2": "An elephant is a mouse with an operating system",
			"bin3": []interface{}{"Go", 2009},
		},
	}
	req, err := json.Marshal(k)
	require.NoError(t, err)
	tests := []struct {
		name            string
		cfg             config.Spec
		setRequest      *types.Request
		getRequest      *types.Request
		wantSetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Spec{
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "127.0.0.1",
					"port": "3000",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetData(req),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("namespace", "test").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("user_key", "user_key1"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantSetErr: false,
			wantGetErr: false,
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
			name: "valid Delete",
			cfg: config.Spec{
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "127.0.0.1",
					"port": "3000",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("user_key", "user_key1").
				SetMetadataKeyValue("namespace", "test"),
		}, {
			name: "invalid Delete - no key",
			cfg: config.Spec{
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "127.0.0.1",
					"port": "3000",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("user_key", "user_key1").
				SetMetadataKeyValue("namespace", "test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_GetBatch(t *testing.T) {
	var keys []*string
	key := "some-key"
	keys = append(keys, &key)
	k := GetBatchRequest{
		KeyNames:  keys,
		Namespace: "test",
		BinNames:  nil,
	}
	req, err := json.Marshal(k)
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid GetBatch",
			cfg: config.Spec{
				Name: "stores-aerospike",
				Kind: "stores.aerospike",
				Properties: map[string]string{
					"host": "127.0.0.1",
					"port": "3000",
				},
			},
			request: types.NewRequest().
				SetData(req).
				SetMetadataKeyValue("method", "get_batch").
				SetMetadataKeyValue("user_key", "user_key1").
				SetMetadataKeyValue("namespace", "test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}
