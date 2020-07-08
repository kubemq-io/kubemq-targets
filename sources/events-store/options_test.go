package events_store

import (
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptions_parseOptions(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "valid options",
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
					"response_channel":           "",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "valid options - used default host",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "",
					"port":                       "50000",
					"client_id":                  "some-client",
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
			name: "valid options - used default port",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "some-group",
					"response_channel":           "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad port range",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "100000",
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
			wantErr: true,
		},
		{
			name: "invalid options - no channel",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "",
					"group":                      "",
					"concurrency":                "1",
					"response_channel":           "",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid options - bad concurrency value",
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
					"concurrency":                "-1",
					"auto_reconnect":             "true",
					"response_channel":           "",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid options - bad concurrency string",
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
					"concurrency":                "some-bad-int",
					"auto_reconnect":             "true",
					"response_channel":           "",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
		{
			name: "valid options - bad auto reconnect",
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
					"response_channel":           "",
					"auto_reconnect":             "some-bad-error",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "valid options - bad reconnect interval",
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
					"response_channel":           "",
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
