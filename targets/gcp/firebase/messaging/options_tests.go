package messaging

import (
	"encoding/json"
	"fmt"
	"testing"

	"firebase.google.com/go/messaging"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

func TestParseMetaData(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		ops     options
		wantErr bool
		Request *types.Request
	}{
		{
			name: "valid method write",
			cfg: config.Spec{
				Name: "google-big-table-target",
				Kind: "",
			},
			ops: options{defult: &messaging.Message{
				Android: &messaging.AndroidConfig{
					Notification: &messaging.AndroidNotification{
						Body: "hello android",
					},
				},
			}},
			Request: types.NewRequest().
				SetMetadataKeyValue("method", "write"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := json.Marshal(tt.ops.defult)
			if err != nil {
				t.Logf("init() error = %v", err)
			}
			fmt.Print(p)
			require.NoError(t, err)
		})
	}
}
