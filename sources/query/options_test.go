package query

import (
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptions_parseOptions(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "valid options",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid options - no channel",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "",
					"group":                      "",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
		{
			name: "valid options - bad auto reconnect",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"auto_reconnect":             "some-bad-error",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "valid options - bad reconnect interval",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "bad value",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseOptions(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}
