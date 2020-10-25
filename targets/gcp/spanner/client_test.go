package spanner

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
	db        string
	query     string
	tableName string
	cred      string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/querySpanner.txt")
	if err != nil {
		return nil, err
	}
	t.query = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/dbSpanner.txt")
	if err != nil {
		return nil, err
	}
	t.db = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/tableName.txt")
	if err != nil {
		return nil, err
	}
	t.tableName = string(dat)
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
				Name: "target.gcp.spanner",
				Kind: "target.gcp.spanner",
				Properties: map[string]string{
					"db":          dat.db,
					"credentials": dat.cred,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - missing db",
			cfg: config.Spec{
				Name: "target.gcp.spanner",
				Kind: "target.gcp.spanner",
				Properties: map[string]string{
					"credentials": dat.cred,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing credentials",
			cfg: config.Spec{
				Name: "target.gcp.spanner",
				Kind: "target.gcp.spanner",
				Properties: map[string]string{
					"db": dat.db,
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

	cfg := config.Spec{
		Name: "target.gcp.spanner",
		Kind: "target.gcp.spanner",
		Properties: map[string]string{
			"db":          dat.db,
			"credentials": dat.cred,
		},
	}
	tests := []struct {
		name         string
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid query",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("query", dat.query),
			wantErr: false,
		}, {
			name: "invalid query - missing query",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "query"),
			wantErr: true,
		}, {
			name: "invalid query- missing method",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("query", dat.query),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClient_Read(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	columnNames := []string{"id", "name"}
	b, err := json.Marshal(columnNames)
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target.gcp.spanner",
		Kind: "target.gcp.spanner",
		Properties: map[string]string{
			"db":          dat.db,
			"credentials": dat.cred,
		},
	}
	tests := []struct {
		name         string
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid read",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "read").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData(b),
			wantErr: false,
		},
		{
			name: "invalid read - missing data",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "read").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		},
		{
			name: "invalid read - missing table_name",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "read"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClient_Insert(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	values := make([]interface{}, 0)
	values = append(values, 17, "name1")
	firstInsUpd := InsertOrUpdate{dat.tableName, []string{"id", "name"}, values, []string{"INT64", "STRING"}}
	values = make([]interface{}, 0)
	values = append(values, 18, "name2")
	scnInsUpd := InsertOrUpdate{dat.tableName, []string{"id", "name"}, values, []string{"INT64", "STRING"}}
	var inputs []InsertOrUpdate
	cfg := config.Spec{
		Name: "target.gcp.spanner",
		Kind: "target.gcp.spanner",
		Properties: map[string]string{
			"db":          dat.db,
			"credentials": dat.cred,
		},
	}
	inputs = append(inputs, firstInsUpd, scnInsUpd)

	bSchema, err := json.Marshal(inputs)
	require.NoError(t, err)

	tests := []struct {
		name         string
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid insert",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetData(bSchema),
			wantErr: false,
		}, {
			name: "invalid insert - missing data",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClient_Update(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	values := make([]interface{}, 0)
	values = append(values, 17, "name3")
	firstInsUpd := InsertOrUpdate{dat.tableName, []string{"id", "name"}, values, []string{"INT64", "STRING"}}
	values = make([]interface{}, 0)
	values = append(values, 18, "name4")
	scnInsUpd := InsertOrUpdate{dat.tableName, []string{"id", "name"}, values, []string{"INT64", "STRING"}}
	var inputs []InsertOrUpdate
	cfg := config.Spec{
		Name: "target.gcp.spanner",
		Kind: "target.gcp.spanner",
		Properties: map[string]string{
			"db":          dat.db,
			"credentials": dat.cred,
		},
	}
	inputs = append(inputs, firstInsUpd, scnInsUpd)

	bSchema, err := json.Marshal(inputs)
	require.NoError(t, err)

	tests := []struct {
		name         string
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid update",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "update").
				SetData(bSchema),
			wantErr: false,
		}, {
			name: "invalid update - missing data",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "update"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClient_UpdateDatabaseDdl(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	var statements []string
	statements = append(statements, "mystatment")
	cfg := config.Spec{
		Name: "target.gcp.spanner",
		Kind: "target.gcp.spanner",
		Properties: map[string]string{
			"db":          dat.db,
			"credentials": dat.cred,
		},
	}

	bSchema, err := json.Marshal(statements)
	require.NoError(t, err)

	tests := []struct {
		name         string
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid UpdateDatabaseDdl",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "update_database_ddl").
				SetData(bSchema),
			wantErr: false,
		}, {
			name: "invalid InsertOrUpdate - missing data",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "update_database_ddl"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClient_InsertOrUpdate(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	values := make([]interface{}, 0)
	values = append(values, 19, "name5")
	firstInsUpd := InsertOrUpdate{dat.tableName, []string{"id", "name"}, values, []string{"INT64", "STRING"}}
	values = make([]interface{}, 0)
	values = append(values, 20, "name6")
	scnInsUpd := InsertOrUpdate{dat.tableName, []string{"id", "name"}, values, []string{"INT64", "STRING"}}
	var inputs []InsertOrUpdate
	cfg := config.Spec{
		Name: "target.gcp.spanner",
		Kind: "target.gcp.spanner",
		Properties: map[string]string{
			"db":          dat.db,
			"credentials": dat.cred,
		},
	}
	inputs = append(inputs, firstInsUpd, scnInsUpd)

	bSchema, err := json.Marshal(inputs)
	require.NoError(t, err)

	tests := []struct {
		name         string
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid InsertOrUpdate",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert_or_update").
				SetData(bSchema),
			wantErr: false,
		}, {
			name: "invalid InsertOrUpdate - missing data",
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert_or_update"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
