package minio

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
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: false,
		},
		{
			name: "init - no endpoint key",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad endpoint key",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "badhost",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no access key",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no secret key",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":      "localhost:9001",
					"access_key_id": "minio",
					"use_ssl":       "false",
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
				t.Errorf("Init() error = %v, wantPutErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
func TestClient_Objects(t *testing.T) {
	tests := []struct {
		name               string
		cfg                config.Spec
		putRequest         *types.Request
		getRequest         *types.Request
		removeRequest      *types.Request
		wantPutResponse    *types.Response
		wantGetResponse    *types.Response
		wantRemoveResponse *types.Response
		wantPutErr         bool
		wantGetErr         bool
		wantRemoveErr      bool
	}{
		{
			name: "valid set get remove request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			putRequest: types.NewRequest().
				SetMetadataKeyValue("method", "put").
				SetMetadataKeyValue("param1", "bucket").
				SetMetadataKeyValue("param2", "test").
				SetData([]byte("test data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("param1", "bucket").
				SetMetadataKeyValue("param2", "test"),
			removeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "remove").
				SetMetadataKeyValue("param1", "bucket").
				SetMetadataKeyValue("param2", "test"),
			wantPutResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok").
				SetData([]byte("test data")),
			wantRemoveResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			wantPutErr:    false,
			wantGetErr:    false,
			wantRemoveErr: false,
		},
		{
			name: "invalid put request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			putRequest: types.NewRequest().
				SetMetadataKeyValue("method", "put").
				SetMetadataKeyValue("param1", "").
				SetMetadataKeyValue("param2", "invalid_test"),
			getRequest:         nil,
			removeRequest:      nil,
			wantPutResponse:    nil,
			wantGetResponse:    nil,
			wantRemoveResponse: nil,
			wantPutErr:         true,
			wantGetErr:         false,
			wantRemoveErr:      false,
		},
		{
			name: "invalid get  request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			putRequest: nil,
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("param1", "").
				SetMetadataKeyValue("param2", "invalid_test"),
			removeRequest:      nil,
			wantPutResponse:    nil,
			wantGetResponse:    nil,
			wantRemoveResponse: nil,
			wantPutErr:         false,
			wantGetErr:         true,
			wantRemoveErr:      false,
		},
		{
			name: "invalid get request - 2",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			putRequest: nil,
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("param1", "invalid_bucket").
				SetMetadataKeyValue("param2", "invalid_test"),
			removeRequest:      nil,
			wantPutResponse:    nil,
			wantGetResponse:    nil,
			wantRemoveResponse: nil,
			wantPutErr:         false,
			wantGetErr:         true,
			wantRemoveErr:      false,
		},
		{
			name: "invalid remove request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			putRequest: nil,
			getRequest: nil,
			removeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "remove").
				SetMetadataKeyValue("param1", "").
				SetMetadataKeyValue("param2", "invalid_test"),
			wantPutResponse:    nil,
			wantGetResponse:    nil,
			wantRemoveResponse: nil,
			wantPutErr:         false,
			wantGetErr:         false,
			wantRemoveErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			if tt.putRequest != nil {
				gotSetResponse, err := c.Do(ctx, tt.putRequest)
				if tt.wantPutErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, gotSetResponse)
				require.EqualValues(t, tt.wantPutResponse, gotSetResponse)
			}
			if tt.getRequest != nil {
				gotGetResponse, err := c.Do(ctx, tt.getRequest)
				if tt.wantGetErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, gotGetResponse)
				require.EqualValues(t, tt.wantGetResponse, gotGetResponse)
			}

			if tt.removeRequest != nil {
				gotRemoveResponse, err := c.Do(ctx, tt.removeRequest)
				if tt.wantRemoveErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, gotRemoveResponse)
				require.EqualValues(t, tt.wantRemoveResponse, gotRemoveResponse)
			}
		})
	}
}
func TestClient_Buckets(t *testing.T) {
	tests := []struct {
		name                string
		cfg                 config.Spec
		makeRequest         *types.Request
		listBucketRequest   *types.Request
		existRequest        *types.Request
		removeRequest       *types.Request
		listObjectsRequest  *types.Request
		makeResponse        *types.Response
		listBucketResponse  *types.Response
		existResponse       *types.Response
		removeResponse      *types.Response
		listObjectsResponse *types.Response
		makeErr             bool
		listBucketErr       bool
		existErr            bool
		removeErr           bool
		listObjectsErr      bool
	}{
		{
			name: "valid make exist list remove request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			makeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "make_bucket").
				SetMetadataKeyValue("param1", "testbucket"),
			listBucketRequest: types.NewRequest().
				SetMetadataKeyValue("method", "list_buckets"),
			existRequest: types.NewRequest().
				SetMetadataKeyValue("method", "bucket_exists").
				SetMetadataKeyValue("param1", "testbucket"),
			removeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "remove_bucket").
				SetMetadataKeyValue("param1", "testbucket"),
			listObjectsRequest: types.NewRequest().
				SetMetadataKeyValue("method", "list_objects").
				SetMetadataKeyValue("param1", "testbucket"),
			makeResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			listBucketResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			existResponse: types.NewResponse().
				SetMetadataKeyValue("exist", "true").
				SetMetadataKeyValue("result", "ok"),
			removeResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			listObjectsResponse: types.NewResponse().
				SetMetadataKeyValue("result", "ok"),
			makeErr:        false,
			listBucketErr:  false,
			existErr:       false,
			removeErr:      false,
			listObjectsErr: false,
		},
		{
			name: "invalid make",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			makeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "make_bucket").
				SetMetadataKeyValue("param1", "__bucket"),
			listBucketRequest:   nil,
			existRequest:        nil,
			removeRequest:       nil,
			listObjectsRequest:  nil,
			makeResponse:        nil,
			listBucketResponse:  nil,
			existResponse:       nil,
			removeResponse:      nil,
			listObjectsResponse: nil,
			makeErr:             true,
			listBucketErr:       false,
			existErr:            false,
			removeErr:           false,
			listObjectsErr:      false,
		},
		{
			name: "invalid exist",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			makeRequest:       nil,
			listBucketRequest: nil,
			existRequest: types.NewRequest().
				SetMetadataKeyValue("method", "bucket_exists").
				SetMetadataKeyValue("param1", "__bucket"),
			removeRequest:       nil,
			listObjectsRequest:  nil,
			makeResponse:        nil,
			listBucketResponse:  nil,
			existResponse:       nil,
			removeResponse:      nil,
			listObjectsResponse: nil,
			makeErr:             false,
			listBucketErr:       false,
			existErr:            true,
			removeErr:           false,
			listObjectsErr:      false,
		},
		{
			name: "invalid remove",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			makeRequest:       nil,
			listBucketRequest: nil,
			existRequest:      nil,
			removeRequest: types.NewRequest().
				SetMetadataKeyValue("method", "remove_bucket").
				SetMetadataKeyValue("param1", "__bucket"),
			listObjectsRequest:  nil,
			makeResponse:        nil,
			listBucketResponse:  nil,
			existResponse:       nil,
			removeResponse:      nil,
			listObjectsResponse: nil,
			makeErr:             false,
			listBucketErr:       false,
			existErr:            false,
			removeErr:           true,
			listObjectsErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			if tt.makeRequest != nil {
				response, err := c.Do(ctx, tt.makeRequest)
				if tt.makeErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, response)
				require.EqualValues(t, tt.removeResponse, response)
			}
			if tt.existRequest != nil {
				response, err := c.Do(ctx, tt.existRequest)
				if tt.existErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, response)
				require.EqualValues(t, tt.existResponse, response)
			}
			if tt.listBucketRequest != nil {
				response, err := c.Do(ctx, tt.listBucketRequest)
				if tt.listBucketErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, response)
				require.EqualValues(t, tt.listBucketResponse.Metadata, response.Metadata)
			}
			if tt.removeRequest != nil {
				response, err := c.Do(ctx, tt.removeRequest)
				if tt.removeErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, response)
				require.EqualValues(t, tt.removeResponse, response)
			}
			if tt.listObjectsRequest != nil {
				response, err := c.Do(ctx, tt.listObjectsRequest)
				if tt.listObjectsErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotNil(t, response)
				require.EqualValues(t, tt.listObjectsResponse.Metadata, response.Metadata)
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Spec
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list_objects").
				SetMetadataKeyValue("param1", "bucket"),
			wantErr: false,
		},
		{
			name: "invalid metadata request",
			cfg: config.Spec{
				Name: "minio",
				Kind: "minio",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("param1", "bucket"),
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
			_, err = c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
