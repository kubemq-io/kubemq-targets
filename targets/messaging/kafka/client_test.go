package kafka

import (
	"context"
	b64 "encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/stretchr/testify/require"
)

func replaceHeaderValues(req *types.Request) {
	r := strings.NewReplacer(
		"_replaceHK_", b64.StdEncoding.EncodeToString([]byte("header1")),
		"_replaceHV_", b64.StdEncoding.EncodeToString([]byte("headervalue1")))
	req.Metadata["Headers"] = r.Replace(req.Metadata["Headers"])

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
				Name: "kafka-target",
				Kind: "",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topic":   "TestTopic",
				},
			},
			wantErr: false,
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
			name: "valid publish request ",
			cfg: config.Metadata{
				Name: "kafka-target",
				Kind: "",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topic":   "NewTestTopic",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("Key", "S2V5").
				SetData([]byte("new-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("partition", "0").
				SetMetadataKeyValue("offset", "0"),
			wantErr: false,
		},
		{
			name: "valid publish request with headers",
			cfg: config.Metadata{
				Name: "kafka-target",
				Kind: "",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topic":   "NewTestTopic",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("Key", "S2V5").
				SetData([]byte("new-data")).SetMetadataKeyValue(
				"Headers", `[{"Key": "_replaceHK_","Value": "_replaceHV_"}]`),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("partition", "0").
				SetMetadataKeyValue("offset", "1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			replaceHeaderValues(tt.request)
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
