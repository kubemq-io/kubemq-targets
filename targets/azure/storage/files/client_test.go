package files

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"testing"
	"time"
)

type testStructure struct {
	storageAccessKey   string
	storageAccount     string
	serviceURL         string
	serviceURLWithMeta string
	file               []byte
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/azure/storage/files/storageAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.storageAccessKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/files/storageAccount.txt")
	if err != nil {
		return nil, err
	}
	t.storageAccount = string(dat)
	contents, err := ioutil.ReadFile("./../../../../credentials/azure/storage/files/contents.txt")
	if err != nil {
		return nil, err
	}
	t.file = contents
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/files/serviceURL.txt")
	if err != nil {
		return nil, err
	}
	t.serviceURL = string(dat)
	t.serviceURLWithMeta = fmt.Sprintf("%sWithMetadata.txt", t.serviceURL)
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
				Name: "azure-storage-files",
				Kind: "azure.storage.files",
				Properties: map[string]string{
					"storage_access_key": dat.storageAccessKey,
					"storage_account":    dat.storageAccount,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing storage_account",
			cfg: config.Spec{
				Name: "azure-storage-files",
				Kind: "azure.storage.files",
				Properties: map[string]string{
					"storage_access_key": dat.storageAccessKey,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing storage_access_key",
			cfg: config.Spec{
				Name: "azure-storage-files",
				Kind: "azure.storage.files",
				Properties: map[string]string{
					"storage_account": dat.storageAccount,
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

			err := c.Init(ctx, tt.cfg, nil)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)

		})
	}
}

func TestClient_Create_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-files",
		Kind: "azure.storage.files",
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
			name: "valid create item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(dat.file),
			wantErr: false,
		}, {
			name: "valid create item with metadata",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("file_metadata", `{"tag":"test","name":"myname"}`).
				SetMetadataKeyValue("service_url", dat.serviceURLWithMeta).
				SetData(dat.file),
			wantErr: false,
		}, {
			name: "invalid create item - missing service url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetData(dat.file),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, cfg, nil)
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

func TestClient_Upload_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-files",
		Kind: "azure.storage.files",
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
			name: "valid upload item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(dat.file),
			wantErr: false,
		}, {
			name: "valid upload item with metadata",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("blob_metadata", `{"tag":"test","name":"myname"}`).
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(dat.file),
			wantErr: false,
		}, {
			name: "invalid upload item - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetData(dat.file),
			wantErr: true,
		}, {
			name: "invalid upload item - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, cfg, nil)
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
		Name: "azure-storage-files",
		Kind: "azure.storage.files",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
		},
	}
	tests := []struct {
		name             string
		request          *types.Request
		wantErr          bool
		wantFileMetadata bool
	}{
		{
			name: "valid get item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "valid get item - with file metadata",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("service_url", dat.serviceURLWithMeta),
			wantFileMetadata: true,
			wantErr:          false,
		}, {
			name: "valid get item with offset and count",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("offset", "2").
				SetMetadataKeyValue("count", "3").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "invalid get item - missing service_url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, cfg, nil)
			require.NoError(t, err)
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got.Data)
			if tt.wantFileMetadata {
				require.NotNil(t, got.Metadata["file_metadata"])
				t.Logf("%s", got.Metadata["file_metadata"])
			} else {
				require.Equal(t, got.Metadata["file_metadata"], "")
			}
		})
	}
}

func TestClient_Delete_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "azure-storage-files",
		Kind: "azure.storage.files",
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
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: false,
		}, {
			name: "valid delete file with tags",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURLWithMeta),
			wantErr: false,
		},
		{
			name: "invalid delete file does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		},
		{
			name: "invalid delete- fake option",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		}, {
			name: "invalid delete- fake url",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", "fakeURL"),
			wantErr: true,
		},
		{
			name: "invalid delete- invalid delete_snapshots_option_type",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, cfg, nil)
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
