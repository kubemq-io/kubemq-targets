package kinesis

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

	shardCount    string
	streamName    string
	partitionKey  string
	shardPosition string
	shardID       string
	limit         string
	streamARN         string
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
			name: "init - missing key",
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
			name: "init - missing region",
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
			name: "init - missing secret key",
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
			require.EqualValues(t, tt.cfg.Name, c.Name())
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
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
