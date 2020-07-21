package sqs

import (
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sqsQueue := os.Getenv("SQS_QUEUE_NAME")
	deadLetter := os.Getenv("DEAD_LETTER")
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "sqs-target",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"queue":                       sqsQueue,
					"region":                      "us-west-2",
					"max_retries":                 "0",
					"max_receive":                 "10",
					"dead_letter":                 deadLetter,
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error no region",
			cfg: config.Spec{
				Name: "sqs-target",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"queue":                       sqsQueue,
					"region":                      "",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init - error no queue",
			cfg: config.Spec{
				Name: "sqs-target",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"queue":                       "",
					"region":                      "us-west-2",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		}, {
			name: "init - error no sqs_key",
			cfg: config.Spec{
				Name: "sqs-target",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                     "",
					"sqs_secret_key":              awsSecret,
					"queue":                       sqsQueue,
					"region":                      "us-west-2",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init -error no sqs_secret_key",
			cfg: config.Spec{
				Name: "sqs-target",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              "",
					"queue":                       sqsQueue,
					"region":                      "us-west-2",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
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

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Do(t *testing.T) {
	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sqsQueue := os.Getenv("SQS_QUEUE_NAME")

	validBody, _ := json.Marshal("valid body")
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid sqs sent",
			cfg: config.Spec{
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"max_receive":                 "10",
					"dead_letter":                 "dead_letter",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", sqsQueue).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		},
		{
			name: "valid sqs sent - tags",
			cfg: config.Spec{
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"max_receive":                 "10",
					"dead_letter":                 "dead_letter",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", sqsQueue).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		},
		{
			name: "valid sqs sent",
			cfg: config.Spec{
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"max_receive":                 "10",
					"dead_letter":                 "dead_letter",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", sqsQueue).
				SetData(validBody),
			want: types.NewResponse().
				SetData(validBody),

			wantErr: false,
		}, {
			name: "incorrect signature",
			cfg: config.Spec{
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     "Incorrect",
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"max_receive":                 "10",
					"dead_letter":                 "dead_letter",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
				SetMetadataKeyValue("delay", "0").
				SetMetadataKeyValue("queue", sqsQueue).
				SetData(validBody),
			want: nil,

			wantErr: true,
		}, {
			name: "incorrect queue",
			cfg: config.Spec{
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"max_receive":                 "10",
					"dead_letter":                 "dead_letter",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "1",
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
			err := c.Init(ctx, tt.cfg)
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
	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sqsQueue := os.Getenv("SQS_QUEUE_NAME")
	deadLetter := os.Getenv("DEAD_LETTER_QUEUE")
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
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"max_receive":                 "10",
					"dead_letter":                 deadLetter,
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "0",
				},
			},
			queueURL: sqsQueue,
			wantErr:  false,
		}, {
			name: "in-valid set queue attribute",
			cfg: config.Spec{
				Name: "target.sqs",
				Kind: "target.sqs",
				Properties: map[string]string{
					"sqs_key":                     aswKey,
					"sqs_secret_key":              awsSecret,
					"region":                      "us-west-2",
					"dead_letter":                 deadLetter,
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
					"retries":                     "0",
				},
			},
			queueURL: sqsQueue,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			err = c.SetQueueAttributes(ctx, sqsQueue)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
