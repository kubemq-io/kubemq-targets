package spanner

import (
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
	db     string
	query         string
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
				Name: "google-spanner-target",
				Kind: "",
				Properties: map[string]string{
					"db": dat.db,
				},
			},
			wantErr: false,
		},
		{
			name: "init",
			cfg: config.Metadata{
				Name: "google-spanner-target",
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
				_ = c.CloseClient()
			}()
			require.NoError(t, err)
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Query(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	cfg := config.Metadata{
		Name: "google-spanner-target",
		Kind: "",
		Properties: map[string]string{
			"db": dat.db,
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
		err = c.CloseClient()
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
	values = append(values, "first", 18)
	firstInsUpd := InsertOrUpdate{"exists_table",[]string{"name","age"},values}
	values = make([]interface{}, 0)
	values = append(values, "first", 18)
	scnInsUpd := InsertOrUpdate{"exists_table",[]string{"name","age"},values}
	var inputs []InsertOrUpdate

	inputs = append(inputs,firstInsUpd,scnInsUpd)

	bSchema, err := json.Marshal(inputs)
	require.NoError(t, err)

	tests := []struct {
		name         string
		cfg          config.Metadata
		queryRequest *types.Request
		wantErr      bool
	}{
		{
			name: "valid insert",
			cfg: config.Metadata{
				Name: "google-spanner-target",
				Kind: "",
				Properties: map[string]string{
					"db": dat.db,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert").
				SetData(bSchema),
			wantErr: false,
		}, {
			name: "invalid valid insert - missing data",
			cfg: config.Metadata{
				Name: "google-spanner-target",
				Kind: "",
				Properties: map[string]string{
					"db": dat.db,
				},
			},
			queryRequest: types.NewRequest().
				SetMetadataKeyValue("method", "insert"),
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
				_ = c.CloseClient()
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
