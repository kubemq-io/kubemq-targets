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
	dat, err := ioutil.ReadFile("./../../../credentials/azure/eventhubs/endPoint.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/sharedAccessKeyName.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKeyName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/sharedAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKey = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/entityPath.txt")
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
				Kind: "azure.eventhubs",
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

		})
	}
}

func TestClient_Send(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-eventhubs",
		Kind: "azure.eventhubs",
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
		{
			name: "valid send with properties",
			request: types.NewRequest().
				SetMetadataKeyValue("properties", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("method", "send").
				SetData(body),
			wantErr: false,
		}, {
			name: "invalid send missing body",
			request: types.NewRequest().
				SetMetadataKeyValue("properties", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("method", "send"),
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

func TestClient_SendBatch(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-azure-eventhubs",
		Kind: "azure.eventhubs",
		Properties: map[string]string{
			"end_point":              dat.endPoint,
			"shared_access_key_name": dat.sharedAccessKeyName,
			"shared_access_key":      dat.sharedAccessKey,
			"entity_path":            dat.entityPath,
		},
	}
	var messages []string
	m1 := "test1"
	m2 := "test2"
	m3 := "test3"
	m4 := "test4"

	messages = append(messages, m1)
	messages = append(messages, m2)
	messages = append(messages, m3)
	messages = append(messages, m4)

	body, err := json.Marshal(messages)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid send no properties",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send_batch").
				SetData(body),
			wantErr: false,
		}, {
			name: "valid send with properties",
			request: types.NewRequest().
				SetMetadataKeyValue("properties", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("method", "send_batch").
				SetData(body),
			wantErr: false,
		}, {
			name: "invalid send missing body",
			request: types.NewRequest().
				SetMetadataKeyValue("properties", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("method", "send_batch"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
