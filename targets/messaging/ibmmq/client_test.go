// +build container

package ibmmq

import (
	"context"
	"github.com/kubemq-io/kubemq-targets/types"
	"io/ioutil"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"

	"github.com/stretchr/testify/require"
	"testing"
)

type testStructure struct {
	applicationChannelName string
	hostname               string
	listenerPort           string
	queueManagerName       string
	apiKey                 string
	mqUsername             string
	password               string
	QueueName              string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/applicationChannelName.txt")
	if err != nil {
		return nil, err
	}
	t.applicationChannelName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/hostname.txt")
	if err != nil {
		return nil, err
	}
	t.hostname = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/listenerPort.txt")
	if err != nil {
		return nil, err
	}
	t.listenerPort = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/queueManagerName.txt")
	if err != nil {
		return nil, err
	}
	t.queueManagerName = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/applicationApiKey/apiKey.txt")
	if err != nil {
		return nil, err
	}
	t.apiKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/applicationApiKey/mqUsername.txt")
	if err != nil {
		return nil, err
	}
	t.mqUsername = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/applicationApiKey/mqPassword.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	t.QueueName = "DEV.QUEUE.1"
	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init incorrect username",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           "fake user name",
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - incorrect password",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           "fake password",
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - incorrect host_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          "fake host name",
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing host_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing queue_manager_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"host_name":      dat.hostname,
					"port_number":    dat.listenerPort,
					"channel_name":   dat.applicationChannelName,
					"username":       dat.mqUsername,
					"key_repository": dat.apiKey,
					"password":       dat.password,
					"queue_name":     dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing channel_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing username",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing queue_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			defer func() {
				_ = c.Stop()
			}()
			require.NoError(t, err)

			err = c.Stop()
			require.NoError(t, err)
		})
	}
}

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid - send",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			request: types.NewRequest().
				SetData([]byte("some-data")),


			wantErr: false,
		}, {
			name: "valid - send override queue",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         "test",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("dynamic_queue", dat.QueueName).
				SetData([]byte("some-data")),


			wantErr: false,
		}, {
			name: "invalid - send missing data",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			request: types.NewRequest(),

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
