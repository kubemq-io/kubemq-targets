package firebase

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_customToken(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
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
		{
			name: "custom token -valid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "custom_token").
				SetMetadataKeyValue("token_id", "some-uid"),
			wantErr: false,
		},
		{
			name: "custom token -valid - missing token",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "custom_token"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Data)

		})
	}
}

func TestClient_verifyToken(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
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
		{
			name: "verify_token -valid",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "verify_token").
				SetMetadataKeyValue("token_id", dat.token),
			wantErr: false,
		},
		{
			name: "verify_token -invalid missing token",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "verify_token"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, r.Data)

		})
	}
}
