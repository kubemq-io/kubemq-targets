package elasticsearch

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

	json     string
	domain   string
	index    string
	endpoint string
	service  string
	id       string
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

	dat, err = ioutil.ReadFile("./../../../credentials/aws/elasticsearch/signer/json.txt")
	if err != nil {
		return nil, err
	}
	t.json = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/elasticsearch/signer/domain.txt")
	if err != nil {
		return nil, err
	}
	t.domain = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/elasticsearch/signer/index.txt")
	if err != nil {
		return nil, err
	}
	t.index = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/elasticsearch/signer/endpoint.txt")
	if err != nil {
		return nil, err
	}
	t.endpoint = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/elasticsearch/signer/service.txt")
	if err != nil {
		return nil, err
	}
	t.service = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/elasticsearch/signer/id.txt")
	if err != nil {
		return nil, err
	}
	t.id = string(dat)

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
				Name: "target-aws-elasticsearch",
				Kind: "target.aws.elasticsearch",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: false,
		}, {
			name: "init - missing key",
			cfg: config.Spec{
				Name: "target-aws-elasticsearch",
				Kind: "target.aws.elasticsearch",
				Properties: map[string]string{
					"aws_secret_key": dat.awsSecretKey,
				},
			},
			wantErr: true,
		}, {
			name: "init - missing secret key",
			cfg: config.Spec{
				Name: "target-aws-elasticsearch",
				Kind: "target.aws.elasticsearch",
				Properties: map[string]string{
					"aws_key": dat.awsKey,
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

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "target-aws-elasticsearch",
		Kind: "target.aws.elasticsearch",
		Properties: map[string]string{
			"aws_key":        dat.awsKey,
			"aws_secret_key": dat.awsSecretKey,
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
			name: "valid post",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: false,
		}, {
			name: "invalid post - missing method",
			request: types.NewRequest().
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		},
		{
			name: "invalid post - incorrect method",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "MADE_UP").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		},
		{
			name: "invalid post - missing region",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		}, {
			name: "invalid post - missing json",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		}, {
			name: "invalid post - missing domain",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		}, {
			name: "invalid post - missing endpoint",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		}, {
			name: "invalid post - missing index",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("service", dat.service).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		}, {
			name: "invalid post - missing service ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("id", dat.id),
			wantErr: true,
		}, {
			name: "invalid post - missing id",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "POST").
				SetMetadataKeyValue("region", dat.region).
				SetMetadataKeyValue("json", dat.json).
				SetMetadataKeyValue("domain", dat.domain).
				SetMetadataKeyValue("endpoint", dat.endpoint).
				SetMetadataKeyValue("index", dat.index).
				SetMetadataKeyValue("service", dat.service),
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
