package amazonmq

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type testStructure struct {
	host        string
	username    string
	password    string
	destination string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/aws/amazonmq/host.txt")
	if err != nil {
		return nil, err
	}
	t.host = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/amazonmq/username.txt")
	if err != nil {
		return nil, err
	}
	t.username = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/amazonmq/password.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/amazonmq/destination.txt")
	if err != nil {
		return nil, err
	}
	t.destination = string(dat)
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
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":     dat.host,
					"username": dat.username,
					"password": dat.password,
				},
			},
			wantErr: false,
		}, {
			name: "init - no host",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"username": dat.username,
					"password": dat.password,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantExecErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name         string
		cfg          config.Spec
		request      *types.Request
		wantResponse *types.Response
		wantErr      bool
	}{
		{
			name: "valid publish request",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":     dat.host,
					"username": dat.username,
					"password": dat.password,
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("destination", dat.destination).
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
		{
			name: "invalid publish request - no destination",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":     dat.host,
					"username": dat.username,
					"password": dat.password,
				},
			},
			request: types.NewRequest().
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
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
			defer func() {
				err = c.Stop()
				require.NoError(t, err)
			}()
			gotResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotResponse)
			require.EqualValues(t, tt.wantResponse, gotResponse)
		})
	}
}
