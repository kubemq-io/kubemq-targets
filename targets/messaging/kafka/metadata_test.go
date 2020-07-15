package kafka

import (
	"strings"
	"testing"

	b64 "encoding/base64"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/stretchr/testify/require"
)

func TestMetadata_parseMeta(t *testing.T) {
	tests := []struct {
		name         string
		meta         types.Metadata
		wantMetadata metadata
		wantErr      bool
	}{

		{
			name: "valid simple",
			meta: map[string]string{
				"Key": "_replaceKey_",
			},
			wantMetadata: metadata{
				Key: []byte("key"),
			},
			wantErr: false,
		},
		{
			name: "valid Headers",
			meta: map[string]string{
				"Headers": `[{"Key": "_replaceHK_","Value": "_replaceHV_"}]`,
				"Key":     "_replaceKey_",
			},
			wantMetadata: metadata{
				Headers: []kafka.RecordHeader{
					{
						Key:   []byte("meta1"),
						Value: []byte("dog"),
					},
				},
				Key: []byte("key"),
			},
			wantErr: false,
		},
		{
			name: "invalid Headers_value is not base64",
			meta: map[string]string{
				"Headers": `[{"Key": "meta1","Value": "badvalue"}]`,
				"Key":     "_replaceKey_",
			},
			wantMetadata: metadata{
				Headers: []kafka.RecordHeader{
					{
						Key:   []byte("meta1"),
						Value: []byte("badvalue"),
					},
				},
				Key: []byte("key"),
			},
			wantErr: true,
		}, {
			name: "invalid Headers_json bad format ",
			meta: map[string]string{
				"Headers": `[{"Key": "meta1""Value": "badvalue"}]`,
				"Key":     "_replaceKey_",
			},
			wantMetadata: metadata{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReplacer(
				"_replaceHV_", b64.StdEncoding.EncodeToString([]byte("meta1")),
				"_replaceHK_", b64.StdEncoding.EncodeToString([]byte("dog")),
				"_replaceKey_", b64.StdEncoding.EncodeToString([]byte("key")))
			tt.meta["Headers"] = r.Replace(tt.meta["Headers"])
			tt.meta["Key"] = r.Replace(tt.meta["Key"])
			meta, err := parseMetadata(tt.meta, options{})
			if tt.wantErr {
				require.Error(t, err)

			} else {
				require.NoError(t, err)

			}
			require.EqualValues(t, tt.wantMetadata, meta)
		})
	}
}
