package cloudfunctions

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

func TestClient_Init(t *testing.T) {

	dat, err := ioutil.ReadFile("./../../../credentials/google_cred.json")
	require.NoError(t, err)
	credentials := fmt.Sprintf("%s", dat)
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
					"project":     "pubsubdemo-281010",
					"credentials": credentials,
				},
			},
			wantErr: false,
		}, {
			name: "init-missing-credentials",
			cfg: config.Metadata{
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"project": "pubsubdemo-281010",
				},
			},
			wantErr: true,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Metadata{
				Name: "google-pubsub-target",
				Kind: "",
				Properties: map[string]string{
					"credentials": credentials,
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

	c, err := ioutil.ReadFile("./../../../credentials/pubsubdemo-281010-11464d7de470.json")
	require.NoError(t, err)
	credentials := fmt.Sprintf("%s", c)
	tests := []struct {
		name    string
		cfg     config.Metadata
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid google-function sent",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project":     "pubsubdemo-281010",
					"credentials": credentials,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "kube-test-pubsub").
				SetData([]byte(`{"data":"dGVzdGluZzEyMy4uLg=="}`)),
			want: types.NewResponse().SetData([]byte("test")),

			wantErr: false,
		},

		{
			name: "valid google-function sent",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project":     "pubsubdemo-281010",
					"credentials": credentials,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "test-kubemq").
				SetData([]byte(`{"message":"test"}`)),
			want: types.NewResponse().SetData([]byte("test")),

			wantErr: false,
		},
		{
			name: "valid google-pubsub sent",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project":     "pubsubdemo-281010",
					"credentials": credentials,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "test-kubemq").
				SetMetadataKeyValue("project", "pubsubdemo-281010").
				SetMetadataKeyValue("location", "us-central1").SetData([]byte(`{"message":"test"}`)),
			want: types.NewResponse().SetData([]byte("test")),

			wantErr: false,
		},
		{
			name: "missing  google-function location with no match",
			cfg: config.Metadata{
				Name: "target.google.pubsub",
				Kind: "target.google.pubsub",
				Properties: map[string]string{
					"project":       "pubsubdemo-281010",
					"credentials":   credentials,
					"locationMatch": "false",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("name", "test-kubemq").
				SetData([]byte(`{"message":"test"}`)),
			want: types.NewResponse().SetData([]byte("test")),

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
