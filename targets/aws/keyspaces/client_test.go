package keyspaces

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type testStructure struct {
	username string
	password string
	endPoint string
	tlsPath  string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/aws/keyspaces/username.txt")
	if err != nil {
		return nil, err
	}
	t.username = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/keyspaces/password.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/keyspaces/endPoint.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = string(dat)

	t.tlsPath = "https://www.amazontrust.com/repository/AmazonRootCA1.pem"

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
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad hosts",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              "",
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad port",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "-1",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			wantErr: true,
		},
		{
			name: "init - pad proto version",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "-1",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad replication factor",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "-1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad consistency",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "bad-value",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
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

		})
	}
}

func TestClient_Set_Get(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
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
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum"),

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
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data-2")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum"),

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
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum").
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
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("table", "bad-table").
				SetMetadataKeyValue("keyspace", "bad-keyspace").
				SetMetadataKeyValue("consistency", "local_quorum"),
			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: nil,
			wantSetErr:      false,
			wantGetErr:      true,
		},
		{
			name: "valid set - empty result",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-id").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "not-exist-key").
				SetMetadataKeyValue("consistency", "local_quorum"),
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
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name              string
		cfg               config.Spec
		execRequest       *types.Request
		queryRequest      *types.Request
		wantExecResponse  *types.Response
		wantQueryResponse *types.Response
		wantExecErr       bool
		wantQueryErr      bool
	}{
		{
			name: "valid exec query request",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte(`INSERT INTO test.test (key, value) VALUES (textAsBlob('some-key'), textAsBlob('some-data'));`)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte(`SELECT value FROM test.test WHERE key = textAsBlob('some-key')`)),
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
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData(nil),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid exec request",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some bad exec query")),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid query request - empty",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte(`INSERT INTO test.test (key, value) VALUES (textAsBlob('some-key'),textAsBlob('some-data'))`)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData(nil),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
		{
			name: "invalid query request - empty",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte(`INSERT INTO test.test (key, value) VALUES (textAsBlob('some-key'),textAsBlob('some-data'))`)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("consistency", "local_quorum").
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
	dat, err := getTestStructure()
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := New()
	err = c.Init(ctx, config.Spec{
		Name: "target-aws-keyspaces",
		Kind: "target.aws.keyspaces",
		Properties: map[string]string{
			"hosts":              dat.endPoint,
			"port":               "9142",
			"username":           dat.username,
			"password":           dat.password,
			"proto_version":      "4",
			"replication_factor": "1",
			"consistency":        "local_quorum",
			"default_table":      "test",
			"default_keyspace":   "test",
			"tls":                dat.tlsPath,
		},
	})
	key := nuid.Next()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("key", key).
		SetMetadataKeyValue("consistency", "local_quorum").
		SetData([]byte("some-data"))

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("key", key).
		SetMetadataKeyValue("consistency", "local_quorum")
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, []byte("some-data"), gotGetResponse.Data)

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("consistency", "local_quorum").
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
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", nuid.Next()).
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data")),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - no key",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("consistency", "local_quorum").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - bad consistency key",
			cfg: config.Spec{
				Name: "target-aws-keyspaces",
				Kind: "target.aws.keyspaces",
				Properties: map[string]string{
					"hosts":              dat.endPoint,
					"port":               "9142",
					"username":           dat.username,
					"password":           dat.password,
					"proto_version":      "4",
					"replication_factor": "1",
					"consistency":        "local_one",
					"default_table":      "test",
					"default_keyspace":   "test",
					"tls":                dat.tlsPath,
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
