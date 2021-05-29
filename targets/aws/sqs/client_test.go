package sqs

import (
	"context"
	"encoding/json"
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

	sqsQueue   string
	deadLetter string
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
	t.region = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/sqs/queue.txt")
	if err != nil {
		return nil, err
	}
	t.sqsQueue = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/sqs/deadLetter.txt")
	if err != nil {
		return nil, err
	}
	t.deadLetter = string(dat)
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
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"token":          dat.token,
					"region":         dat.region,
					"max_retries":    "0",
					"max_receive":    "10",
					"dead_letter":    dat.deadLetter,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - no region",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"token":          dat.token,
					"region":         "",
					"max_retries":    "0",
				},
			},
			wantErr: true,
		}, {
			name: "invalid  init - no aws_key",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        "",
					"aws_secret_key": dat.awsSecretKey,
					"token":          dat.token,
					"queue":          dat.sqsQueue,
					"region":         dat.region,
					"max_retries":    "0",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid  init - no aws_secret_key",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": "",
					"token":          dat.token,
					"region":         dat.region,
					"max_retries":    "0",
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

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)

	validBody, err := json.Marshal("valid body")
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid sqs sent without tags",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"max_retries":    "0",

					"retries": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", dat.sqsQueue).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		},
		{
			name: "valid sqs sent - with tags",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"max_retries":    "0",
					"retries":        "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", dat.sqsQueue).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		},
		{
			name: "valid sqs sent",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"max_retries":    "0",
					"retries":        "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", dat.sqsQueue).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		}, {
			name: "invalid send - incorrect signature",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        "Incorrect",
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"max_retries":    "0",
					"retries":        "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", dat.sqsQueue).
				SetData(validBody),
			want: nil,

			wantErr: true,
		}, {
			name: "invalid send - incorrect queue",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"max_retries":    "0",
					"retries":        "1",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetData(validBody),
			want: nil,

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
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

func TestClient_SetQueueAttributes(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name     string
		cfg      config.Spec
		queueURL string
		want     *types.Response
		wantErr  bool
	}{
		{
			name: "valid set queue attribute",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"max_receive":    "10",
					"dead_letter":    dat.deadLetter,
					"max_retries":    "0",
					"retries":        "0",
				},
			},
			queueURL: dat.sqsQueue,
			wantErr:  false,
		}, {
			name: "invalid - set queue attribute",
			cfg: config.Spec{
				Name: "aws-sqs",
				Kind: "aws.sqs",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"region":         dat.region,
					"dead_letter":    dat.deadLetter,
					"max_retries":    "0",
					"retries":        "0",
				},
			},
			queueURL: dat.sqsQueue,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			err = c.SetQueueAttributes(ctx, dat.sqsQueue)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
