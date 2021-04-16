package mongodb

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testDocument struct {
	Id      string `json:"_id,omitempty"`
	String  string `json:"string,omitempty"`
	Integer string `json:"integer,omitempty"`
}

func (t *testDocument) data() []byte {
	data, _ := json.Marshal(t)
	return data
}
func (t *testDocument) dataString() string {
	data, _ := json.MarshalToString(t)
	return data
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
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "majority",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error connection",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "bad-host:32017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "majority",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad host",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad database",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad collection",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad write concurrency",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "bad-concurrency",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad read concurrency",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "bad-concurrency",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad operation timeout",
			cfg: config.Spec{
				Name: "mongodb-target",
				Kind: "",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "-2",
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

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
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
				Name: "mongodb",
				Kind: "mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set_by_key").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_by_key").
				SetMetadataKeyValue("key", "some-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "valid set , no key get request",
			cfg: config.Spec{
				Name: "mongodb",
				Kind: "mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set_by_key").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_by_key").
				SetMetadataKeyValue("key", "bad-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
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
			err := c.Init(ctx, tt.cfg)
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
		Name: "mongodb",
		Kind: "mongodb",
		Properties: map[string]string{
			"host":                      "localhost:27017",
			"username":                  "admin",
			"password":                  "password",
			"database":                  "admin",
			"collection":                "test",
			"write_concurrency":         "",
			"read_concurrency":          "",
			"params":                    "",
			"operation_timeout_seconds": "2",
		},
	})
	key := uuid.New().String()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set_by_key").
		SetMetadataKeyValue("key", key).
		SetData([]byte("some-data"))

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get_by_key").
		SetMetadataKeyValue("key", key)
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, []byte("some-data"), gotGetResponse.Data)

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete_by_key").
		SetMetadataKeyValue("key", key)
	_, err = c.Do(ctx, delRequest)
	require.NoError(t, err)
	gotGetResponse, err = c.Do(ctx, getRequest)
	require.Error(t, err)
	require.Nil(t, gotGetResponse)
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Spec{
				Name: "mongodb",
				Kind: "mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set_by_key").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Spec{
				Name: "mongodb",
				Kind: "mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - no key",
			cfg: config.Spec{
				Name: "mongodb",
				Kind: "mongodb",
				Properties: map[string]string{
					"host":                      "localhost:27017",
					"username":                  "admin",
					"password":                  "password",
					"database":                  "admin",
					"collection":                "test",
					"write_concurrency":         "",
					"read_concurrency":          "",
					"params":                    "",
					"operation_timeout_seconds": "2",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set_by_key").
				SetData([]byte("some-data")),
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
			_, err = c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestClient_Insert_Find_Update_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := New()
	err := c.Init(ctx, config.Spec{
		Name: "mongodb",
		Kind: "mongodb",
		Properties: map[string]string{
			"host":                      "localhost:27017",
			"username":                  "admin",
			"password":                  "password",
			"database":                  "admin",
			"collection":                "test",
			"write_concurrency":         "",
			"read_concurrency":          "",
			"params":                    "",
			"operation_timeout_seconds": "2",
		},
	})

	require.NoError(t, err)
	doc := &testDocument{
		Id:      uuid.New().String(),
		String:  "s",
		Integer: "1",
	}
	var docs []*testDocument
	docs = append(docs, doc)
	docsData, err := json.Marshal(&docs)
	require.NoError(t, err)
	insertRequest := types.NewRequest().
		SetMetadataKeyValue("method", "insert_many").
		SetData(docsData)
	fmt.Println(insertRequest.String())
	_, err = c.Do(ctx, insertRequest)
	require.NoError(t, err)

	findFilter := &testDocument{
		Id: doc.Id,
	}

	findRequest := types.NewRequest().
		SetMetadataKeyValue("method", "find").
		SetMetadataKeyValue("filter", findFilter.dataString())

	findResponse, err := c.Do(ctx, findRequest)
	require.NoError(t, err)
	require.NotNil(t, findResponse)
	receivedDoc := &testDocument{}
	err = json.Unmarshal(findResponse.Data, receivedDoc)
	require.NoError(t, err)
	require.EqualValues(t, doc, receivedDoc)

	updateFilter := &testDocument{
		Id: doc.Id,
	}
	updateDoc := &testDocument{
		Id:      doc.Id,
		String:  "b",
		Integer: "2",
	}
	updateRequest := types.NewRequest().
		SetMetadataKeyValue("method", "update").
		SetMetadataKeyValue("filter", updateFilter.dataString()).SetData(updateDoc.data())
	updateResponse, err := c.Do(ctx, updateRequest)
	require.NoError(t, err)
	require.NotNil(t, updateResponse)

	findResponse, err = c.Do(ctx, findRequest)
	require.NoError(t, err)
	require.NotNil(t, findResponse)
	receivedDoc2 := &testDocument{}
	err = json.Unmarshal(findResponse.Data, receivedDoc2)
	require.NoError(t, err)
	require.EqualValues(t, updateDoc, receivedDoc2)
	upsertDoc := &testDocument{
		Id:      uuid.New().String(),
		String:  "c",
		Integer: "3",
	}
	upsertFilter := &testDocument{
		Id: upsertDoc.Id,
	}

	upsertRequest := types.NewRequest().
		SetMetadataKeyValue("method", "update").
		SetMetadataKeyValue("set_upsert", "true").
		SetMetadataKeyValue("filter", upsertFilter.dataString()).SetData(upsertDoc.data())

	upsertResponse, err := c.Do(ctx, upsertRequest)
	require.NoError(t, err)
	require.NotNil(t, upsertResponse)

	findUpsertFilter := &testDocument{
		Id: upsertDoc.Id,
	}
	findUpsertRequest := types.NewRequest().
		SetMetadataKeyValue("method", "find").
		SetMetadataKeyValue("filter", findUpsertFilter.dataString())
	findUpsertResponse, err := c.Do(ctx, findUpsertRequest)
	require.NoError(t, err)
	require.NotNil(t, findUpsertResponse)
	receivedUpsertDoc2 := &testDocument{}
	err = json.Unmarshal(findUpsertResponse.Data, receivedUpsertDoc2)
	require.NoError(t, err)
	require.EqualValues(t, upsertDoc, receivedUpsertDoc2)

	delFilter := &testDocument{
		Id: upsertDoc.Id,
	}
	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("filter", delFilter.dataString())
	delResponse, err := c.Do(ctx, delRequest)
	require.NoError(t, err)
	require.NotNil(t, delResponse)
	delFilter.Id = doc.Id
	delRequest = types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("filter", delFilter.dataString())
	delResponse, err = c.Do(ctx, delRequest)
	require.NoError(t, err)
	require.NotNil(t, delResponse)

}
