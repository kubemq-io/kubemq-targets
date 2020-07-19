package query

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

	err := c.Init(ctx, config.Metadata{
		Name: "kubemq-rpc",
		Kind: "",
		Properties: map[string]string{
			"host":                       "localhost",
			"port":                       "50000",
			"client_id":                  "",
			"auth_token":                 "some-auth token",
			"channel":                    "query",
			"group":                      "group",
			"concurrency":                "1",
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
func sendQuery(ctx context.Context, req *types.Request, sendChannel string, timeout time.Duration) (*types.Response, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		return nil, err
	}
	queryResponse, err := client.SetQuery(req.ToQuery()).SetChannel(sendChannel).SetTimeout(timeout).Send(ctx)
	if err != nil {
		return nil, err
	}
	return types.ParseResponse(queryResponse.Body)

}
func TestClient_processQuery(t *testing.T) {
	tests := []struct {
		name     string
		target   targets.Target
		req      *types.Request
		wantResp *types.Response
		timeout  time.Duration
		sendCh   string
		wantErr  bool
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
			timeout:  5 * time.Second,
			sendCh:   "query",
			wantErr:  false,
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
			timeout:  5 * time.Second,
			sendCh:   "query",

			wantErr: false,
		},
		{
			name: "request with target do error",
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("error"),
				ResponseError: nil,
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetError(fmt.Errorf("error")),
			timeout:  5 * time.Second,
			sendCh:   "query",

			wantErr: false,
		},
		{
			name: "request with timeout error",
			target: &null.Client{
				Delay:         3 * time.Second,
				DoError:       nil,
				ResponseError: nil,
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse(),
			timeout:  2 * time.Second,
			sendCh:   "query",

			wantErr: true,
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
			gotResp, err := sendQuery(ctx, tt.req, tt.sendCh, tt.timeout)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
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
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Metadata{
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "start",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			cfg: config.Metadata{
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
			cfg: config.Metadata{
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
