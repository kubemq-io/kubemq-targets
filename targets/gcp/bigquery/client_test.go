package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type testStructure struct {
	projectID     string
	tableName     string
	query         string
	emptyTable    string
	emptyTableQry string
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
	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			wantErr: false,
		},
		{
			name: "google-big-query-target",
			cfg: config.Metadata{
				Name: "google-big_table-target",
				Kind: "",
				Properties: map[string]string{
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
				_=c.CloseClient()
			}()
			require.NoError(t, err)
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Query(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name            string
		cfg             config.Metadata
		queryRequest    *types.Request
		wantErr         bool
		wantEmptyData   bool
		wantReaderError bool
	}{
		{
			name: "valid query",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("query", dat.query),
			wantErr:         false,
			wantEmptyData:   false,
			wantReaderError: false,
		}, {
			name: "invalid query - missing query",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query"),
			wantErr:         true,
			wantEmptyData:   false,
			wantReaderError: false,
		}, {
			name: "valid query- empty table",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("query", dat.emptyTableQry),
			wantErr:         false,
			wantEmptyData:   true,
			wantReaderError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				_=c.CloseClient()
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			if tt.wantEmptyData {
				require.EqualValues(t, gotSetResponse.Metadata["error"], "true")
				require.EqualValues(t, gotSetResponse.Metadata["message"], "no rows found for this query")
				t.Logf("init() error = %v, wantErr %v", gotSetResponse.Metadata["message"], tt.wantErr)
				return
			}
			if tt.wantReaderError {
				require.EqualValues(t, gotSetResponse.Metadata["error"], "true")
				t.Logf("init() error = %v, wantErr %v", gotSetResponse.Metadata["message"], tt.wantErr)
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
		name           string
		cfg            config.Metadata
		queryRequest   *types.Request
		wantErr        bool
		wantWriteError bool
	}{
		{
			name: "valid create table",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("dataset_id", "my_data_set").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bSchema),
			wantErr:        false,
			wantWriteError: false,
		}, {
			name: "invalid create_table - missing tableName",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("dataset_id", "my_data_set").
				SetData(bSchema),
			wantErr:        true,
			wantWriteError: false,
		}, {
			name: "invalid create_table - table already exists",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("dataset_id", "my_data_set").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bSchema),
			wantErr:        false,
			wantWriteError: true,
		}, {
			name: "invalid create_table - missing dataset_id",
			cfg: config.Metadata{
				Name: "google-big-query-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bSchema),
			wantErr:        true,
			wantWriteError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer func() {
				_=c.CloseClient()
			}()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.queryRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			if tt.wantWriteError {
				require.EqualValues(t, gotSetResponse.Metadata["error"], "true")
				t.Logf("init() error = %v, wantErr %v", gotSetResponse.Metadata["message"], tt.wantWriteError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}
