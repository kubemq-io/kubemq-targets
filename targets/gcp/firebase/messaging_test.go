package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
		},
		Data: map[string]string{"key": "val"},
	}

	multib, err := json.Marshal(multi)
	if err != nil {
		return
	}
	fmt.Print(string(multib))

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
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendMessage").
				SetData(mb),
			wantmsg: messages{single: m},
			wantErr: true,
		},
		{
			name:    "parse multicast msg",
			isMulti: true,
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendBatch").
				SetData(multib),
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
				m, err := parseMetadataMessages(tt.request.Data, options{
					defaultMessaging: &messages{},
				}, SendBatch)
				require.NoError(t, err)
				require.EqualValues(t, tt.wantmsg.multicast, m.multicast)

			} else {
				m, err := parseMetadataMessages(tt.request.Data, options{
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
					"messaging_client": "true",
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
func TestDefultMessage(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantmsg *messages
		wantErr bool
	}{
		{
			name: "missing data on SendMessage",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messaging_client": "true",
					"defaultmsg":      `{"topic":"defult"}`,
				},
			},
			wantmsg: &messages{single: &messaging.Message{
				Topic: "defult",
				Data:  map[string]string{"key1": "val1"},
			},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendMessage"),
			wantErr: true,
		},
		{
			name: "combine default and SendMessage",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messaging_client": "true",
					"defaultmsg":      `{"topic":"defult"}`,
				},
			},
			wantmsg: &messages{single: &messaging.Message{
				Topic: "defult",
				Data:  map[string]string{"key1": "val1"},
			},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendMessage").SetData([]byte(`{"Topic":"defult","data":{"key1":"val1"}}`)),
			wantErr: false,
		},
		{
			name: "combine and replace defult SendMessage",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messaging_client": "true",
					"defaultmsg":      `{"Topic":"defult","token":"1234"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendMessage").
				SetData([]byte(`{"Topic":"newTopic"}`)),
			wantmsg: &messages{single: &messaging.Message{
				Topic: "newTopic",
				Token: "1234",
			},
			},
			wantErr: false,
		},
		{
			name: "replace defult SendMessage",
			cfg: config.Spec{
				Name: "test",
				Kind: "test",
				Properties: map[string]string{
					"project_id":      "123",
					"credentials":     "noc",
					"messaging_client": "true",
					"defaultmsg":      `{"Topic":"defult"}`,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendMessage").
				SetData([]byte(`{"Topic":"newTopic"}`)),
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
			m, err := parseMetadataMessages(tt.request.Data, o, SendMessage)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, tt.wantmsg, m)
			}
			//defult message not changed
			oc, _ := parseOptions(tt.cfg)
			require.EqualValues(t, oc.defaultMessaging, o.defaultMessaging)

		})
	}
}

func TestClientDo(t *testing.T) {

	cred := `{
	  }`

	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "test do",
			cfg: config.Spec{
				Kind: "test",
				Name: "test",
				Properties: map[string]string{
					"project_id":      "pubsubdemo-281010",
					"credentials":     cred,
					"messaging_client": "true",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "SendMessage").
				SetData([]byte(`{"Topic":"test"}`)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			res, err := c.Do(ctx, tt.request)
			if err != nil {
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
			}
			fmt.Print(res)
		})
	}
}
