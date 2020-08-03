package openfass

import (
	"context"
	"time"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		cfg        config.Spec
		request    *types.Request
		wantStatus string
		wantErr    bool
	}{
		{
			name: "valid request",
			cfg: config.Spec{
				Name: "target.openfass",
				Kind: "target.openfass",
				Properties: map[string]string{
					"gateway":  "http://127.0.0.1:31112",
					"username": "admin",
					"password": "password",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("topic", "function/nslookup").
				SetData([]byte("kubemq.io")),
			wantStatus: "200",
			wantErr:    false,
		},
		{
			name: "invalid request - execution error",
			cfg: config.Spec{
				Name: "target.openfass",
				Kind: "target.openfass",
				Properties: map[string]string{
					"gateway":  "http://bad:31112",
					"username": "admin",
					"password": "password",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("topic", "function/nslookup").
				SetData([]byte("kubemq.io")),
			wantStatus: "200",
			wantErr:    true,
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
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.EqualValues(t, tt.wantStatus, got.Metadata.Get("status"))
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
				Name: "openfass-target",
				Kind: "",
				Properties: map[string]string{
					"gateway":  "http://127.0.0.1:31112",
					"username": "admin",
					"password": "password",
				},
			},
			wantErr: false,
		},
		{
			name: "init - no gateway",
			cfg: config.Spec{
				Name: "openfass-target",
				Kind: "",
				Properties: map[string]string{
					"username": "admin",
					"password": "password",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no username",
			cfg: config.Spec{
				Name: "openfass-target",
				Kind: "",
				Properties: map[string]string{
					"gateway":  "http://127.0.0.1:31112",
					"password": "password",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no password",
			cfg: config.Spec{
				Name: "openfass-target",
				Kind: "",
				Properties: map[string]string{
					"gateway":  "http://127.0.0.1:31112",
					"username": "admin",
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
