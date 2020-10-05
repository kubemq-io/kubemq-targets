package firebase

import (
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_db(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "google-firebase-target",
		Kind: "target.gcp.firebase",
		Properties: map[string]string{
			"project_id":  dat.projectID,
			"credentials": dat.cred,
			"db_client":   "true",
			"db_url":      dat.dbName,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	m := make(map[string]interface{})
	m["some_key"] = "some_value"
	b, err := json.Marshal(m)
	require.NoError(t, err)
	m2 := make(map[string]interface{})
	m2["some_key"] = "some_value_new"
	b2, err := json.Marshal(m2)
	require.NoError(t, err)
	k := "newValue1"
	bk1, err := json.Marshal(k)
	require.NoError(t, err)
	k2 := "newValue2"
	bk2, err := json.Marshal(k2)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		{
			name: "valid db-get",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_db").
				SetMetadataKeyValue("ref_path", "test"),
			wantErr: false,
		},
		{
			name: "valid db-set - no child ref",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set_db").
				SetMetadataKeyValue("ref_path", "test").
				SetData(bk1),
			wantErr: false,
		}, {
			name: "valid db-set - child ref",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set_db").
				SetMetadataKeyValue("ref_path", "test").
				SetMetadataKeyValue("child_ref", "test_child_ref").
				SetData(bk2),
			wantErr: false,
		},
		{
			name: "valid db-update- no child ref",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update_db").
				SetMetadataKeyValue("ref_path", "test").
				SetData(b),
			wantErr: false,
		},
		{
			name: "valid db-update - child ref ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update_db").
				SetMetadataKeyValue("ref_path", "test").
				SetMetadataKeyValue("child_ref", "test_child_ref").
				SetData(b2),
			wantErr: false,
		},
		{
			name: "valid db-update - child ref ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_db").
				SetMetadataKeyValue("ref_path", "test").
				SetMetadataKeyValue("child_ref", "test_child_ref"),
			wantErr: false,
		}, {
			name: "valid db-update - no child ref ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_db").
				SetMetadataKeyValue("ref_path", "test"),
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

			t.Logf("received response: %s for test: %s", r.Data, tt.name)
		})
	}
}
