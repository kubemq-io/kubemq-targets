package amqp

import (
	"context"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"
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
				Name: "messaging-amqp",
				Kind: "messaging.amqp",
				Properties: map[string]string{
					"url":      "amqp://192.168.50.95:5672/",
					"username": "artemis",
					"password": "simetraehcapa",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad url",
			cfg: config.Spec{
				Name: "messaging-amqp",
				Kind: "messaging-amqp",
				Properties: map[string]string{
					"url": "amqp://amqp:amqp@localhost:6000/",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no auth",
			cfg: config.Spec{
				Name: "messaging-amqp",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://localhost:5672/",
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
			name: "valid publish request without confirmation",
			cfg: config.Spec{
				Name: "messaging-amqp",
				Kind: "messaging.amqp",
				Properties: map[string]string{
					"url":      "amqp://localhost:5672/",
					"username": "artemis",
					"password": "simetraehcapa",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("address", "some-address").
				SetMetadataKeyValue("to", "q2").
				SetMetadataKeyValue("durable", "true").
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
			time.Sleep(2 * time.Second)
		})
	}
}
