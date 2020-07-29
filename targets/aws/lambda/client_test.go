package lambda

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
	
	zipFileName  string
	functionName string
	handlerName  string
	role         string
	runtime      string
	description  string
	
	lambdaExp []byte
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
	
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/zipFileName.txt")
	if err != nil {
		return nil, err
	}
	t.zipFileName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/functionName.txt")
	if err != nil {
		return nil, err
	}
	t.functionName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/handlerName.txt")
	if err != nil {
		return nil, err
	}
	t.handlerName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/role.txt")
	if err != nil {
		return nil, err
	}
	t.role = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/runtime.txt")
	if err != nil {
		return nil, err
	}
	t.runtime = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/lambda/description.txt")
	if err != nil {
		return nil, err
	}
	t.description = fmt.Sprintf("%s", dat)
	contents, err := ioutil.ReadFile("./../../../credentials/aws/lambda/lambdaCode.zip")
	if err != nil {
		return nil, err
	}
	t.lambdaExp = contents
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
				Name: "aws-lambda",
				Kind: "aws.lambda",
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
				Name: "aws-lambda",
				Kind: "aws.lambda",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing region",
			cfg: config.Spec{
				Name: "aws-lambda",
				Kind: "aws.lambda",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "aws-lambda",
				Kind: "aws.lambda",
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

func TestClient_List(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
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
			name: "valid list",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list"),
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

func TestClient_Create(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
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
			name: "valid create",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("memorySize", "256").
				SetMetadataKeyValue("timeout", "15").
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime).
				SetData(dat.lambdaExp),
			wantErr: false,
		},
		{
			name: "invalid create- already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime).
				SetData(dat.lambdaExp),
			wantErr: true,
		},
		{
			name: "invalid create- missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}


func TestClient_Delete(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-lambda",
		Kind: "aws.lambda",
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
			name: "valid create",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("memorySize", "256").
				SetMetadataKeyValue("timeout", "15").
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime).
				SetData(dat.lambdaExp),
			wantErr: false,
		},
		{
			name: "invalid create- already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime).
				SetData(dat.lambdaExp),
			wantErr: true,
		},
		{
			name: "invalid create- missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create").
				SetMetadataKeyValue("zip_file_name", dat.zipFileName).
				SetMetadataKeyValue("description", dat.description).
				SetMetadataKeyValue("handler_name", dat.handlerName).
				SetMetadataKeyValue("role", dat.role).
				SetMetadataKeyValue("function_name", dat.functionName).
				SetMetadataKeyValue("runtime", dat.runtime),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
