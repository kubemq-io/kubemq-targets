package queue

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
					"host":             "localhost",
					"port":             "50000",
					"client_id":        "some-client-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"concurrency":      "1",
					"batch_size":       "1",
					"wait_timeout":     "60",
				},
			},
			want: options{
				host:            "localhost",
				port:            50000,
				clientId:        "some-client-id",
				authToken:       "some-auth token",
				channel:         "some-channel",
				responseChannel: "some-response-channel",
				concurrency:     1,
				batchSize:       1,
				waitTimeout:     60,
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad port",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host": "localhost",
					"port": "-1",
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
					"host":    "localhost",
					"port":    "50000",
					"channel": "",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad concurrency",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":        "localhost",
					"port":        "50000",
					"channel":     "channel",
					"concurrency": "-1",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad batch size",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":        "localhost",
					"port":        "50000",
					"channel":     "channel",
					"concurrency": "1",
					"batch_size":  "-1",
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
					"host":         "localhost",
					"port":         "50000",
					"channel":      "channel",
					"concurrency":  "1",
					"batch_size":   "1",
					"wait_timeout": "-1",
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
