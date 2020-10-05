package servicebus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"testing"
	"time"
)

type testStructure struct {
	endPoint            string
	sharedAccessKey     string
	sharedAccessKeyName string
	data                []byte
	queue               string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/azure/servicebus/endPoint.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/servicebus/sharedAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKey = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/servicebus/sharedAccessKeyName.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKeyName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/servicebus/message.txt")
	if err != nil {
		return nil, err
	}
	t.data = dat
	t.queue = "myqueue"
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
			name: "init ",
			cfg: config.Spec{
				Name: "target-azure-servicebus",
				Kind: "target.azure.servicebus",
				Properties: map[string]string{
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
					"queue_name":             dat.queue,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing shared_access_key_name",
			cfg: config.Spec{
				Name: "target-azure-servicebus",
				Kind: "target.azure.servicebus",
				Properties: map[string]string{
					"end_point":         dat.endPoint,
					"shared_access_key": dat.sharedAccessKey,
					"queue_name":        dat.queue,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing shared_access_key",
			cfg: config.Spec{
				Name: "target-azure-servicebus",
				Kind: "target.azure.servicebus",
				Properties: map[string]string{
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"queue_name":             dat.queue,
				},
			},
			wantErr: true,
		},
		{
			name: "init - missing queue_name",
			cfg: config.Spec{
				Name: "target-azure-servicebus",
				Kind: "target.azure.servicebus",
				Properties: map[string]string{
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
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
			require.NoError(t, err)

		})
	}
}

func TestClient_Send_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-servicebus",
		Kind: "target.azure.servicebus",
		Properties: map[string]string{
			"end_point":              dat.endPoint,
			"shared_access_key_name": dat.sharedAccessKeyName,
			"shared_access_key":      dat.sharedAccessKey,
			"queue_name":             dat.queue,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid send item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send").
				SetData(dat.data),
			wantErr: false,
		}, {
			name: "invalid send item missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, cfg)
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

func TestClient_Send_Batch_Items(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	var data []string
	data = append(data, "message1")
	data = append(data, "message2")
	b, err := json.Marshal(data)
	cfg := config.Spec{
		Name: "target-azure-servicebus",
		Kind: "target.azure.servicebus",
		Properties: map[string]string{
			"end_point":              dat.endPoint,
			"shared_access_key_name": dat.sharedAccessKeyName,
			"shared_access_key":      dat.sharedAccessKey,
			"queue_name":             dat.queue,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid send batch item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send_batch").
				SetMetadataKeyValue("label", "my_label").
				SetMetadataKeyValue("content_type", "content_type").
				SetData(b),
			wantErr: false,
		}, {
			name: "invalid send batch item missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send_batch"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, cfg)
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
