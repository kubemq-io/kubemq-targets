package filesystem

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "filesystem-target",
				Kind: "",
				Properties: map[string]string{
					"base_path": "./",
				},
			},
			wantErr: false,
		},
		{
			name: "init",
			cfg: config.Spec{
				Name: "filesystem-target",
				Kind: "",
				Properties: map[string]string{
					"base_path": "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSaveErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
func TestClient_Files(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.Spec{
		Name: "filesystem-target",
		Kind: "",
		Properties: map[string]string{
			"base_path": "../../../test-fs",
		},
	}
	c := New()
	err := c.Init(ctx, cfg)
	require.NoError(t, err)
	save := types.NewRequest().
		SetMetadataKeyValue("method", "save").
		SetMetadataKeyValue("path", "test").
		SetMetadataKeyValue("filename", "test.txt").
		SetData([]byte("test"))
	resp, err := c.Do(ctx, save)
	require.NoError(t, err)
	require.Equal(t, false, resp.IsError)

	load := types.NewRequest().
		SetMetadataKeyValue("method", "load").
		SetMetadataKeyValue("path", "test").
		SetMetadataKeyValue("filename", "test.txt")
	resp, err = c.Do(ctx, load)
	require.NoError(t, err)
	require.Equal(t, false, resp.IsError)
	require.Equal(t, []byte("test"), resp.Data)

	del := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("path", "test").
		SetMetadataKeyValue("filename", "test.txt")
	resp, err = c.Do(ctx, del)
	require.NoError(t, err)
	require.Equal(t, false, resp.IsError)

	saveErr := types.NewRequest().
		SetMetadataKeyValue("method", "save").
		SetMetadataKeyValue("path", "*").
		SetMetadataKeyValue("filename", "bad-filename")
	resp, err = c.Do(ctx, saveErr)
	require.NoError(t, err)
	require.Equal(t, true, resp.IsError)

	loadErr := types.NewRequest().
		SetMetadataKeyValue("method", "load").
		SetMetadataKeyValue("path", "test").
		SetMetadataKeyValue("filename", "bad-filename")
	resp, err = c.Do(ctx, loadErr)
	require.NoError(t, err)
	require.Equal(t, true, resp.IsError)

	delErr := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("path", "test").
		SetMetadataKeyValue("filename", "bad-filename")
	resp, err = c.Do(ctx, delErr)
	require.NoError(t, err)
	require.Equal(t, true, resp.IsError)
}
