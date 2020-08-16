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

	query   string
	catalog string
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
				Name: "aws-athena",
				Kind: "aws.athena",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing key",
			cfg: config.Spec{
				Name: "aws-athena",
				Kind: "aws.athena",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing region",
			cfg: config.Spec{
				Name: "aws-athena",
				Kind: "aws.athena",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "aws-athena",
				Kind: "aws.athena",
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
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_ListCatalogs(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-athena",
		Kind: "aws.athena",
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
		Name: "aws-athena",
		Kind: "aws.athena",
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
		Name: "aws-athena",
		Kind: "aws.athena",
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
				SetMetadataKeyValue("query", dat.query),
			wantErr: false,
		},
		{
			name: "invalid query - missing query",
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
