package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
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
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	require.NoError(t, err)
	credentials := fmt.Sprintf("%s", dat)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "gcp-firestore",
				Kind: "gcp.firestore",
				Properties: map[string]string{
					"project_id":  projectID,
					"credentials": credentials,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init-missing-credentials",
			cfg: config.Spec{
				Name: "gcp-firestore",
				Kind: "gcp.firestore",
				Properties: map[string]string{
					"project_id": projectID,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init-missing-project-id",
			cfg: config.Spec{
				Name:       "gcp-firestore",
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

			err := c.Init(ctx, tt.cfg, nil)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)

		})
	}
}
func TestClient_Set_Get(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	require.NoError(t, err)
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	require.NoError(t, err)
	credentials := fmt.Sprintf("%s", dat)
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
		cfg             config.Spec
		setRequest      *types.Request
		getRequest      *types.Request
		getAllRequest   *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
	}{
		{
			name: "valid set get request",
			cfg: config.Spec{
				Name: "gcp-firestore",
				Kind: "gcp.firestore",
				Properties: map[string]string{
					"project_id":  projectID,
					"credentials": credentials,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "add").
				SetMetadataKeyValue("collection", "myCollection").
				SetData(bUser),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "document_key").
				SetMetadataKeyValue("item", objKey).
				SetMetadataKeyValue("collection", "myCollection"),
			getAllRequest: types.NewRequest().
				SetMetadataKeyValue("method", "documents_all").
				SetMetadataKeyValue("collection", "myCollection"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetMetadataKeyValue("collection", "myCollection"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("collection", "myCollection").
				SetMetadataKeyValue("item", objKey).
				SetData(bUser),
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
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantSetResponse, gotSetResponse)
			gotGetResponse, err := c.Do(ctx, tt.getRequest)
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.EqualValues(t, tt.wantGetResponse, gotGetResponse)
			gotGetAllResponse, err := c.Do(ctx, tt.getAllRequest)
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.Equal(t, gotGetAllResponse.Error, "")
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
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	require.NoError(t, err)
	credentials := fmt.Sprintf("%s", dat)
	tests := []struct {
		name              string
		cfg               config.Spec
		deleteRequest     *types.Request
		wantDeleteRequest *types.Response
		wantErr           bool
	}{
		{
			name: "valid delete request",
			cfg: config.Spec{
				Name: "gcp-firestore",
				Kind: "gcp.firestore",
				Properties: map[string]string{
					"project_id":  projectID,
					"credentials": credentials,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_document_key").
				SetMetadataKeyValue("item", deleteKey).
				SetMetadataKeyValue("collection", "myCollection"),
			wantDeleteRequest: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetMetadataKeyValue("item", deleteKey).
				SetMetadataKeyValue("collection", "myCollection"),
			wantErr: false,
		}, {
			name: "invalid delete request - missing item",
			cfg: config.Spec{
				Name: "gcp-firestore",
				Kind: "gcp.firestore",
				Properties: map[string]string{
					"project_id":  projectID,
					"credentials": credentials,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_document_key").
				SetMetadataKeyValue("collection", "myCollection"),
			wantDeleteRequest: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetMetadataKeyValue("error", "false").
				SetMetadataKeyValue("item", "fake-key").
				SetMetadataKeyValue("collection", "myCollection"),
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
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	require.NoError(t, err)
	credentials := fmt.Sprintf("%s", dat)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "valid google-firestore-list",
			cfg: config.Spec{
				Name: "google.firestore",
				Kind: "google.firestore",
				Properties: map[string]string{
					"project_id":  projectID,
					"credentials": credentials,
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
			err := c.Init(ctx, tt.cfg, nil)
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
