package s3

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
	address string
	file    []byte
	user    string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	t.address = "localhost:9000"
	contents, err := ioutil.ReadFile("./../../../examples/storage/hdfs/exampleFile.txt")
	if err != nil {
		return nil, err
	}
	t.file = contents
	t.user = "test_user"
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
				Name: "storage-hdfs",
				Kind: "storage.hdfs",
				Properties: map[string]string{
					"address": dat.address,
					"user":    dat.address,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - incorrect port",
			cfg: config.Spec{
				Name: "storage-hdfs",
				Kind: "storage.hdfs",
				Properties: map[string]string{
					"address": "localhost:123123",
					"user":    dat.address,
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

func TestClient_Mkdir(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "storage-hdfs",
		Kind: "storage.hdfs",
		Properties: map[string]string{
			"address": dat.address,
			"user":    dat.address,
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
			name: "valid mkdir",
			request: types.NewRequest().
				SetMetadataKeyValue("file_path", "/hadoop/dfs/name").
				SetMetadataKeyValue("file_mode", "0755").
				SetMetadataKeyValue("method", "mkdir"),
			wantErr: false,
		}, {
			name: "invalid mkdir - missing path",
			request: types.NewRequest().
				SetMetadataKeyValue("file_mode", "0755").
				SetMetadataKeyValue("method", "mkdir"),
			wantErr: true,
		}, {
			name: "invalid mkdir invalid file_mode",
			request: types.NewRequest().
				SetMetadataKeyValue("file_path", "/test").
				SetMetadataKeyValue("file_mode", "99999").
				SetMetadataKeyValue("method", "mkdir"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestClient_Upload(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	cfg := config.Spec{
		Name: "storage-hdfs",
		Kind: "storage.hdfs",
		Properties: map[string]string{
			"address": dat.address,
			"user":    dat.address,
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
			name: "valid upload",
			request: types.NewRequest().
				SetMetadataKeyValue("file_path", "/test/foo.txt").
				SetData(dat.file).
				SetMetadataKeyValue("method", "write_file"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

//
//func TestClient_List_Bucket_Items(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "false",
//		},
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//	defer cancel()
//	c := New()
//
//	err = c.Init(ctx, cfg)
//	require.NoError(t, err)
//	tests := []struct {
//		name    string
//		request *types.Request
//		wantErr bool
//	}{
//		{
//			name: "valid list",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "list_bucket_items").
//				SetMetadataKeyValue("bucket_name", dat.bucketName),
//			wantErr: false,
//		},
//		{
//			name: "invalid list - missing bucket_name",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "list_bucket_items"),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Create_Bucket(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "false",
//		},
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//	defer cancel()
//	c := New()
//
//	err = c.Init(ctx, cfg)
//	require.NoError(t, err)
//	tests := []struct {
//		name    string
//		request *types.Request
//		wantErr bool
//	}{
//		{
//			name: "valid create",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "create_bucket").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr: false,
//		},
//		{
//			name: "invalid create - missing bucket_name",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "create_bucket"),
//			wantErr: true,
//		},
//		{
//			name: "invalid create - Already exists",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "create_bucket").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Upload_Item(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "true",
//		},
//	}
//	tests := []struct {
//		name        string
//		request     *types.Request
//		wantErr     bool
//		setUploader bool
//	}{
//		{
//			name: "valid upload item",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "upload_item").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetMetadataKeyValue("item_name", dat.itemName).
//				SetData(dat.file),
//			wantErr:     false,
//			setUploader: true,
//		},
//		{
//			name: "invalid upload - missing item",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "upload_item").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetData(dat.file),
//			wantErr:     true,
//			setUploader: true,
//		},
//		{
//			name: "invalid upload - missing uploader",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "upload_item").
//				SetMetadataKeyValue("item_name", dat.itemName).
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr:     true,
//			setUploader: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c := New()
//
//			if !tt.setUploader {
//				cfg.Properties["uploader"] = "false"
//			} else {
//				cfg.Properties["uploader"] = "true"
//			}
//			err = c.Init(ctx, cfg)
//			require.NoError(t, err)
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Get_Item(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "true",
//			"uploader":       "false",
//		},
//	}
//	tests := []struct {
//		name          string
//		request       *types.Request
//		wantErr       bool
//		setDownloader bool
//	}{
//		{
//			name: "valid get item",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "get_item").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetMetadataKeyValue("item_name", dat.itemName),
//			wantErr:       false,
//			setDownloader: true,
//		},
//		{
//			name: "invalid get - missing item",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "get_item").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr:       true,
//			setDownloader: true,
//		}, {
//			name: "invalid get - item does not exists",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "get_item").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetMetadataKeyValue("item_name", "fakeItemName"),
//			wantErr:       true,
//			setDownloader: true,
//		},
//		{
//			name: "invalid get - missing bucketName",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "get_item").
//				SetMetadataKeyValue("item_name", dat.itemName),
//			wantErr:       true,
//			setDownloader: true,
//		}, {
//			name: "invalid upload - missing downloader",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "get_item").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetMetadataKeyValue("item_name", dat.itemName),
//			wantErr:       true,
//			setDownloader: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c := New()
//
//			if !tt.setDownloader {
//				cfg.Properties["downloader"] = "false"
//			} else {
//				cfg.Properties["downloader"] = "true"
//			}
//			err = c.Init(ctx, cfg)
//			require.NoError(t, err)
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Delete_Item(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "false",
//		},
//	}
//	tests := []struct {
//		name    string
//		request *types.Request
//		wantErr bool
//	}{
//		{
//			name: "valid delete item",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_item_from_bucket").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetMetadataKeyValue("item_name", dat.itemName),
//			wantErr: false,
//		},
//		{
//			name: "invalid delete - missing item",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_item_from_bucket").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr: true,
//		},
//		{
//			name: "invalid delete - missing bucket_name",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_item_from_bucket").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("item_name", dat.itemName),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c := New()
//
//			err = c.Init(ctx, cfg)
//			require.NoError(t, err)
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Copy_Items(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "false",
//		},
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	c := New()
//
//	err = c.Init(ctx, cfg)
//	require.NoError(t, err)
//	tests := []struct {
//		name    string
//		request *types.Request
//		wantErr bool
//	}{
//		{
//			name: "valid copy items",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "copy_item").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.bucketName).
//				SetMetadataKeyValue("item_name", dat.itemName).
//				SetMetadataKeyValue("copy_source", dat.dstBucketName),
//			wantErr: false,
//		},
//		{
//			name: "invalid copy items - missing copy_source ",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "copy_item").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName).
//				SetMetadataKeyValue("item_name", dat.itemName),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Delete_All_Items(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "false",
//		},
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	c := New()
//
//	err = c.Init(ctx, cfg)
//	require.NoError(t, err)
//	tests := []struct {
//		name    string
//		request *types.Request
//		wantErr bool
//	}{
//		{
//			name: "valid delete all items",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_all_items_from_bucket").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr: false,
//		},
//		{
//			name: "invalid valid delete all items - missing bucket",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_all_items_from_bucket"),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
//
//func TestClient_Delete_Bucket(t *testing.T) {
//	dat, err := getTestStructure()
//	require.NoError(t, err)
//	cfg := config.Spec{
//		Name: "storage-hdfs",
//		Kind: "storage.hdfs",
//		Properties: map[string]string{
//			"aws_key":        dat.awsKey,
//			"aws_secret_key": dat.awsSecretKey,
//			"region":         dat.region,
//			"downloader":     "false",
//			"uploader":       "false",
//		},
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//	defer cancel()
//	c := New()
//
//	err = c.Init(ctx, cfg)
//	require.NoError(t, err)
//	tests := []struct {
//		name    string
//		request *types.Request
//		wantErr bool
//	}{
//		{
//			name: "valid delete",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_bucket").
//				SetMetadataKeyValue("wait_for_completion", "true").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr: false,
//		},
//		{
//			name: "invalid delete - missing bucket_name",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "create_bucket"),
//			wantErr: true,
//		},
//		{
//			name: "invalid delete - bucket does not exists",
//			request: types.NewRequest().
//				SetMetadataKeyValue("method", "delete_bucket").
//				SetMetadataKeyValue("bucket_name", dat.testBucketName),
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := c.Do(ctx, tt.request)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//			require.NotNil(t, got)
//		})
//	}
//}
