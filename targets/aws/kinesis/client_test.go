package kinesis

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
	awsKey       string
	awsSecretKey string
	region       string
	token        string

	shardCount    string
	streamName    string
	partitionKey  string
	shardPosition string
	shardID       string
	limit         string
	streamARN     string
	shardIterator string
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

	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/shardCount.txt")
	if err != nil {
		return nil, err
	}
	t.shardCount = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/streamName.txt")
	if err != nil {
		return nil, err
	}
	t.streamName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/partitionKey.txt")
	if err != nil {
		return nil, err
	}
	t.partitionKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/shardPosition.txt")
	if err != nil {
		return nil, err
	}
	t.shardPosition = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/shardID.txt")
	if err != nil {
		return nil, err
	}
	t.shardID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/limit.txt")
	if err != nil {
		return nil, err
	}
	t.limit = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/streamARN.txt")
	if err != nil {
		return nil, err
	}
	t.streamARN = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/shardIterator.txt")
	if err != nil {
		return nil, err
	}
	t.shardIterator = string(dat)
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
				Name: "aws-kinesis",
				Kind: "aws.kinesis",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing aws_key",
			cfg: config.Spec{
				Name: "aws-kinesis",
				Kind: "aws.kinesis",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing region",
			cfg: config.Spec{
				Name: "aws-kinesis",
				Kind: "aws.kinesis",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing aws_secret_key",
			cfg: config.Spec{
				Name: "aws-kinesis",
				Kind: "aws.kinesis",
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

func TestClient_ListStreams(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
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
			name: "valid list_streams",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_streams"),
			wantErr: false,
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

func TestClient_ListStreamConsumers(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
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
			name: "valid list_streams",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_stream_consumers").
				SetMetadataKeyValue("stream_arn", dat.streamARN),
			wantErr: false,
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

func TestClient_CreateStream(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
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
			name: "valid create_stream",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_stream").
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetMetadataKeyValue("shard_count", dat.shardCount),
			wantErr: false,
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

func TestClient_ListShards(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
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
			name: "valid list_shards",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_shards").
				SetMetadataKeyValue("stream_name", dat.streamName),
			wantErr: false,
		}, {
			name: "invalid list_shards - missing stream name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_shards"),
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

func TestClient_GetShardIterator(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
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
			name: "valid get_shard_iterator",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_shard_iterator").
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetMetadataKeyValue("shard_iterator_type", "LATEST").
				SetMetadataKeyValue("shard_id", dat.shardID),
			wantErr: false,
		}, {
			name: "invalid get_shard_iterator - missing stream",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_shard_iterator").
				SetMetadataKeyValue("shard_iterator_type", "LATEST").
				SetMetadataKeyValue("shard_id", dat.shardID),
			wantErr: true,
		}, {
			name: "invalid get_shard_iterator - missing type",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_shard_iterator").
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetMetadataKeyValue("shard_id", dat.shardID),
			wantErr: true,
		}, {
			name: "invalid get_shard_iterator - missing shard_id",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_shard_iterator").
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetMetadataKeyValue("shard_iterator_type", "LATEST"),
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
			require.NotNil(t, got.Metadata["shard_iterator"])
			t.Logf("got iterator: %s", got.Metadata["shard_iterator"])
		})
	}
}

func TestClient_PutRecord(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	data := `{"my_result":"ok"}`
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid put_record",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_record").
				SetMetadataKeyValue("partition_key", dat.partitionKey).
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetData([]byte(data)),
			wantErr: false,
		}, {
			name: "invalid put_record - missing partition_key",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_record").
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetData([]byte(data)),
			wantErr: true,
		}, {
			name: "invalid put_record  - missing stream_name ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_record").
				SetMetadataKeyValue("partition_key", dat.partitionKey).
				SetData([]byte(data)),
			wantErr: true,
		}, {
			name: "invalid put_record - missing data",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_record").
				SetMetadataKeyValue("partition_key", dat.partitionKey).
				SetMetadataKeyValue("stream_name", dat.streamName),
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

func TestClient_PutRecords(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	rm := make(map[string][]byte)
	data := `{"my_result":"ok"}`
	data2 := `{"my_result2":"ok!"}`

	rm["1"] = []byte(data)
	rm["2"] = []byte(data2)
	b, err := json.Marshal(rm)
	require.NoError(t, err)
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid putRecords",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_record").
				SetMetadataKeyValue("partition_key", dat.partitionKey).
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetData(b),
			wantErr: false,
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

func TestClient_GetRecords(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-kinesis",
		Kind: "aws.kinesis",
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
			name: "valid get records",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_records").
				SetMetadataKeyValue("stream_name", dat.streamName).
				SetMetadataKeyValue("shard_iterator_id", dat.shardIterator),
			wantErr: false,
		},
		{
			name: "invalid get records - missing stream name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_records").
				SetMetadataKeyValue("shard_iterator_id", dat.shardIterator),
			wantErr: true,
		},
		{
			name: "invalid get records - missing iterator",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "get_records").
				SetMetadataKeyValue("stream_name", dat.streamName),
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
