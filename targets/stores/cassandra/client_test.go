package cassandra

import (
	"context"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad hosts",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad port",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "-1",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			wantErr: true,
		},
		{
			name: "init - pad proto version",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "-1",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad replication factor",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "-1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad consistency",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "bad-value",
					"default_table":      "test",
					"default_keyspace":   "test",
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
				t.Errorf("Init() error = %v, wantExecErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Set_Get(t *testing.T) {
	tests := []struct {
		name            string
		cfg             config.Metadata
		setRequest      *types.Request
		getRequest      *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "strong"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
				SetData([]byte("some-data")),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "update set get request",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "eventual").
				SetData([]byte("some-data-2")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "eventual"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
				SetData([]byte("some-data-2")),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "invalid set",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "strong").
				SetMetadataKeyValue("table", "bad-table").
				SetMetadataKeyValue("keyspace", "bad-keyspace").
				SetData([]byte("some-data")),
			getRequest:      nil,
			wantSetResponse: nil,
			wantGetResponse: nil,
			wantSetErr:      true,
			wantGetErr:      false,
		},
		{
			name: "valid set - bad get table",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("table", "bad-table").
				SetMetadataKeyValue("keyspace", "bad-keyspace").
				SetMetadataKeyValue("consistency", "strong"),
			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: nil,
			wantSetErr:      false,
			wantGetErr:      true,
		},
		{
			name: "valid set - empty result",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "not-exist-key").
				SetMetadataKeyValue("consistency", "strong"),
			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
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

func TestClient_Query_Exec(t *testing.T) {
	tests := []struct {
		name              string
		cfg               config.Metadata
		execRequest       *types.Request
		queryRequest      *types.Request
		wantExecResponse  *types.Response
		wantQueryResponse *types.Response
		wantExecErr       bool
		wantQueryErr      bool
	}{
		{
			name: "valid exec query request",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte(`INSERT INTO test.test (key, value) VALUES ('some-key',textAsBlob('some-data'))`)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte(`SELECT value FROM test.test WHERE key = 'some-key'`)),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetData([]byte("some-data")),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "invalid exec request - empty",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "eventual").
				SetData(nil),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid exec request",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "eventual").
				SetData([]byte("some bad exec query")),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid query request - empty",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte(`INSERT INTO test.test (key, value) VALUES ('some-key',textAsBlob('some-data'))`)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("consistency", "strong").
				SetData(nil),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
		{
			name: "invalid query request - empty",
			cfg: config.Metadata{
				Name: "cassandra-target",
				Kind: "",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "eventual").
				SetData([]byte(`INSERT INTO test.test (key, value) VALUES ('some-key',textAsBlob('some-data'))`)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("consistency", "eventual").
				SetData([]byte("some bad query")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.execRequest)
			if tt.wantExecErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantExecResponse, gotSetResponse)
			gotGetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantQueryErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.EqualValues(t, tt.wantQueryResponse, gotGetResponse)
		})
	}
}
func TestClient_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := New()
	err := c.Init(ctx, config.Metadata{
		Name: "target.cassandra",
		Kind: "target.cassandra",
		Properties: map[string]string{
			"hosts":              "localhost",
			"port":               "9042",
			"username":           "cassandra",
			"password":           "cassandra",
			"proto_version":      "4",
			"replication_factor": "1",
			"consistency":        "All",
			"default_table":      "test",
			"default_keyspace":   "test",
		},
	})
	key := nuid.Next()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("key", key).
		SetMetadataKeyValue("consistency", "strong").
		SetData([]byte("some-data"))

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("key", key).
		SetMetadataKeyValue("consistency", "strong")
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, []byte("some-data"), gotGetResponse.Data)

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("consistency", "strong").
		SetMetadataKeyValue("key", key)
	_, err = c.Do(ctx, delRequest)
	require.NoError(t, err)
	gotGetResponse, err = c.Do(ctx, getRequest)
	require.Error(t, err)
	require.Nil(t, gotGetResponse)

	delRequest = types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("consistency", "").
		SetMetadataKeyValue("table", "bad-table").
		SetMetadataKeyValue("keyspace", "bad-keyspace").
		SetMetadataKeyValue("key", key)
	_, err = c.Do(ctx, delRequest)
	require.Error(t, err)

}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Metadata{
				Name: "target.cassandra",
				Kind: "target.cassandra",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", nuid.Next()).
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte("some-data")),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Metadata{
				Name: "target.cassandra",
				Kind: "target.cassandra",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - no key",
			cfg: config.Metadata{
				Name: "target.cassandra",
				Kind: "target.cassandra",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("consistency", "strong").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - bad consistency key",
			cfg: config.Metadata{
				Name: "target.cassandra",
				Kind: "target.cassandra",
				Properties: map[string]string{
					"hosts":              "localhost",
					"port":               "9042",
					"username":           "cassandra",
					"password":           "cassandra",
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "All",
					"default_table":      "test",
					"default_keyspace":   "test",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("consistency", "not-valid").
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
