package queue_stream

import (
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptions_parseOptions(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		want    options
		wantErr bool
	}{
		{
			name: "valid options",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "some-clients-id",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"response_channel":           "some-response-channel",
					"visibility_timeout_seconds": "2",
					"wait_timeout":               "60",
				},
			},
			want: options{
				host:              "localhost",
				port:              50000,
				clientId:          "some-clients-id",
				authToken:         "some-auth token",
				channel:           "some-channel",
				responseChannel:   "some-response-channel",
				sources:           1,
				waitTimeout:       60,
				visibilityTimeout: 2,
			},
			wantErr: false,
		},
		{
			name: "invalid options - address",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address": "bad_address",
					"channel": "",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad channel",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
					"channel": "",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad max visibility",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"channel":                    "channel",
					"sources":                    "1",
					"visibility_timeout_seconds": "-1",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad wait timeout",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"channel":                    "channel",
					"sources":                    "1",
					"visibility_timeout_seconds": "1",
					"wait_timeout":               "-1",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad sources",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"channel":                    "channel",
					"sources":                    "-1",
					"visibility_timeout_seconds": "1",
					"wait_timeout":               "1",
				},
			},
			want:    options{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptions(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)

			} else {
				require.NoError(t, err)

			}
			require.EqualValues(t, got, tt.want)
		})
	}
}
