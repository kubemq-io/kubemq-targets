package events

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

	detail      string
	detailType  string
	source      string
	rule        string
	resourceARN string
	id          string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../../credentials/aws/awsKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/awsSecretKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsSecretKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/region.txt")
	if err != nil {
		return nil, err
	}
	t.region = fmt.Sprintf("%s", dat)
	t.token = ""
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/events/detail.txt")
	if err != nil {
		return nil, err
	}
	t.detail = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/events/detailType.txt")
	if err != nil {
		return nil, err
	}
	t.detailType = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/events/source.txt")
	if err != nil {
		return nil, err
	}
	t.source = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/events/rule.txt")
	if err != nil {
		return nil, err
	}
	t.rule = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../../credentials/aws/cloudwatch/events/resourceARN.txt")
	if err != nil {
		return nil, err
	}
	t.resourceARN = fmt.Sprintf("%s", dat)

	t.id = "my_arn_id"

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
				Name: "aws-cloudwatch-events",
				Kind: "aws.cloudwatch.events",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "aws-cloudwatch-events",
				Kind: "aws.cloudwatch.events",
				Properties: map[string]string{
					"aws_key": dat.awsKey,
					"region":  dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing key",
			cfg: config.Spec{
				Name: "aws-cloudwatch-events",
				Kind: "aws.cloudwatch.events",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
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

func TestClient_PutTargets(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-cloudwatch-events",
		Kind: "aws.cloudwatch.events",
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
	m := make(map[string]string)
	m[dat.id] = dat.resourceARN
	b, err := json.Marshal(m)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid put targets",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_targets").
				SetMetadataKeyValue("rule", dat.rule).
				SetData(b),
			wantErr: false,
		},
		{
			name: "invalid put targets - missing rule",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "put_targets"),
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

func TestClient_ListBuses(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-cloudwatch-events",
		Kind: "aws.cloudwatch.events",
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
			name: "valid list buses",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_buses"),
			wantErr: false,
		},
		{
			name: "valid list buses with prefix",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_buses").
				SetMetadataKeyValue("limit", "1"),
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

func TestClient_SendEvent(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-cloudwatch-events",
		Kind: "aws.cloudwatch.events",
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
	arn := []string{dat.resourceARN}
	b, err := json.Marshal(arn)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid send event",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send_event").
				SetMetadataKeyValue("detail", "{ \"key1\": \"value1\", \"key2\": \"value2\" }").
				SetMetadataKeyValue("detail_type", "appRequestSubmitted").
				SetMetadataKeyValue("source", "kubemq_testing").
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
