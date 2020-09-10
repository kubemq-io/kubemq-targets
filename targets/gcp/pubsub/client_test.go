package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type testStructure struct {
	projectID   string
	credentials string
	topicID     string
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
	t.credentials = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../credentials/topicID.txt")
	if err != nil {
		return nil, err
	}
	t.topicID = fmt.Sprintf("%s", dat)
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
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"retries":     "0",
					"credentials": dat.credentials,
				},
			},
			wantErr: false,
		}, {
			name: "init-missing-credentials",
			cfg: config.Spec{
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"retries":    "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Spec{
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"retries":     "0",
					"credentials": dat.credentials,
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

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	validBody, _ := json.Marshal("valid body")
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid google-pubsub sent",
			cfg: config.Spec{
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"retries":     "0",
					"topic_id":    dat.topicID,
					"credentials": dat.credentials,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("topic_id", dat.topicID).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		}, {
			name: "missing topic google-pubsub sent",
			cfg: config.Spec{
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"retries":     "0",
					"credentials": dat.credentials,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

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

func TestClient_list(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "valid google-pubsub-list",
			cfg: config.Spec{
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.credentials,
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
