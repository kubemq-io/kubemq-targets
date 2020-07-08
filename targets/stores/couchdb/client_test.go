package couchdb

import (
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type doc struct {
	Data string `json:"data"`
}

func newDoc(data string) *doc {
	return &doc{Data: data}
}
func (d *doc) binary() []byte {
	b, _ := json.Marshal(d)
	return b
}
func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "some-collection",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad url",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"url":              "couchdb://localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no url",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad username and password",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "bad-couchdb",
					"password":         "bad-couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no bucket",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad num to replicate",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "-1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad num to persist",
			cfg: config.Metadata{
				Name: "couchdb-target",
				Kind: "",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "-1",
					"collection":       "",
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
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key-2").
				SetData(newDoc("some-data").binary()),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-key-2"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key-2").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key-2").
				SetMetadataKeyValue("error", "false").
				SetData(newDoc("some-data").binary()),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "valid set , no key get request",
			cfg: config.Metadata{
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key-2").
				SetData(newDoc("some-data").binary()),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "bad-key"),
			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key-2").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: nil,
			wantSetErr:      false,
			wantGetErr:      true,
		},
		{
			name: "invalid set request",
			cfg: config.Metadata{
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "100",
					"num_to_persist":   "100",
					"collection":       "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key-2"),

			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-key-2"),
			wantSetResponse: nil,
			wantGetResponse: nil,
			wantSetErr:      true,
			wantGetErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		Name: "target.couchdb",
		Kind: "target.couchdb",
		Properties: map[string]string{
			"url":              "localhost",
			"username":         "couchdb",
			"password":         "couchdb",
			"bucket":           "bucket",
			"num_to_replicate": "1",
			"num_to_persist":   "1",
			"collection":       "",
		},
	})
	key := nuid.Next()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("key", key).
		SetData(newDoc("some-data").binary())

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("key", key)
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, newDoc("some-data").binary(), gotGetResponse.Data)

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
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData(newDoc("some-data").binary()),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Metadata{
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
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
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - bad cas",
			cfg: config.Metadata{
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("cas", "-1").
				SetData(newDoc("some-data").binary()),
			wantErr: true,
		},
		{
			name: "invalid request - bad expiry",
			cfg: config.Metadata{
				Name: "target.couchdb",
				Kind: "target.couchdb",
				Properties: map[string]string{
					"url":              "localhost",
					"username":         "couchdb",
					"password":         "couchdb",
					"bucket":           "bucket",
					"num_to_replicate": "1",
					"num_to_persist":   "1",
					"collection":       "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("expiry_seconds", "-1").
				SetData(newDoc("some-data").binary()),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
