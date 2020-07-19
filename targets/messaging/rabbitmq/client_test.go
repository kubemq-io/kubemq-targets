package rabbitmq

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
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
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad url",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:6000/",
				},
			},
			wantErr: true,
		}, {
			name: "init - no url",
			cfg: config.Metadata{
				Name:       "rabbitmq-target",
				Kind:       "",
				Properties: map[string]string{},
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
			name: "valid publish request with confirmation",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("queue", "some-queue").
				SetMetadataKeyValue("exchange", "").
				SetMetadataKeyValue("confirm", "true").
				SetMetadataKeyValue("delivery_mode", "2").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("ack", "true").
				SetMetadataKeyValue("delivery_tag", "1"),
			wantErr: false,
		},
		{
			name: "valid publish request without confirmation",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("queue", "some-queue").
				SetMetadataKeyValue("exchange", "").
				SetMetadataKeyValue("confirm", "false").
				SetMetadataKeyValue("delivery_mode", "1").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
		{
			name: "invalid publish request - no queue",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("exchange", "").
				SetData([]byte("some-data")),
			wantResponse: nil,
			wantErr:      true,
		},
		{
			name: "invalid publish request - bad delivery mode",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("queue", "some-queue").
				SetMetadataKeyValue("exchange", "").
				SetMetadataKeyValue("confirm", "false").
				SetMetadataKeyValue("delivery_mode", "3"),
			wantResponse: nil,
			wantErr:      true,
		},
		{
			name: "invalid publish request - bad priority",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("queue", "some-queue").
				SetMetadataKeyValue("exchange", "").
				SetMetadataKeyValue("confirm", "false").
				SetMetadataKeyValue("priority", "12"),
			wantResponse: nil,
			wantErr:      true,
		},
		{
			name: "invalid publish request - bad expiry",
			cfg: config.Metadata{
				Name: "rabbitmq-target",
				Kind: "",
				Properties: map[string]string{
					"url": "amqp://rabbitmq:rabbitmq@localhost:5672/",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("queue", "some-queue").
				SetMetadataKeyValue("exchange", "").
				SetMetadataKeyValue("confirm", "false").
				SetMetadataKeyValue("expiry_seconds", "-1"),
			wantResponse: nil,
			wantErr:      true,
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
