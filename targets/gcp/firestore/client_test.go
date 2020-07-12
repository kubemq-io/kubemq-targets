package firestore

import (
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "google-firestore-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": projectID,
				},
			},
			wantErr: false,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Metadata{
				Name:       "google-firestore-target",
				Kind:       "",
				Properties: map[string]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
			defer cancel()
			c := New()

			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
func TestClient_Set_Get(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	require.NoError(t, err)
	user := map[string]interface{}{
		"first": "kubemq",
		"last":  "kubemq-last",
		"id":    123,
	}
	bUser, err := json.Marshal(user)
	require.NoError(t, err)
	dat, err = ioutil.ReadFile("./../../../credentials/objKey.txt")
	require.NoError(t, err)
	objKey := string(dat)
	tests := []struct {
		name            string
		cfg             config.Metadata
		setRequest      *types.Request
		getRequest      *types.Request
		getAllRequest   *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Metadata{
				Name: "google-firestore-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": projectID,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "add").
				SetMetadataKeyValue("collection", "myCollection").
				SetData(bUser),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "document_key").
				SetMetadataKeyValue("item", objKey).
				SetMetadataKeyValue("error", "false").
				SetMetadataKeyValue("collection", "myCollection"),
			getAllRequest: types.NewRequest().
				SetMetadataKeyValue("method", "documents_all").
				SetMetadataKeyValue("error", "false").
				SetMetadataKeyValue("collection", "myCollection"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetMetadataKeyValue("error", "false").
				SetMetadataKeyValue("collection", "myCollection"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("collection", "myCollection").
				SetMetadataKeyValue("item", objKey).
				SetMetadataKeyValue("error", "false").
				SetData(bUser),
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
			require.EqualValues(t, tt.wantGetResponse, gotGetResponse)
			gotGetAllResponse, err := c.Do(ctx, tt.getAllRequest)
			if tt.wantGetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.Equal(t, gotGetAllResponse.Metadata["error"], "false")
		})
	}
}

func TestClient_Delete(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	require.NoError(t, err)
	dat, err = ioutil.ReadFile("./../../../credentials/deleteKey.txt")
	require.NoError(t, err)
	deleteKey := string(dat)
	tests := []struct {
		name              string
		cfg               config.Metadata
		deleteRequest     *types.Request
		wantDeleteRequest *types.Response
		wantErr           bool
	}{
		{
			name: "valid delete request",
			cfg: config.Metadata{
				Name: "google-firestore-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": projectID,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_document_key").
				SetMetadataKeyValue("item", deleteKey).
				SetMetadataKeyValue("collection", "myCollection"),
			wantDeleteRequest: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetMetadataKeyValue("error", "false").
				SetMetadataKeyValue("item", deleteKey).
				SetMetadataKeyValue("collection", "myCollection"),
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
			deleteResponse, err := c.Do(ctx, tt.deleteRequest)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, deleteResponse)
			require.EqualValues(t, tt.wantDeleteRequest, deleteResponse)
		})
	}
}

func TestClient_list(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "valid google-firestore-list",
			cfg: config.Metadata{
				Name: "target.google.firestore",
				Kind: "target.google.firestore",
				Properties: map[string]string{
					"project_id": projectID,
				},
			},

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
			got, err := c.list(ctx)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, "ok", got.Metadata["result"])
			require.NotNil(t, got)
		})
	}
}
