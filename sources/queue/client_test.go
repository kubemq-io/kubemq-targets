package queue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/middleware"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/targets/null"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"

	"github.com/kubemq-io/kubemq-targets/targets"
)

func setupClient(ctx context.Context, target middleware.Middleware, channel string) (*Client, error) {
	c := New()

	err := c.Init(ctx, config.Spec{
		Name: "kubemq-queue",
		Kind: "",
		Properties: map[string]string{
			"address":          "localhost:50000",
			"client_id":        "some-clients-id",
			"auth_token":       "",
			"channel":          channel,
			"response_channel": "queue.stream.response",
			"batch_size":       "1",
			"wait_timeout":     "60",
			"sources":          "1",
		},
	}, "", nil)
	if err != nil {
		return nil, err
	}
	err = c.Start(ctx, target)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return c, nil
}

func sendQueueMessage(t *testing.T, ctx context.Context, req *types.Request, sendChannel, respChannel string) (*types.Response, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		return nil, err
	}
	go func() {
		time.Sleep(time.Second)
		result, err := client.SetQueueMessage(
			req.ToQueueMessage()).
			SetChannel(sendChannel).
			SetPolicyMaxReceiveCount(2).
			Send(ctx)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.False(t, result.IsError)
	}()
	if respChannel != "" {
		respMsgs, err := client.ReceiveQueueMessages(ctx,
			client.NewReceiveQueueMessagesRequest().
				SetChannel(respChannel).
				SetClientId(uuid.New().String()).
				SetMaxNumberOfMessages(1).
				SetWaitTimeSeconds(5))
		if err != nil {
			return nil, err
		}
		if len(respMsgs.Messages) == 0 {
			return nil, fmt.Errorf("no messages")
		}
		return types.ParseResponse(respMsgs.Messages[0].Body)
	}
	return nil, nil
}

func TestClient_processQueue(t *testing.T) {
	tests := []struct {
		name        string
		target      targets.Target
		respChannel string
		req         *types.Request
		wantResp    *types.Response
		sendCh      string
		wantErr     bool
	}{
		{
			name: "request",
			target: &null.Client{
				Delay:         time.Second,
				DoError:       nil,
				ResponseError: nil,
			},
			respChannel: "queue.stream.response",
			req:         types.NewRequest().SetData([]byte("some-data")),
			wantResp:    types.NewResponse().SetData([]byte("some-data")),
			sendCh:      uuid.New().String(),
			wantErr:     false,
		},
		{
			name: "request with target do error",
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("do-error"),
				ResponseError: nil,
			},
			respChannel: "queue.stream.response",
			req:         types.NewRequest().SetData([]byte("some-data")),
			wantResp:    types.NewResponse().SetError(fmt.Errorf("do-error")),
			sendCh:      uuid.New().String(),
			wantErr:     false,
		},
		{
			name: "request with target remote error",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: fmt.Errorf("do-error"),
			},
			respChannel: "queue.stream.response",
			req:         types.NewRequest().SetData([]byte("some-data")),
			wantResp:    types.NewResponse().SetError(fmt.Errorf("do-error")),
			sendCh:      uuid.New().String(),
			wantErr:     false,
		},
		{
			name: "bad request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:      nil,
			wantResp: nil,
			sendCh:   uuid.New().String(),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupClient(ctx, tt.target, tt.sendCh)
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()
			gotResp, err := sendQueueMessage(t, ctx, tt.req, tt.sendCh, tt.respChannel)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "kubemq-queue",
				Kind: "",
				Properties: map[string]string{
					"address":          "localhost:50000",
					"client_id":        "some-clients-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"batch_size":       "1",
					"wait_timeout":     "60",
					"sources":          "2",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Spec{
				Name: "kubemq-queue",
				Kind: "",
				Properties: map[string]string{
					"host": "localhost",
					"port": "-1",
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
			if err := c.Init(ctx, tt.cfg, "", nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Start(t *testing.T) {
	tests := []struct {
		name    string
		target  targets.Target
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "start",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			cfg: config.Spec{
				Name: "kubemq-queue",
				Kind: "",
				Properties: map[string]string{
					"address":          "localhost:50000",
					"client_id":        "some-clients-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"batch_size":       "1",
					"wait_timeout":     "60",
					"sources":          "2",
				},
			},
			wantErr: false,
		},
		{
			name:   "start - bad target",
			target: nil,
			cfg: config.Spec{
				Name: "kubemq-queue",
				Kind: "",
				Properties: map[string]string{
					"address":          "localhost:50000",
					"client_id":        "some-clients-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"batch_size":       "1",
					"wait_timeout":     "60",
					"sources":          "2",
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
			_ = c.Init(ctx, tt.cfg, "", nil)

			if err := c.Start(ctx, tt.target); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
