package google

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
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"project_id":          projectID,
					"retries":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Metadata{
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"retries": "0",
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
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/topicID.txt")
	require.NoError(t, err)
	TopicID:= string(dat)
	validBody, _ := json.Marshal("valid body")
	tests := []struct {
		name    string
		cfg     config.Metadata
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid google-pubsub sent",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project_id":          projectID,
					"retries":             "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("topic_id", TopicID).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		}, {
			name: "missing topic google-pubsub sent",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project_id":          projectID,
					"retries":             "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("topic_id", `{"tag-1":"test","tag-2":"test2"}`).
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
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "valid google-pubsub-list",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project_id":          projectID,
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
			require.EqualValues(t,"ok",got.Metadata["result"])
			require.NotNil(t, got)
		})
	}
}
