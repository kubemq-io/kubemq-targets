package azuresql

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type testStructure struct {
	connectionString string
	connectionStringBadPort string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/azure/stores/azuresql/connectionString.txt")
	if err != nil {
		return nil, err
	}
	t.connectionString = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/azure/stores/azuresql/connectionStringBadPort.txt")
	if err != nil {
		return nil, err
	}
	t.connectionStringBadPort = string(dat)
	return t, nil
}

type post struct {
	Id        int64  `json:"id"`
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	BigNumber int64  `json:"bignumber,omitempty"`
	BoolValue bool   `json:"boolvalue"`
}
type posts []*post

func (p *posts) marshal() []byte {
	b, _ := json.Marshal(p)
	return b
}
func unmarshal(data []byte) *posts {
	if data == nil {
		return nil
	}
	p := &posts{}
	_ = json.Unmarshal(data, p)
	return p
}

var allPosts = posts{
	&post{
		Id:        0,
		Content:   "Content One",
		BigNumber: 1231241241231231123,
		BoolValue: true,
	},
	&post{
		Id:        1,
		Title:     "Title Two",
		Content:   "Content Two",
		BigNumber: 123125241231231123,
		BoolValue: false,
	},
}

const (
	createPostTable = `DROP TABLE IF EXISTS post;
	       CREATE TABLE post (
	         ID bigint,
	         TITLE varchar(40),
	         CONTENT varchar(255),
			 BIGNUMBER bigint,
			 BOOLVALUE bit,
	         CONSTRAINT pk_post PRIMARY KEY(ID)
	       );
	       INSERT INTO post(ID,TITLE,CONTENT,BIGNUMBER,BOOLVALUE) VALUES
	                       (0,NULL,'Content One',1231241241231231123,1),
	                       (1,'Title Two','Content Two',123125241231231123,0);`
	selectPostTable = `SELECT id,title,content,bignumber,boolvalue FROM post;`
)

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
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad connection string",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      "bad connection string",
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad port connection string",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionStringBadPort,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no connection string",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad max idle connections",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "-1",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad max open connections",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "-1",
					"connection_max_lifetime_seconds": "",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad connection max lifetime seconds",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "-1",
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

func TestClient_Query_Exec_Transaction(t *testing.T) {
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
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte(selectPostTable)),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetData(allPosts.marshal()),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "empty exec request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec"),

			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid exec request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte("bad statement")),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "valid exec empty query request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
		{
			name: "valid exec bad query request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("some bad query")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: nil,
			wantExecErr:       false,
			wantQueryErr:      true,
		},
		{
			name: "valid exec valid query - no results",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte("SELECT id,title,content FROM post where id=100")),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "valid exec query request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "exec").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetData([]byte(selectPostTable)),
			wantExecResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantQueryResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetData(allPosts.marshal()),
			wantExecErr:  false,
			wantQueryErr: false,
		},
		{
			name: "empty transaction request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "transaction"),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "invalid transaction request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetData([]byte("bad statement")),
			queryRequest:      nil,
			wantExecResponse:  nil,
			wantQueryResponse: nil,
			wantExecErr:       true,
			wantQueryErr:      false,
		},
		{
			name: "valid transaction empty query request",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			execRequest: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetData([]byte(createPostTable)),
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query"),

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

			if tt.wantQueryResponse != nil {
				wantPosts := unmarshal(tt.wantQueryResponse.Data)
				var gotPosts *posts
				if gotGetResponse != nil {
					gotPosts = unmarshal(gotGetResponse.Data)
				}
				require.EqualValues(t, wantPosts, gotPosts)
			} else {
				require.EqualValues(t, tt.wantQueryResponse, gotGetResponse)
			}

		})
	}
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
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "read_uncommitted").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "valid request - 2",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "read_committed").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "valid request - 3",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "repeatable_read").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "valid request - 3",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "serializable").
				SetData([]byte(createPostTable)),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method"),
			wantErr: true,
		},
		{
			name: "invalid request - bad isolation level",
			cfg: config.Spec{
				Name: "target-azure-stores-azuresql",
				Kind: "target.azure.stores.azuresql",
				Properties: map[string]string{
					"connection":                      dat.connectionString,
					"max_idle_connections":            "",
					"max_open_connections":            "",
					"connection_max_lifetime_seconds": "",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "transaction").
				SetMetadataKeyValue("isolation_level", "bad_level").
				SetData([]byte(createPostTable)),
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
