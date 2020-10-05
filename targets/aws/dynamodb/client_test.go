package dynamodb

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

	tableName string
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

	dat, err = ioutil.ReadFile("./../../../credentials/aws/dynamodb/tableName.txt")
	if err != nil {
		return nil, err
	}
	t.tableName = string(dat)
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
				Name: "target-aws-dynamodb",
				Kind: "target.aws.dynamodb",
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
				Name: "target-aws-dynamodb",
				Kind: "target.aws.dynamodb",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing region",
			cfg: config.Spec{
				Name: "target-aws-dynamodb",
				Kind: "target.aws.dynamodb",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "target-aws-dynamodb",
				Kind: "target.aws.dynamodb",
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

func TestClient_List(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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
				SetMetadataKeyValue("method", "list_tables"),
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

func TestClient_CreateTable(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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
	input := fmt.Sprintf(`{
					"AttributeDefinitions": [
						{
							"AttributeName": "Year",
							"AttributeType": "N"
						},
						{
							"AttributeName": "Title",
							"AttributeType": "S"
						}
					],
					"BillingMode": null,
					"GlobalSecondaryIndexes": null,
					"KeySchema": [
						{
							"AttributeName": "Year",
							"KeyType": "HASH"
						},
						{
							"AttributeName": "Title",
							"KeyType": "RANGE"
						}
					],
					"LocalSecondaryIndexes": null,
					"ProvisionedThroughput": {
						"ReadCapacityUnits": 10,
						"WriteCapacityUnits": 10
					},
					"SSESpecification": null,
					"StreamSpecification": null,
					"TableName": "%s",
					"Tags": null
				}`, dat.tableName)

	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetData([]byte(input)),
			wantErr: false,
		},
		{
			name: "invalid create - table already exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_table").
				SetData([]byte(input)),
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

func TestClient_InsertItem(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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

	input := `{
		"Plot": {
			"S": "some plot"
		},
		"Rating": {
			"N": "10.1"
		},
		"Title": {
			"S": "KubeMQ test Movie"
		},
		"Year": {
			"N": "2020"
		}
	}`

	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid put item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "insert_item").
				SetMetadataKeyValue("table_name", dat.tableName).
				SetData([]byte(input)),
			wantErr: false,
		},
		{
			name: "invalid put item - missing table_name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "insert_item").
				SetData([]byte(input)),
			wantErr: true,
		},
		{
			name: "invalid put item - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "insert_item").
				SetMetadataKeyValue("table_name", dat.tableName),
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

func TestClient_GetItem(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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

	input := fmt.Sprintf(`
		{
			"AttributesToGet": null,
			"ConsistentRead": null,
			"ExpressionAttributeNames": null,
			"Key": {
				"Title": {
					"S": "KubeMQ test Movie"
				},
				"Year": {
					"N": "2020"
				}
			},
			"ProjectionExpression": null,
			"ReturnConsumedCapacity": null,
			"TableName": "%s"
		}`,
		dat.tableName)

	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid put get item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_item").
				SetData([]byte(input)),
			wantErr: false,
		},
		{
			name: "invalid get item - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_item"),
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

func TestClient_UpdateItem(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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

	input := fmt.Sprintf(`
		{
			"ExpressionAttributeValues": {
				":r": {
					"N": "0.9"
				}
			},
			"Key": {
				"Title": {
					"S": "KubeMQ test Movie"
				},
				"Year": {
					"N": "2020"
				}
			},
			"ReturnValues": "UPDATED_NEW",
			"TableName": "%s",
			"UpdateExpression": "set Rating = :r"
		}`,
		dat.tableName)

	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid update item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update_item").
				SetData([]byte(input)),
			wantErr: false,
		},
		{
			name: "invalid update item - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "update_item"),
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

func TestClient_DeleteItem(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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
	input := fmt.Sprintf(
		`{
					"Key": {
						"Title": {
							"S": "KubeMQ test Movie"
						},
						"Year": {
							"N": "2020"
						}
					},
					"TableName": "%s"
				}`,
		dat.tableName)

	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete_item item",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_item").
				SetData([]byte(input)),
			wantErr: false,
		},
		{
			name: "invalid delete_item - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_item"),
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

func TestClient_DeleteTable(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-dynamodb",
		Kind: "target.aws.dynamodb",
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
			name: "valid delete_table ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_table").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: false,
		},
		{
			name: "invalid delete_table - table does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_table").
				SetMetadataKeyValue("table_name", dat.tableName),
			wantErr: true,
		},
		{
			name: "invalid delete_table - missing table name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_table"),
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
