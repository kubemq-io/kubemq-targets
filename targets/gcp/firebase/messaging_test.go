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

func TestClientDo(t *testing.T) {

	cred := `{
		"type": "service_account",
		"project_id": "pubsubdemo-281010",
		"private_key_id": "997d2ab5df9f2e869e15c2a550b111b6a89cec20",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC40ai1DNA4lE65\njn/Pjnt3pr7AtT11nRukjCGw7q3TevQrKT5KQ75Z+eo+sf8G6MnBTy+tK6EGt5ed\nKFUdngXnh4Ipo7JI3rfwGX22hMnC8iTnCIZygHy5ArA+2LpkqQTveQnbz+MFchrI\nEdg9LlG1p0oYSQMt4g3AOhl9mcjLbkLxRlDi8+xM4G8e2mm5+QiscwXV9I6AC/Fh\nriPXs/1Vw3nGdjHajYqdQj8Uc63kBB1x5MsA+YFyFZl1bU8ViapatUrG9eDsMr8o\nXbA+uoJ2p6nJGCd8ISixLh4+YzIvAYclM1Tj7U05AwYpC4e8DAXyjyLYHbiZV7Lc\nFgba2JF9AgMBAAECggEACTSv0hYtBoOO2r9iuSxwMS5EPVfbhVf0aSTDTYRacobb\n0hs3MpPmelCFwmZ6izmdu4MzEM2ez4rJpiGD5EibXJJ9Mx2nwyI3+gvvcpiN4G8I\n8E4yrm5kTx+FxRjz9Z8OTER0i95V8s7hsvTwJUVWbLHgqo46r5arjQIGEncHIVpj\n6KZJCz1/o/pUAg2YOqTy8EuFPyvtwolSdLITA37tieWnh0D8lE5Dj5UwmDBm0YQi\nybh+ey1BaTG8DIS/Wy/dRliZDKoAd3st40ypr5vYVOoF9kk3Opl8tgB0xQAheSVb\nXsJa2EUBDPS8IwymjYu3ijvQQ9nCakdnlOg39eVlMwKBgQD8KANbU1Ia4yTuNZMy\n350aAs4e9trWcgGwUIrrOw+yKt/ugcxjao9y0il7OiVtZ5V1r2VQn0/GiLrwGTDB\nell6eoI/g99qn/KPsXt6RT7S+OYbCxUo0EMPyMgwPBSlFERY2rt10jAO7Ok++tJh\nYuvxf9JMeNIS5ThD19CCXm3YNwKBgQC7ouBNXZdWFCwOxd8NUqUAhPe2VKMQJtLd\n9j7EBfi0Gemqsk4KIOlwb3H+07AilDBhvXlpINm1dNjIyaWZvraqlvntFFyS/06p\nAdN5FjxG+lhfc4b+zs4YvR9LNDRClPJoBa3IX7CDjSjtkVAlQ+0EHBsZzEdGJdBi\n7+FLu/wh6wKBgEVCdFGUXDv4Yf9wBcN2ej9Xv+fvZAJ9BAu6w72C1nfYoPNxAYPZ\nFBe0tCIdwYQAbKQLjieL6qych8RFFwg9o/ApUDdD8Izn7Acd982I0Y2/QezxqVkx\ngwoF2z6scfs5yuAhDFZ7ainfVt2upTSMqEQIGOpaUVFRVpgD4ki8yS0XAoGAPBEk\nJSA09kV25TPK+ATg9Y2bjy8BFIaZMp1F8pLGz0EMYKy79toaYPgMUjuKQ0eVRXTW\njSULDN/fFkgXT2SSLYIveAnwqM46bDg9bqIDoeU6rTPan2+s4paIkhagNEBiaZKH\n04FujG6AD61ZLtTT52Dn/BY9KuOoFkQcp5YCXQkCgYAR2B7t1ntfrW4YbtGSLly2\nEtn2FW+JukaRuZAO0i4mgemoN/qERZh3jybgPjCmXvDQGDOHMnUUVwf0neQKQzzs\n6VapjhxJwpM2H7jaVUv5GDrXAZpkkN+OAdc/O9WgpgEGecdKZzUKIwwW/D0/lnC8\nO5JORl6i10mk/SyD+yo8Iw==\n-----END PRIVATE KEY-----\n",
		"client_email": "pubsub@pubsubdemo-281010.iam.gserviceaccount.com",
		"client_id": "107517457611651099401",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/pubsub%40pubsubdemo-281010.iam.gserviceaccount.com"
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
					"messagingClient": "true",
				},
			},
			request: types.NewRequest().SetMetadataKeyValue("message", `{"Topic":"test"}`).SetMetadataKeyValue("method", "SendMessage"),
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
