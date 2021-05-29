package nats

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io/ioutil"
	"time"

	"github.com/kubemq-hub/kubemq-targets/config"

	"github.com/stretchr/testify/require"
	"testing"
)

type testStructure struct {
	url                string
	subject            string
	username           string
	password           string
	token              string
	sslcertificatefile string
	sslcertificatekey  string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/nats/url.txt")
	if err != nil {
		return nil, err
	}
	t.url = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/subject.txt")
	if err != nil {
		return nil, err
	}
	t.subject = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/username.txt")
	if err != nil {
		return nil, err
	}
	t.username = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/password.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/token.txt")
	if err != nil {
		return nil, err
	}
	t.token = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/certFile.pem")
	if err != nil {
		return nil, err
	}
	t.sslcertificatefile = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/certKey.pem")
	if err != nil {
		return nil, err
	}
	t.sslcertificatekey = string(dat)

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
			name: "init -without tls",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":      dat.url,
					"username": dat.username,
					"password": dat.password,
					"token":    dat.password,
					"tls":      "false",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - no  url",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"username":  dat.username,
					"password":  dat.password,
					"token":     dat.password,
					"tls":       "true",
					"cert_file": dat.sslcertificatefile,
					"cert_key":  dat.sslcertificatekey,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing cert key",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":       dat.url,
					"username":  dat.username,
					"password":  dat.password,
					"token":     dat.password,
					"tls":       "true",
					"cert_file": dat.sslcertificatefile,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing cert file",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":      dat.url,
					"username": dat.username,
					"password": dat.password,
					"token":    dat.password,
					"tls":      "true",
					"cert_key": dat.sslcertificatekey,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "valid -without tls",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":      dat.url,
					"username": dat.username,
					"password": dat.password,
					"token":    dat.password,
					"tls":      "false",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("subject", "foo").
				SetData([]byte("some-data")),
			wantResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
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
