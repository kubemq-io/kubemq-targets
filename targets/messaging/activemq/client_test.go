package activemq

import (
	"context"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "activemq-target",
				Kind: "",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "",
					"password": "",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad host",
			cfg: config.Metadata{
				Name: "activemq-target",
				Kind: "",
				Properties: map[string]string{
					"host":     "localhost:6000",
					"username": "",
					"password": "",
				},
			},
			wantErr: true,
		}, {
			name: "init - no host",
			cfg: config.Metadata{
				Name: "activemq-target",
				Kind: "",
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

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantExecErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.Metadata
		request      *types.Request
		wantResponse *types.Response
		wantErr      bool
	}{
		{
			name: "valid publish request",
			cfg: config.Metadata{
				Name: "activemq-target",
				Kind: "",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "",
					"password": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("destination", "some-destination").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
		{
			name: "invalid publish request - no destination",
			cfg: config.Metadata{
				Name: "activemq-target",
				Kind: "",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "",
					"password": "",
				},
			},
			request: types.NewRequest().
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
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
