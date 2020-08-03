package events

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/targets/null"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-targets/targets"
)

func setupClient(ctx context.Context, target middleware.Middleware) (*Client, error) {
	c := New()
	err := c.Init(ctx, config.Spec{
		Name: "kubemq-rpc",
		Kind: "",
		Properties: map[string]string{
			"host":                       "localhost",
			"port":                       "50000",
			"client_id":                  "",
			"auth_token":                 "",
			"channel":                    "events",
			"group":                      "some-group",
			"concurrency":                "1",
			"response_channel":           "events.response",
			"auto_reconnect":             "true",
			"reconnect_interval_seconds": "1",
			"max_reconnects":             "0",
		},
	})
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
func sendEvent(t *testing.T, ctx context.Context, req *types.Request, sendChannel, respChannel string) (*types.Response, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	if err != nil {
		return nil, err
	}
	go func() {
		time.Sleep(time.Second)
		err = client.SetEvent(req.ToEvent()).SetChannel(sendChannel).Send(ctx)
		require.NoError(t, err)
	}()
	if respChannel != "" {
		errCh := make(chan error, 1)
		eventCh, err := client.SubscribeToEvents(ctx, respChannel, "", errCh)
		if err != nil {
			return nil, err
		}
		select {
		case event := <-eventCh:
			return types.ParseResponse(event.Body)
		case err := <-errCh:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, nil
}
func TestClient_processEvent(t *testing.T) {
	tests := []struct {
		name     string
		target   targets.Target
		req      *types.Request
		wantResp *types.Response
		sendCh   string
		respCh   string

		wantErr bool
	}{
		{
			name: "request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetData([]byte("some-data")),
			wantErr:  false,
			sendCh:   "events",
			respCh:   "events.response",
		},
		{
			name: "request with target error",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: fmt.Errorf("error"),
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetMetadataKeyValue("error", "error"),
			wantErr:  false,
			sendCh:   "events",
			respCh:   "events.response",
		},
		{
			name: "request with target error 2",
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("error"),
				ResponseError: nil,
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetError(fmt.Errorf("error")),
			wantErr:  false,
			sendCh:   "events",
			respCh:   "events.response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupClient(ctx, tt.target)
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()

			gotResp, err := sendEvent(t, ctx, tt.req, tt.sendCh, tt.respCh)
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
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"response_channel":           "",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Spec{
				Name: "kubemq-rpc",
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
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
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "false",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name:   "start - bad target",
			target: nil,
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
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
			_ = c.Init(ctx, tt.cfg)

			if err := c.Start(ctx, tt.target); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
