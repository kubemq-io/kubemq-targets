package blob

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"testing"
	"time"
)

type testStructure struct {
	storageAccessKey string
	storageAccount   string
	fileName         string
	serviceURL       string
	file             []byte
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/azure/storage/storageAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.storageAccessKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/storageAccount.txt")
	if err != nil {
		return nil, err
	}
	t.storageAccount = fmt.Sprintf("%s", dat)
	contents, err := ioutil.ReadFile("./../../../../credentials/azure/storage/contents.txt")
	if err != nil {
		return nil, err
	}
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/fileName.txt")
	if err != nil {
		return nil, err
	}
	t.fileName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/storage/serviceURL.txt")
	if err != nil {
		return nil, err
	}
	t.serviceURL = fmt.Sprintf("%s", dat)

	t.file = contents
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
				Name: "target-azure-storage-blob",
				Kind: "target.azure.storage.blob",
				Properties: map[string]string{
					"storage_access_key": dat.storageAccessKey,
					"storage_account":    dat.storageAccount,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing account",
			cfg: config.Spec{
				Name: "target-azure-storage-blob",
				Kind: "target.azure.storage.blob",
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
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Upload_Item(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-storage-blob",
		Kind: "target.azure.storage.blob",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
		},
	}
	tests := []struct {
		name        string
		request     *types.Request
		wantErr     bool
	}{
		{
			name: "valid upload item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("service_url", dat.serviceURL).
				SetData(dat.file),
			wantErr:     false,
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
		Name: "target-azure-storage-blob",
		Kind: "target.azure.storage.blob",
		Properties: map[string]string{
			"storage_access_key": dat.storageAccessKey,
			"storage_account":    dat.storageAccount,
		},
	}
	tests := []struct {
		name        string
		request     *types.Request
		wantErr     bool
	}{
		{
			name: "valid upload item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("file_name", dat.fileName).
				SetMetadataKeyValue("service_url", dat.serviceURL),
			wantErr:     false,
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
