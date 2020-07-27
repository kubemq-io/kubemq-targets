package firebase

import (
	"encoding/json"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
)

func TestMessageMetadata(t *testing.T) {

	m := &messaging.Message{
		Topic: "test",
		Token: "1231231",
	}
	mb, err := json.Marshal(m)
	if err != nil {
		return
	}
	multi := &messaging.MulticastMessage{
		Tokens: []string{"123", "456"},
		Notification: &messaging.Notification{
			Title: "title",
		}}

	multib, err := json.Marshal(multi)
	if err != nil {
		return
	}

	tests := []struct {
		name    string
		isMulti bool
		request *types.Request
		wantmsg messages
		wantErr bool
	}{
		{
			name:    "parse message",
			isMulti: false,
			request: types.NewRequest().SetMetadataKeyValue("message", string(mb)),
			wantmsg: messages{single: m},
			wantErr: true,
		},
		{
			name:    "parse multicast msg",
			isMulti: true,
			request: types.NewRequest().SetMetadataKeyValue("multicast", string(multib)),
			wantmsg: messages{multicast: &messaging.MulticastMessage{
				Tokens: []string{"123", "456"},
				Notification: &messaging.Notification{
					Title: "title",
				},
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isMulti {
				m, err := parseMetadataMessages(tt.request.Metadata, options{
					defaultMessaging: &messages{},
				}, SendBatch)
				require.NoError(t, err)
				require.EqualValues(t, tt.wantmsg.multicast, m.multicast)

			} else {
				m, err := parseMetadataMessages(tt.request.Metadata, options{
					defaultMessaging: &messages{},
				}, SendMessage)
				require.NoError(t, err)
				require.EqualValues(t, tt.wantmsg.single, m.single)
			}
		})
	}
}

func TestOptionsParse(t *testing.T) {
	m := messages{
		single: &messaging.Message{
			Topic: "newmsg",
		},
	}
	mb, err := json.Marshal(m.single)
	if err != nil {
		return
	}
	ms := string(mb)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantmsg *messages
		wantErr bool
	}{
		{
			name: "parse options",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messagingClient": "true",
					"defaultmsg":      ms,
				},
			},
			wantmsg: &m,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := parseOptions(tt.cfg)
			require.NoError(t, err)
			require.EqualValues(t, tt.wantmsg, o.defaultMessaging)
		})
	}
}
func TestDefultMessge(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantmsg *messages
		wantErr bool
	}{
		{
			name: "parse options",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messagingClient": "true",
					"defaultmsg":      `{"topic":"defult"}`,
				},
			},
			wantmsg: &messages{single: &messaging.Message{
				Topic: "defult",
			},
			},
			request: &types.Request{},
			wantErr: false,
		},
		{
			name: "defult single add request fields",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messagingClient": "true",
					"defaultmsg":      `{"topic":"defult","token":"1234"}`,
				},
			},
			request: types.NewRequest().SetMetadataKeyValue("message", `{"topic":"newTopic"}`),
			wantmsg: &messages{single: &messaging.Message{
				Topic: "newTopic",
				Token: "1234",
			},
			},
			wantErr: false,
		},
		{
			name: "parse defult and msg with ",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messagingClient": "true",
					"defaultmsg":      `{"topic":"defult"}`,
				},
			},
			request: types.NewRequest().SetMetadataKeyValue("message", `{"topic":"newTopic"}`),
			wantmsg: &messages{single: &messaging.Message{
				Topic: "newTopic",
			},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := parseOptions(tt.cfg)
			require.NoError(t, err)
			m, err := parseMetadataMessages(tt.request.Metadata, o, SendMessage)
			require.NoError(t, err)
			require.EqualValues(t, tt.wantmsg, m)
		})
	}
}
