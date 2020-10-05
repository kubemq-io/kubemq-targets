package athena

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"testing"
	"time"
)

type testStructure struct {
	awsKey       string
	awsSecretKey string
	region       string
	token        string

	query          string
	catalog        string
	db             string
	outputLocation string
	executionID    string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/aws/awsKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/awsSecretKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsSecretKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/region.txt")
	if err != nil {
		return nil, err
	}
	t.region = fmt.Sprintf("%s", dat)
	t.token = ""

	dat, err = ioutil.ReadFile("./../../../credentials/aws/athena/query.txt")
	if err != nil {
		return nil, err
	}
	t.query = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/athena/catalog.txt")
	if err != nil {
		return nil, err
	}
	t.catalog = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/athena/db.txt")
	if err != nil {
		return nil, err
	}
	t.db = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/athena/outputLocation.txt")
	if err != nil {
		return nil, err
	}
	t.outputLocation = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/athena/queryExecutionID.txt")
	if err != nil {
		return nil, err
	}
	t.executionID = string(dat)
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
				Name: "target-aws-athena",
				Kind: "target.aws.athena",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing aws key",
			cfg: config.Spec{
				Name: "target-aws-athena",
				Kind: "target.aws.athena",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing region",
			cfg: config.Spec{
				Name: "target-aws-athena",
				Kind: "target.aws.athena",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "target-aws-athena",
				Kind: "target.aws.athena",
				Properties: map[string]string{
					"aws_key": dat.awsKey,
					"region":  dat.region,
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
			require.NoError(t, err)

		})
	}
}

func TestClient_ListCatalogs(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-athena",
		Kind: "target.aws.athena",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list_data_catalogs",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_data_catalogs"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_ListDatabases(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-athena",
		Kind: "target.aws.athena",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list_databases",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_databases").
				SetMetadataKeyValue("catalog", dat.catalog),
			wantErr: false,
		},
		{
			name: "invalid list_databases - missing catalog",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_databases"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_Query(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-athena",
		Kind: "target.aws.athena",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid query",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("db", dat.db).
				SetMetadataKeyValue("output_location", dat.outputLocation).
				SetMetadataKeyValue("catalog", dat.catalog).
				SetMetadataKeyValue("query", dat.query),
			wantErr: false,
		},
		{
			name: "invalid query - missing query",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("db", dat.db).
				SetMetadataKeyValue("output_location", dat.outputLocation).
				SetMetadataKeyValue("catalog", dat.catalog),
			wantErr: true,
		},
		{
			name: "invalid query - missing catalog",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("db", dat.db).
				SetMetadataKeyValue("output_location", dat.outputLocation).
				SetMetadataKeyValue("query", dat.query),
			wantErr: true,
		},
		{
			name: "invalid query - missing db",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("catalog", dat.catalog).
				SetMetadataKeyValue("output_location", dat.outputLocation).
				SetMetadataKeyValue("query", dat.query),
			wantErr: true,
		},
		{
			name: "invalid query - missing output_location",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "query").
				SetMetadataKeyValue("db", dat.db).
				SetMetadataKeyValue("catalog", dat.catalog).
				SetMetadataKeyValue("query", dat.query),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

func TestClient_GetQueryResult(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-athena",
		Kind: "target.aws.athena",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()

	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid get query",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_query_result").
				SetMetadataKeyValue("execution_id", dat.executionID),
			wantErr: false,
		},
		{
			name: "invalid get query - missing execution_id",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_query_result"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
