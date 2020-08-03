package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

type payload struct {
	Body string `json:"body"`
}

func newPayload(data string) *payload {
	return &payload{
		Body: data,
	}
}
func (p *payload) Marshal() []byte {
	b, _ := json.Marshal(p)

	return b
}
func unmarshalPayload(data []byte) *payload {
	if len(data) == 0 {
		return nil
	}

	p := &payload{}
	_ = json.Unmarshal(data, p)
	return p
}

type mockHttpServer struct {
	echo      *echo.Echo
	port      string
	postError bool
	getError  bool
}

func (m *mockHttpServer) start() {
	m.echo = echo.New()
	m.echo.POST("/post", func(c echo.Context) error {
		if m.postError {
			return c.String(500, "")
		}
		p := &payload{}
		err := c.Bind(p)
		if err != nil {
			panic(err)
		}
		return c.JSON(200, p)
	})
	m.echo.GET("/get", func(c echo.Context) error {
		if m.getError {
			return c.JSON(500, nil)
		}
		return c.JSON(200, nil)
	})
	go func() {
		_ = m.echo.Start(fmt.Sprintf(":%s", m.port))
	}()
	time.Sleep(time.Second)

}
func (m *mockHttpServer) stop(ctx context.Context) {
	_ = m.echo.Shutdown(ctx)
}
func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		mock    *mockHttpServer
		cfg     config.Spec
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid request - post",
			mock: &mockHttpServer{
				port:      "30000",
				postError: false,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "post").
				SetMetadataKeyValue("url", "http://localhost:30000/post").
				SetData(newPayload("some-data").Marshal()),
			want: types.NewResponse().
				SetData(newPayload("some-data").Marshal()),

			wantErr: false,
		},
		{
			name: "valid request - post error",
			mock: &mockHttpServer{
				port:      "30001",
				postError: true,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "post").
				SetMetadataKeyValue("url", "http://localhost:30001/post").
				SetData(newPayload("some-data").Marshal()),
			want: types.NewResponse(),

			wantErr: false,
		},
		{
			name: "valid request - error on send",
			mock: &mockHttpServer{
				port:      "30002",
				postError: true,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "post").
				SetMetadataKeyValue("url", "http://localhost:40001/post").
				SetData(newPayload("some-data").Marshal()),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid request - method not supported",
			mock: &mockHttpServer{
				port:      "30003",
				postError: true,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "invalid-method").
				SetMetadataKeyValue("url", "http://localhost:30003/post").
				SetData(newPayload("some-data").Marshal()),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid request - no method",
			mock: &mockHttpServer{
				port:      "30004",
				postError: true,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("url", "http://localhost:30004/post").
				SetData(newPayload("some-data").Marshal()),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid request - no url",
			mock: &mockHttpServer{
				port:      "30005",
				postError: true,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "post").
				SetMetadataKeyValue("url", "http://localhost:30005/post").
				SetMetadataKeyValue("headers", `invalid-format`).
				SetData(newPayload("some-data").Marshal()),

			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid request - bad headers format",
			mock: &mockHttpServer{
				port:      "30003",
				postError: true,
				getError:  false,
			},
			cfg: config.Spec{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"auth_type":       "no_auth",
					"username":        "",
					"password":        "",
					"token":           "",
					"default_headers": `{"Content-Type":"application/json"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "post").
				SetData(newPayload("some-data").Marshal()),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tt.mock.start()
			defer tt.mock.stop(ctx)
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)

			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.EqualValues(t, unmarshalPayload(tt.want.Data), unmarshalPayload(got.Data))
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
				Name: "http-target",
				Kind: "",
				Properties: map[string]string{
					"auth_type":          "basic",
					"username":           "username",
					"password":           "password",
					"token":              "token",
					"proxy":              "proxy",
					"retry_count":        "1",
					"retry_wait_seconds": "1",
					"root_certificate":   "some-certificate",
					"client_private_key": "",
					"client_public_key":  "",
					"default_headers":    "",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error on client certificate",
			cfg: config.Spec{
				Name: "http-target",
				Kind: "",
				Properties: map[string]string{
					"auth_type":          "auth_token",
					"username":           "username",
					"password":           "password",
					"token":              "token",
					"proxy":              "proxy",
					"retry_count":        "1",
					"retry_wait_seconds": "1",
					"root_certificate":   "some-certificate",
					"client_private_key": "some-certificate",
					"client_public_key":  "some-certificate",
					"default_headers":    "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - error on bad options 1",
			cfg: config.Spec{
				Name: "http-target",
				Kind: "",
				Properties: map[string]string{
					"retry_wait_seconds": "-1",
				},
			},
			wantErr: true,
		},
		{
			name: "init - error on bad options 2",
			cfg: config.Spec{
				Name: "http-target",
				Kind: "",
				Properties: map[string]string{
					"default_headers": "bad format",
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
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
