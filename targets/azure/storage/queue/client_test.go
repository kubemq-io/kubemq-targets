package queue

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
	storageAccessKey      string
	storageAccount        string
	serviceURL            string
	queueName             string
	queueNameWithMetadata string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/azure/storage/queue/storageAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.storageAccessKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/queue/storageAccount.txt")
	if err != nil {
		return nil, err
	}
	t.storageAccount = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/queue/serviceURL.txt")
	if err != nil {
		return nil, err
	}
	t.serviceURL = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/queue/queueName.txt")
	if err != nil {
		return nil, err
	}
	t.queueName = fmt.Sprintf("%s", dat)
	t.queueNameWithMetadata = fmt.Sprintf("%smetadata", t.queueName)
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
				Name: "azure-storage-queue",
				Kind: "azure.storage.queue",
				Properties: map[string]string{
					"storage_access_key": dat.storageAccessKey,
					"storage_account":    dat.storageAccount,
					"policy":             "retry_policy_exponential",
					"max_tries":          "1",
					"try_timeout":        "10000",
					"retry_delay":        "60",
					"max_retry_delay":    "180",
				},
			},
			wantErr: false,
		}, {
			name: "init - missing account",
			cfg: config.Spec{
				Name: "azure-storage-queue",
				Kind: "azure.storage.queue",
				Properties: map[string]string{
					"storage_access_key": dat.storageAccessKey,
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

func TestClient_Create_Queue(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-queue",
		Kind: "azure.storage.queue",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
			"policy":             "retry_policy_exponential",
			"max_tries":          "2",
			"try_timeout":        "10000",
			"retry_delay":        "60",
			"max_retry_delay":    "180",
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create queue",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "valid create queue with queue metadata",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("queue_metadata", `{"tag":"test","name":"myname"}`).
				SetMetadataKeyValue("queue_name", dat.queueNameWithMetadata).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		},
		{
			name: "invalid create queue - queue already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid create queue - missing queue name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid create queue - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("queue_name", dat.queueName),
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

func TestClient_Push(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	myMessage := "my message to send to queue"
	b, err := json.Marshal(myMessage)
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-queue",
		Kind: "azure.storage.queue",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
			"policy":             "retry_policy_exponential",
			"max_tries":          "1",
			"try_timeout":        "10000",
			"retry_delay":        "60",
			"max_retry_delay":    "180",
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid push item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "push").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(b),
			wantErr: false,
		}, {
			name: "invalid push item - fake queue",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "push").
				SetMetadataKeyValue("queue_name", "fakequeue").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(b),
			wantErr: true,
		}, {
			name: "invalid push item - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "push").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid push item - missing queue name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "push").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(b),
			wantErr: true,
		}, {
			name: "invalid push item - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "push").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetData(b),
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

func TestClient_Peek(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-queue",
		Kind: "azure.storage.queue",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
			"policy":             "retry_policy_exponential",
			"max_tries":          "1",
			"try_timeout":        "10000",
			"retry_delay":        "60",
			"max_retry_delay":    "180",
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid peek",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "peek").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "invalid peek - missing queue name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "peek").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid peek - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "peek").
				SetMetadataKeyValue("queue_name", dat.queueName),
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

			t.Logf("received message :%s", got.Data)
		})
	}
}

func TestClient_GetMessageCount(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-queue",
		Kind: "azure.storage.queue",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
			"policy":             "retry_policy_exponential",
			"max_tries":          "1",
			"try_timeout":        "10000",
			"retry_delay":        "60",
			"max_retry_delay":    "180",
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid message get_messages_count",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_messages_count").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "invalid get_messages_count - missing queue_name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "peek").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid get_messages_count - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "peek").
				SetMetadataKeyValue("queue_name", dat.queueName),
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

			t.Logf("count messages :%s", got.Metadata["count"])
		})
	}
}

func TestClient_Pop(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-queue",
		Kind: "azure.storage.queue",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
			"policy":             "retry_policy_exponential",
			"max_tries":          "1",
			"try_timeout":        "10000",
			"retry_delay":        "60",
			"max_retry_delay":    "180",
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid pop",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "pop").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "invalid pop - missing queue name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "pop").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid pop - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "pop").
				SetMetadataKeyValue("queue_name", dat.queueName),
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

			t.Logf("received message :%s", got.Data)
		})
	}
}

func TestClient_Delete_Queue(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-queue",
		Kind: "azure.storage.queue",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
			"policy":             "retry_policy_exponential",
			"max_tries":          "2",
			"try_timeout":        "10000",
			"retry_delay":        "60",
			"max_retry_delay":    "180",
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create queue",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "valid create queue with queue metadata",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("queue_name", dat.queueNameWithMetadata).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		},
		{
			name: "invalid create queue - queue does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("queue_name", dat.queueName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid create queue - missing queue name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid create queue - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("queue_name", dat.queueName),
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
