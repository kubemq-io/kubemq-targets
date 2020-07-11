package bigtable

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

type testStructure struct {
	projectID    string
	instance     string
	tableName    string
	tempTable    string
	columnFamily string
	rowKeyPrefix string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/instance.txt")
	if err != nil {
		return nil, err
	}
	t.instance = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/projectID.txt")
	if err != nil {
		return nil, err
	}
	t.projectID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/tableName.txt")
	if err != nil {
		return nil, err
	}
	t.tableName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/tempTable.txt")
	if err != nil {
		return nil, err
	}
	t.tempTable = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/columnFamily.txt")
	if err != nil {
		return nil, err
	}
	t.columnFamily = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/rowKeyPrefix.txt")
	if err != nil {
		return nil, err
	}
	t.rowKeyPrefix = string(dat)

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
				Name: "google-big_table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			wantErr: false,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Metadata{
				Name: "google-big_table-target",
				Kind: "",
				Properties: map[string]string{
					"instance": dat.instance,
				},
			},
			wantErr: true,
		},
		{
			name: "init-missing-instance",
			cfg: config.Metadata{
				Name: "google-big_table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
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
			defer c.CloseAdminClient()
			require.NoError(t, err)
			require.EqualValues(t, tt.cfg.Name, c.Name())
			err = c.CloseAdminClient()
			require.NoError(t, err)
		})
	}
}

func TestClient_Create_Column_Family(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Metadata{
		Name: "google-big-table-target",
		Kind: "",
		Properties: map[string]string{
			"project_id": dat.projectID,
			"instance":   dat.instance,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name              string
		cfg               config.Metadata
		request           *types.Request
		wantResponseError bool
		wantAttributesErr bool
	}{
		{
			name: "valid create create-column-family",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_column_family").
				SetMetadataKeyValue("column_family", dat.columnFamily).
				SetMetadataKeyValue("table_name", dat.tableName),
			wantResponseError: false,
			wantAttributesErr: false,
		}, {
			name: "invalid create-column-family -invalid type",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"instance":   dat.instance,
					"project_id": dat.projectID,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_column_family2").
				SetMetadataKeyValue("column_family", dat.columnFamily).
				SetMetadataKeyValue("table_name", dat.tableName),
			wantResponseError: false,
			wantAttributesErr: true,
		}, {
			name: "invalid create-column-family- already exists",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_column_family").
				SetMetadataKeyValue("column_family", dat.columnFamily).
				SetMetadataKeyValue("table_name", dat.tableName),
			wantAttributesErr: false,
			wantResponseError: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantAttributesErr {
				t.Logf("init() error = %v, wantAttributesErr %v", err, tt.wantAttributesErr)
				require.Error(t, err)
				return
			}
			if tt.wantResponseError {
				require.Equal(t, gotSetResponse.Metadata["error"], "true")
				t.Logf("init() error = %v, wantResponseError %v", err, gotSetResponse.Metadata["message"])
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
	err = c.CloseAdminClient()
	require.NoError(t, err)
}

func TestClient_Create_Delete_Table(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Metadata{
		Name: "google-big-table-target",
		Kind: "",
		Properties: map[string]string{
			"project_id": dat.projectID,
			"instance":   dat.instance,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name              string
		cfg               config.Metadata
		request           *types.Request
		wantResponseError bool
		wantAttributesErr bool
	}{
		{
			name: "valid create table",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantResponseError: false,
			wantAttributesErr: false,
		}, {
			name: "invalid create table -invalid type",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"instance":   dat.instance,
					"project_id": dat.projectID,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_table2").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantResponseError: false,
			wantAttributesErr: true,
		}, {
			name: "valid create table- already exists",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantAttributesErr: false,
			wantResponseError: false,
		}, {
			name: "valid delete table",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_table").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantAttributesErr: false,
			wantResponseError: false,
		}, {
			name: "invalid delete table - table does not exists",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_table").
				SetMetadataKeyValue("table_name", dat.tempTable),
			wantAttributesErr: false,
			wantResponseError: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantAttributesErr {
				t.Logf("init() error = %v, wantAttributesErr %v", err, tt.wantAttributesErr)
				require.Error(t, err)
				return
			}
			if tt.wantResponseError {
				require.Equal(t, gotSetResponse.Metadata["error"], "true")
				t.Logf("init() error = %v, wantResponseError %v", err, gotSetResponse.Metadata["message"])
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
	err = c.CloseAdminClient()
	require.NoError(t, err)
}

func TestClient_write(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	singleRow := map[string]interface{}{"set_row_key": fmt.Sprintf("%d", rand.Int()), "id": 1, "name": "test1"}
	var rows []map[string]interface{}
	rowOne := map[string]interface{}{"set_row_key": fmt.Sprintf("%d", rand.Int()), "id": 2, "name": "test2"}
	rowTwo := map[string]interface{}{"set_row_key": fmt.Sprintf("%d", rand.Int()), "id": 3, "name": "test3"}
	rows = append(rows, rowOne)
	rows = append(rows, rowTwo)

	singleB, err := json.Marshal(singleRow)
	require.NoError(t, err)
	multiB, err := json.Marshal(rows)
	require.NoError(t, err)
	tests := []struct {
		name              string
		cfg               config.Metadata
		writeRequest      *types.Request
		wantWriteResponse *types.Response
		wantWriteErr      bool
	}{
		{
			name: "valid single write",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "write").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("column_family", dat.columnFamily).
				SetData(singleB),
			wantWriteResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantWriteErr: false,
		},
		{
			name: "valid single write",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "write_batch").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("column_family", dat.columnFamily).
				SetData(multiB),
			wantWriteResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantWriteErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer c.CloseAdminClient()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.writeRequest)
			if tt.wantWriteErr {
				t.Logf("init() error = %v, wantWriteErr %v", err, tt.wantWriteErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantWriteResponse, gotSetResponse)
		})
	}
}

func TestClient_Delete_Rows(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name          string
		cfg           config.Metadata
		deleteRequest *types.Request
		wantErr       bool
		wantDeleteErr bool
	}{
		{
			name: "valid delete rows",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_row").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("row_key_prefix", dat.rowKeyPrefix),
			wantErr: false,
		},
		{
			name: "invalid delete rows - missing row_key_prefix",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_row").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		},
		{
			name: "invalid delete rows - missing table name",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_row").
				SetMetadataKeyValue("row_key_prefix", dat.rowKeyPrefix),
			wantErr: true,
		},
		{
			name: "invalid delete rows - table doesnt exists",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			deleteRequest: types.NewRequest().
				SetMetadataKeyValue("method", "delete_row").
				SetMetadataKeyValue("table_name", "fake_table").
				SetMetadataKeyValue("row_key_prefix", dat.rowKeyPrefix),
			wantDeleteErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer c.CloseAdminClient()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.deleteRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			if tt.wantDeleteErr {
				require.Equal(t, gotSetResponse.Metadata["error"], "true")
				t.Logf("init() error = %v, wantErr %v", err, gotSetResponse.Metadata["message"])
				return
			}
			require.NotEqual(t, gotSetResponse.Metadata["error"], "true")
			t.Logf("init() error = %v, response %v", err, fmt.Sprintf("%s", gotSetResponse.Data))
		})
	}
}

func TestClient_Read_Rows(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	keys := []string{
		dat.rowKeyPrefix,
	}
	bKeys, err := json.Marshal(keys)
	require.NoError(t, err)
	tests := []struct {
		name         string
		cfg          config.Metadata
		writeRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid read all rows",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_all_rows").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: false,
		},
		{
			name: "valid read all rows by keys",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_all_rows").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bKeys),
			wantErr: false,
		}, {
			name: "valid read all rows - column_filter",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_all_rows_by_column").
				SetMetadataKeyValue("column_name", "id").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: false,
		},
		{
			name: "valid read all rows by keys - column_filter",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_all_rows_by_column").
				SetMetadataKeyValue("column_name", "id").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: false,
		}, {
			name: "invalid read all rows - column_filter - missing column_name",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_all_rows_by_column").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(bKeys),
			wantErr: true,
		}, {
			name: "valid read row",
			cfg: config.Metadata{
				Name: "google-big-table-target",
				Kind: "",
				Properties: map[string]string{
					"project_id": dat.projectID,
					"instance":   dat.instance,
				},
			},
			writeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get_row").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetMetadataKeyValue("row_key_prefix", "my_id"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			defer c.CloseAdminClient()
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.writeRequest)
			if tt.wantErr {
				t.Logf("init() error = %v, wantWriteErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			t.Logf("init() error = %v, response %v", err, fmt.Sprintf("%s", gotSetResponse.Data))
		})
	}
}
