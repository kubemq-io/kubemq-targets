package elastic

import (
	"context"
	"encoding/json"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const mapping = `{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"id": {
				"type": "keyword"
			},
			"data": {
				"type": "text"
			}
		}
	}
}`

const (
	elasticUrl = "http://localhost:9200"
	testIndex  = "log"
)

type log struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func (l *log) marshal() []byte {
	b, _ := json.Marshal(l)
	return b
}

func newLog(id, data string) *log {
	return &log{Id: id,
		Data: data}
}
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := elastic.NewClient(elastic.SetURL(elasticUrl), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	_, _, err = client.Ping(elasticUrl).Do(ctx)
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(testIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	if exists {
		_, err := client.DeleteIndex(testIndex).Do(ctx)
		if err != nil {
			panic(err)
		}

	}

	_, err = client.CreateIndex(testIndex).BodyString(mapping).Do(ctx)
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad urls",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":                    "scheme://localhost:9200",
					"sniff":                   "false",
					"username":                "",
					"password":                "",
					"retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad options - no urls",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":                    "",
					"sniff":                   "false",
					"username":                "",
					"password":                "",
					"retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_Set_Get(t *testing.T) {
	tests := []struct {
		name            string
		cfg             config.Spec
		setRequest      *types.Request
		getRequest      *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("id", "some-id").
				SetMetadataKeyValue("index", "log").
				SetData(newLog("some-id", "some-data").marshal()),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("index", "log").
				SetMetadataKeyValue("id", "some-id"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("id", "some-id").
				SetMetadataKeyValue("result", "created"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("id", "some-id").
				SetData(newLog("some-id", "some-data").marshal()),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "update set get request",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("id", "some-id").
				SetMetadataKeyValue("index", "log").
				SetData(newLog("some-id", "some-data-2").marshal()),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("index", "log").
				SetMetadataKeyValue("id", "some-id"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("id", "some-id").
				SetMetadataKeyValue("result", "updated"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("id", "some-id").
				SetData(newLog("some-id", "some-data-2").marshal()),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "bad set - invalid index",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("id", "some-id").
				SetMetadataKeyValue("index", "bad-index").
				SetData(nil),
			getRequest: nil,

			wantSetResponse: nil,
			wantGetResponse: nil,
			wantSetErr:      true,
			wantGetErr:      false,
		},
		{
			name: "valid set - not found index",
			cfg: config.Spec{
				Name: "elastic-target",
				Kind: "",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("id", "some-id-2").
				SetMetadataKeyValue("index", "log").
				SetData(newLog("some-id-2", "some-data").marshal()),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("index", "log").
				SetMetadataKeyValue("id", "bad-id"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("id", "some-id-2").
				SetMetadataKeyValue("result", "created"),
			wantGetResponse: nil,
			wantSetErr:      false,
			wantGetErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.setRequest)
			if tt.wantSetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantSetResponse, gotSetResponse)
			gotGetResponse, err := c.Do(ctx, tt.getRequest)
			if tt.wantGetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.EqualValues(t, tt.wantGetResponse, gotGetResponse)
		})
	}
}

func TestClient_Delete(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := New()
	err := c.Init(ctx, config.Spec{
		Name: "elastic-target",
		Kind: "",
		Properties: map[string]string{
			"urls":     "http://localhost:9200",
			"sniff":    "false",
			"username": "",
			"password": "",
		},
	}, nil)
	key := uuid.New().String()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("index", testIndex).
		SetMetadataKeyValue("id", key).
		SetData(newLog(key, "some-data").marshal())

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("index", testIndex).
		SetMetadataKeyValue("id", key)
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, newLog(key, "some-data").marshal(), gotGetResponse.Data)

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("index", testIndex).
		SetMetadataKeyValue("id", key)
	_, err = c.Do(ctx, delRequest)
	require.NoError(t, err)
	gotGetResponse, err = c.Do(ctx, getRequest)
	require.Error(t, err)
	require.Nil(t, gotGetResponse)
	delRequest = types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("index", "bad-index").
		SetMetadataKeyValue("id", key)
	_, err = c.Do(ctx, delRequest)
	require.Error(t, err)
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "bad request - bad method",
			cfg: config.Spec{
				Name: "elastic",
				Kind: "elastic",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad method").
				SetMetadataKeyValue("id", "some-id"),
			wantErr: true,
		},
		{
			name: "bad request - no index",
			cfg: config.Spec{
				Name: "elastic",
				Kind: "elastic",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("id", "some-id"),
			wantErr: true,
		},
		{
			name: "bad request - no id",
			cfg: config.Spec{
				Name: "elastic",
				Kind: "elastic",
				Properties: map[string]string{
					"urls":     "http://localhost:9200",
					"sniff":    "false",
					"username": "",
					"password": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("index", testIndex),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			_, err = c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
