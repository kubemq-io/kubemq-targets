package cloudfunctions

import (
	"context"

	"io/ioutil"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

type testStructure struct {
	projectID string
	cred      string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	if err != nil {
		return nil, err
	}
	t.projectID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	if err != nil {
		return nil, err
	}
	t.cred = string(dat)
	return t, nil
}

func TestClient_Init(t *testing.T) {

	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init-missing-credentials",
			cfg: config.Spec{
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init-missing-project-id",
			cfg: config.Spec{
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"credentials": dat.cred,
				},
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

func TestClient_Do(t *testing.T) {

	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid google-function sent",
			cfg: config.Spec{
				Name: "target-google-cloudfunctions",
				Kind: "google.cloudfunctions",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "kube-test").
				SetData([]byte(`{"data":"dGVzdGluZzEyMy4uLg=="}`)),
			want: types.NewResponse().SetData([]byte("test")),

			wantErr: false,
		},
		{
			name: "valid google-function sent",
			cfg: config.Spec{
				Name: "target-google-cloudfunctions",
				Kind: "google.cloudfunctions",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "kube-test").
				SetMetadataKeyValue("project_id", dat.projectID).
				SetMetadataKeyValue("location", "us-central1").SetData([]byte(`{"data":"dGVzdGluZzEyMy4uLg=="}`)),
			want: types.NewResponse().SetData([]byte("test")),

			wantErr: false,
		},
		{
			name: "invalid  google-function location with no match",
			cfg: config.Spec{
				Name: "target-google-cloudfunctions",
				Kind: "google.cloudfunctions",
				Properties: map[string]string{
					"project_id":     dat.projectID,
					"credentials":    dat.cred,
					"location_match": "false",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "test-kubemq").
				SetData([]byte(`{"message":"test"}`)),
			want: types.NewResponse().SetData([]byte("test")),

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
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
