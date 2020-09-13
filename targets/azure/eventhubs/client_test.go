package eventhubs

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
	sharedAccessKeyName string
	sharedAccessKey     string
	entityPath          string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/azure/eventhubs/endPoint.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/eventhubs/sharedAccessKeyName.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKeyName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/eventhubs/sharedAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKey = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/eventhubs/entityPath.txt")
	if err != nil {
		return nil, err
	}
	t.entityPath = fmt.Sprintf("%s", dat)

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
				Name: "target-azure-eventhubs",
				Kind: "target.azure.eventhubs",
				Properties: map[string]string{
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
					"entity_path":            dat.entityPath,
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

			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Send(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-eventhubs",
		Kind: "target.azure.eventhubs",
		Properties: map[string]string{
			"end_point":              dat.endPoint,
			"shared_access_key_name": dat.sharedAccessKeyName,
			"shared_access_key":      dat.sharedAccessKey,
			"entity_path":            dat.entityPath,
		},
	}
	body, err := json.Marshal("test")
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid send no properties",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send").
				SetData(body),
			wantErr: false,
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

func TestClient_Get_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-eventhubs",
		Kind: "target.azure.eventhubs",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid get item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "valid get item with offset and count",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("offset", "2").
				SetMetadataKeyValue("count", "3").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "invalid get item - missing file_name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid get item - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("file_name", dat.fileName),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
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
			require.NotNil(t, got.Data)
		})
	}
}

func TestClient_Delete_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-eventhubs",
		Kind: "target.azure.eventhubs",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetMetadataKeyValue("delete_snapshots_option_type", ""),
			wantErr: false,
		},
		{
			name: "invalid delete file does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetMetadataKeyValue("delete_snapshots_option_type", ""),
			wantErr: true,
		},
		{
			name: "invalid delete- missing file_name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetMetadataKeyValue("delete_snapshots_option_type", ""),
			wantErr: true,
		},
		{
			name: "invalid delete- fake option",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("delete_snapshots_option_type", ""),
			wantErr: true,
		}, {
			name: "invalid delete- fake url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", "fakeURL").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("delete_snapshots_option_type", ""),
			wantErr: true,
		},
		{
			name: "invalid delete- invalid delete_snapshots_option_type",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("delete_snapshots_option_type", "test"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
