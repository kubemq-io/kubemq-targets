package storage

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
	object       string
	renameObject string
	bucket       string
	dstBucket    string
	filePath     string
	projectID    string
	storageClass string
	location     string
	cred         string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/storage/object.txt")
	if err != nil {
		return nil, err
	}
	t.object = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/renameObject.txt")
	if err != nil {
		return nil, err
	}
	t.renameObject = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/bucket.txt")
	if err != nil {
		return nil, err
	}
	t.bucket = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/dstBucket.txt")
	if err != nil {
		return nil, err
	}
	t.dstBucket = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/filePath.txt")
	if err != nil {
		return nil, err
	}
	t.filePath = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/projectID.txt")
	if err != nil {
		return nil, err
	}
	t.projectID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/storage_class.txt")
	if err != nil {
		return nil, err
	}
	t.storageClass = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/storage/location.txt")
	if err != nil {
		return nil, err
	}
	t.location = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	if err != nil {
		return nil, err
	}
	t.cred = fmt.Sprintf("%s", dat)
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
				Name: "google-storage-target",
				Kind: "",
				Properties: map[string]string{
					"credentials": dat.cred,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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

func TestClient_Create_Bucket(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid create bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "create_bucket").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("storage_class", dat.storageClass).
				SetMetadataKeyValue("project_id", dat.projectID).
				SetMetadataKeyValue("location", dat.location),
			wantErr: false,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
	err = c.client.Close()
	require.NoError(t, err)
}

func TestClient_Upload_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid upload object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("path", dat.filePath).
				SetMetadataKeyValue("object", dat.object),
			wantErr: false,
		}, {
			name: "invalid upload object - missing bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("path", dat.filePath).
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		}, {
			name: "invalid upload object - bucket dont exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("bucket", "not-real-bucket").
				SetMetadataKeyValue("path", dat.filePath).
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		}, {
			name: "invalid upload object - missing file path",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		}, {
			name: "invalid upload object - incorrect file path",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("path", "not/real/path").
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		}, {
			name: "invalid upload object - missing object name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "upload").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("path", dat.filePath),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	require.NoError(t, err)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Delete_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid delete object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object),
			wantErr: false,
		}, {
			name: "invalid delete object - object doesn't exist ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", "madeUpObject"),
			wantErr: true,
		}, {
			name: "invalid delete object - missing object ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("bucket", dat.bucket),
			wantErr: true,
		}, {
			name: "invalid delete object - missing bucket ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("object", "madeUpObject"),
			wantErr: true,
		}, {
			name: "invalid delete object - bucket does not exists ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "delete").
				SetMetadataKeyValue("bucket", "madeup-123").
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Download_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid download object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "download").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object),
			wantErr: false,
		}, {
			name: "invalid download  - missing bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "download").
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		}, {
			name: "invalid download object - bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "download").
				SetMetadataKeyValue("bucket", "notreal-123").
				SetMetadataKeyValue("object", dat.object),
			wantErr: true,
		}, {
			name: "invalid download object - missing object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "download").
				SetMetadataKeyValue("bucket", dat.bucket),
			wantErr: true,
		}, {
			name: "invalid download object - object does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "download").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", "not-real-object"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.NotNil(t, gotSetResponse.Data)
		})
	}
}

func TestClient_List_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid list object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list").
				SetMetadataKeyValue("bucket", dat.bucket),
			wantErr: false,
		}, {
			name: "invalid list object - missing bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list"),
			wantErr: true,
		}, {
			name: "invalid list object - bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "list").
				SetMetadataKeyValue("bucket", "not-real-123"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.NotNil(t, gotSetResponse.Data)
		})
	}
}

func TestClient_Rename_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid rename object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "rename").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: false,
		}, {
			name: "invalid rename object - missing bucket name",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "rename").
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid rename object - bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "rename").
				SetMetadataKeyValue("bucket", "bucket-not-real-123").
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid rename object - missing object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "rename").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid rename object - object does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "rename").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", "object-not-exits-123").
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid rename object - missing rename object",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "rename").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", "object-not-exits-123"),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Copy_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid copy object- same bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("dst_bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: false,
		}, {
			name: "valid copy object- another bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("dst_bucket", dat.dstBucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: false,
		}, {
			name: "invalid copy object - missing origin bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("dst_bucket", dat.dstBucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid copy object -  origin bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("dst_bucket", dat.dstBucket).
				SetMetadataKeyValue("bucket", "not-real-123").
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid copy object -  missing dst bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid copy object -  dst bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("dst_bucket", "not-real-123").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}

func TestClient_Move_Object(t *testing.T) {
	dat, err := getTestStructure()
	cfg2 := config.Spec{
		Name: "google-storage-target",
		Kind: "",
		Properties: map[string]string{
			"credentials": dat.cred,
		},
	}
	require.NoError(t, err)
	tests := []struct {
		name    string
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid move object ",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "move").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("dst_bucket", dat.dstBucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: false,
		}, {
			name: "invalid move object - missing origin bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("dst_bucket", dat.dstBucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid move object -  origin bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("dst_bucket", dat.dstBucket).
				SetMetadataKeyValue("bucket", "not-real-123").
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid move object -  missing dst bucket",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		}, {
			name: "invalid move object -  dst bucket does not exists",
			request: types.NewRequest().
				SetMetadataKeyValue("method", "copy").
				SetMetadataKeyValue("dst_bucket", "not-real-123").
				SetMetadataKeyValue("bucket", dat.bucket).
				SetMetadataKeyValue("object", dat.object).
				SetMetadataKeyValue("rename_object", dat.renameObject),
			wantErr: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := New()
	err = c.Init(ctx, cfg2)
	defer func() {
		err = c.client.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSetResponse, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				t.Logf("init() error = %v, wantErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
		})
	}
}
