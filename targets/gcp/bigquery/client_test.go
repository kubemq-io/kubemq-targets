package bigquery

import (
	"cloud.google.com/go/bigquery"
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
	projectID     string
	tableName     string
	query         string
	dataSetID     string
	emptyTable    string
	emptyTableQry string
	cred          string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	if err != nil {
		return nil, err
	}
	t.projectID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/query.txt")
	if err != nil {
		return nil, err
	}
	t.query = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/tableName.txt")
	if err != nil {
		return nil, err
	}
	t.tableName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/emptyTable.txt")
	if err != nil {
		return nil, err
	}
	t.emptyTable = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/emptyTableQry.txt")
	if err != nil {
		return nil, err
	}
	t.emptyTableQry = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/dataSetID.txt")
	if err != nil {
		return nil, err
	}
	t.dataSetID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	if err != nil {
		return nil, err
	}
	t.cred = fmt.Sprintf("%s", dat)
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
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing credentials",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing project_id",
			cfg: config.Spec{
				Name: "target-gcp-bigtable",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"credentials": dat.cred,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c := New()

			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			defer func() {
				_ = c.Stop()
			}()
			require.NoError(t, err)

		})
	}
}

func TestClient_Query(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name         string
		cfg          config.Spec
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid query",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("query", dat.query),
			wantErr: false,
		}, {
			name: "invalid query - missing query",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query"),
			wantErr: true,
		}, {
			name: "valid query- empty table",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("query", dat.emptyTableQry),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Create_Table(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	mySchema := bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema:         mySchema,
		ExpirationTime: time.Now().AddDate(2, 1, 0), // Table will deleted in 2 years and 1 month.
	}

	bSchema, err := json.Marshal(metaData)
	require.NoError(t, err)

	tests := []struct {
		name         string
		cfg          config.Spec
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid create table",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("dataset_id", "my_data_set").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bSchema),
			wantErr: false,
		}, {
			name: "invalid create_table - missing tableName",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("dataset_id", "my_data_set").
				SetData(bSchema),
			wantErr: true,
		}, {
			name: "invalid create_table - table already exists",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("dataset_id", "my_data_set").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bSchema),
			wantErr: true,
		}, {
			name: "invalid create_table - missing dataset_id",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bSchema),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Get_Data_Sets(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name         string
		cfg          config.Spec
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid get Data-Sets",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_data_sets"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.EqualValues(t, gotSetResponse.Metadata["result"], "ok")
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Get_Table_Info(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name         string
		cfg          config.Spec
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid get table-info",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_table_info").
				SetMetadataKeyValue("dataset_id", dat.dataSetID).
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: false,
		}, {
			name: "invalid get table-info - missing dataset_id",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_table_info").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		}, {
			name: "invalid get table-info - missing table_name",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_table_info").
				SetMetadataKeyValue("dataset_id", dat.dataSetID),
			wantErr: true,
		}, {
			name: "valid get table-info - missing NotExistingTable",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_table_info").
				SetMetadataKeyValue("dataset_id", dat.dataSetID).
				SetMetadataKeyValue("table_name", "NotExistingTable"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.EqualValues(t, gotSetResponse.Metadata["result"], "ok")
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Insert_To_Table(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	var rows []map[string]bigquery.Value
	firstRow := make(map[string]bigquery.Value)
	firstRow["name"] = "myName4"
	firstRow["age"] = 25
	rows = append(rows, firstRow)
	secondRow := make(map[string]bigquery.Value)
	secondRow["name"] = "myName5"
	secondRow["age"] = 28
	rows = append(rows, secondRow)
	bRows, err := json.Marshal(&rows)
	require.NoError(t, err)
	tests := []struct {
		name         string
		cfg          config.Spec
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid insert to table",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetMetadataKeyValue("dataset_id", dat.dataSetID).
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bRows),
			wantErr: false,
		}, {
			name: "invalid insert to table - missing table_name",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetMetadataKeyValue("dataset_id", dat.dataSetID).
				SetData(bRows),
			wantErr: true,
		}, {
			name: "invalid insert to table - missing dataset_id",
			cfg: config.Spec{
				Name: "target-gcp-bigquery",
				Kind: "target.gcp.bigquery",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bRows),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, gotSetResponse.Metadata["result"], "ok")
			require.NotNil(t, gotSetResponse)
		})
	}
}
