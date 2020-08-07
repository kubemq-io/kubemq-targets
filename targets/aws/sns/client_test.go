package sns

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

	topic              string
	message            string
	endPoint           string
	protocol           string
	returnSubscription string
}

func createSendMessageAttributes(store string, event string) ([]byte, error) {
	var at []Attributes
	aStore := Attributes{
		Name:        "store",
		StringValue: store,
		DataType:    "String",
	}
	at = append(at, aStore)
	aEvent := Attributes{
		Name:        "event",
		StringValue: event,
		DataType:    "String",
	}
	at = append(at, aEvent)
	b, err := json.Marshal(at)
	if err != nil {
		return nil, err
	}
	return b, nil
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
	dat, err = ioutil.ReadFile("./../../../credentials/aws/sns/topic.txt")
	if err != nil {
		return nil, err
	}
	t.topic = fmt.Sprintf("%s", dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/sns/message.txt")
	if err != nil {
		return nil, err
	}
	t.message = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/sns/email.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = fmt.Sprintf("%s", dat)

	t.protocol = "email"

	t.returnSubscription = "true"
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
			name: "init ",
			cfg: config.Spec{
				Name: "aws-sns",
				Kind: "aws.sns",
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
				Name: "aws-sns",
				Kind: "aws.sns",
				Properties: map[string]string{
					"aws_key": dat.awsKey,
					"region":  dat.region,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing key",
			cfg: config.Spec{
				Name: "aws-sns",
				Kind: "aws.sns",
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
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_List_Topics(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_topics"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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

func TestClient_List_Subscriptions(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_subscriptions"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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

func TestClient_List_Subscriptions_By_Topic(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list by topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_subscriptions_by_topic").
				SetMetadataKeyValue("topic", dat.topic),
			wantErr: false,
		},
		{
			name: "invalid list by topic - missing topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_subscriptions_by_topic"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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

func TestClient_Create_Topic(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	attributes := make(map[string]*string)
	DisplayName := "my-display-name"
	attributes["DisplayName"] = &DisplayName
	b, err := json.Marshal(attributes)
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create topic -with attributes",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_topic").
				SetMetadataKeyValue("topic", dat.topic).
				SetData(b),
			wantErr: false,
		}, {
			name: "valid create topic -without attributes",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_topic").
				SetMetadataKeyValue("topic", fmt.Sprintf("%s-another", dat.topic)),
			wantErr: false,
		}, {
			name: "invalid create topic - missing topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_topic"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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

func TestClient_SendMessage(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	attributes, err := createSendMessageAttributes("my_store", "my_event")
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid send Message- target_arn",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send_message").
				SetMetadataKeyValue("target_arn", dat.topic).
				SetMetadataKeyValue("message", dat.message).
				SetData(attributes),
			wantErr: false,
		}, {
			name: "valid send message - topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "send_message").
				SetMetadataKeyValue("topic", dat.topic).
				SetMetadataKeyValue("message", dat.message).
				SetData(attributes),
			wantErr: false,
		}, {
			name: "invalid send message - missing target_arn",
			request: types.NewRequest().
				SetMetadataKeyValue("message", dat.message).
				SetMetadataKeyValue("method", "send_message"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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

func TestClient_Subscribe(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	attributes := make(map[string]*string)
	RawMessageDelivery := `{
	"store": ["mystore"],
    "event": [{"anything-but": "my-event"}],
}`
	attributes["FilterPolicy"] = &RawMessageDelivery
	b, err := json.Marshal(attributes)
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid subscribe topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "subscribe").
				SetMetadataKeyValue("topic", dat.topic).
				SetMetadataKeyValue("protocol", dat.protocol).
				SetMetadataKeyValue("return_subscription", dat.returnSubscription).
				SetMetadataKeyValue("end_point", dat.endPoint).
				SetData(b),
			wantErr: false,
		}, {
			name: "invalid subscribe topic - missing topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "subscribe").
				SetMetadataKeyValue("protocol", dat.protocol).
				SetMetadataKeyValue("return_subscription", dat.returnSubscription).
				SetMetadataKeyValue("end_point", dat.endPoint),
			wantErr: true,
		},
		{
			name: "invalid subscribe topic - missing protocol",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "subscribe").
				SetMetadataKeyValue("topic", dat.topic).
				SetMetadataKeyValue("return_subscription", dat.returnSubscription).
				SetMetadataKeyValue("end_point", dat.endPoint),
			wantErr: true,
		}, {
			name: "invalid subscribe topic - missing return_subscription",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "subscribe").
				SetMetadataKeyValue("topic", dat.topic).
				SetMetadataKeyValue("protocol", dat.protocol).
				SetMetadataKeyValue("end_point", dat.endPoint),
			wantErr: true,
		}, {
			name: "invalid subscribe topic - missing end_point",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "subscribe").
				SetMetadataKeyValue("topic", dat.topic).
				SetMetadataKeyValue("protocol", dat.protocol).
				SetMetadataKeyValue("return_subscription", dat.returnSubscription),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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

func TestClient_Delete_Topic(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "aws-sns",
		Kind: "aws.sns",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
			"region":         dat.region,
		},
	}
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_topic").
				SetMetadataKeyValue("topic", dat.topic),
			wantErr: false,
		}, {
			name: "invalid delete topic - missing topic",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete_topic"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err = c.Init(ctx, cfg)
			require.NoError(t, err)
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
