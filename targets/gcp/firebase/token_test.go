package firebase

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestClient_customToken(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	type fields struct {
		name string
		opts options
	}
	type args struct {
		ctx  context.Context
		meta metadata
		data []byte
	}
	cfg := config.Spec{
		Name: "google-firebase-target",
		Kind: "target.gcp.firebase",
		Properties: map[string]string{
			"project_id":  dat.projectID,
			"credentials": dat.cred,
			"auth_client": "true",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg)
	require.NoError(t, err)
	tests := []struct {
		name    string
		wantErr bool
		request *types.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				name:       tt.fields.name,
				opts:       tt.fields.opts,
				clientAuth: tt.fields.clientAuth,
				dbClient:   tt.fields.dbClient,
			}
			got, err := c.customToken(tt.args.ctx, tt.args.meta, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("customToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_verifyToken(t *testing.T) {
	type fields struct {
		name       string
		opts       options
		clientAuth *auth.Client
		dbClient   *db.Client
	}
	type args struct {
		ctx  context.Context
		meta metadata
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				name:       tt.fields.name,
				opts:       tt.fields.opts,
				clientAuth: tt.fields.clientAuth,
				dbClient:   tt.fields.dbClient,
			}
			got, err := c.verifyToken(tt.args.ctx, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifyToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
