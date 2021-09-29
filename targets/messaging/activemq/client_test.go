package activemq

import (
	"context"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
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
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "admin",
					"password": "admin",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - bad host",
			cfg: config.Spec{
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"host":     "localhost:6000",
					"username": "",
					"password": "",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - no host",
			cfg: config.Spec{
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"username": "",
					"password": "",
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

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantExecErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.Spec
		request      *types.Request
		wantResponse *types.Response
		wantErr      bool
	}{
		{
			name: "valid publish request",
			cfg: config.Spec{
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "",
					"password": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("destination", "some-queue").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
		{
			name: "invalid publish request - no destination",
			cfg: config.Spec{
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "admin",
					"password": "admin",
				},
			},
			request: types.NewRequest().
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: true,
		},
		{
			name: "valid publish request with default destination",
			cfg: config.Spec{
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"host":                "localhost:61613",
					"username":            "",
					"password":            "",
					"default_destination": "q1",
				},
			},
			request: types.NewRequest().
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
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
			gotResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotResponse)
			require.EqualValues(t, tt.wantResponse, gotResponse)
		})
	}
}
