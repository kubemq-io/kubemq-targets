package storage

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"upload":        "upload",
	"create_bucket": "create_bucket",
	"download":      "download",
	"delete":        "delete",
	"rename":        "rename",
	"copy":          "copy",
	"move":          "move",
	"list":          "list",
}

type metadata struct {
	method       string
	object       string
	renameObject string
	bucket       string
	dstBucket    string
	path         string
	projectID    string
	storageClass string
	location     string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method == "create_bucket" {
		m.bucket, err = meta.MustParseString("bucket")
		if err != nil {
			return metadata{}, fmt.Errorf("bucket is required for method:%s , error on parsing bucket, %w", m.method, err)
		}
		m.projectID, err = meta.MustParseString("project_id")
		if err != nil {
			return metadata{}, fmt.Errorf("project_id is required for method:%s, error on project_id, %w", m.method, err)
		}
		m.storageClass, err = meta.MustParseString("storage_class")
		if err != nil {
			return metadata{}, fmt.Errorf("storage_class is required for method:%s, error on storage_class, %w", m.method, err)
		}
		m.location, err = meta.MustParseString("location")
		if err != nil {
			return metadata{}, fmt.Errorf("location is required for method:%s, error on location, %w", m.method, err)
		}
	}
	if m.method == "upload" || m.method == "download" || m.method == "delete" || m.method == "rename" || m.method == "copy" || m.method == "move" {
		m.object, err = meta.MustParseString("object")
		if err != nil {
			return metadata{}, fmt.Errorf("object is required for method:%s, error on parsing upload, %w", m.method, err)
		}
		m.bucket, err = meta.MustParseString("bucket")
		if err != nil {
			return metadata{}, fmt.Errorf("bucket is required for method:%s, error on parsing bucket, %w", m.method, err)
		}
		if m.method == "upload" {
			m.path, err = meta.MustParseString("path")
			if err != nil {
				return metadata{}, fmt.Errorf("path is required for method:%s,error on path, %w", m.method, err)
			}
		} else if m.method == "rename" || m.method == "copy" || m.method == "move" {
			m.renameObject, err = meta.MustParseString("rename_object")
			if err != nil {
				return metadata{}, fmt.Errorf("rename_object is required for method:%s,error on rename_object, %w", m.method, err)
			}
			if m.method == "copy" || m.method == "move" {
				m.dstBucket, err = meta.MustParseString("dst_bucket")
				if err != nil {
					return metadata{}, fmt.Errorf("dst_bucket is required for method:%s,error on dst_bucket, %w", m.method, err)
				}
			}
		}
	}
	if m.method == "list" {
		m.bucket, err = meta.MustParseString("bucket")
		if err != nil {
			return metadata{}, fmt.Errorf("bucket is required for method:%s,error on parsing bucket, %w", m.method, err)
		}
	}

	return m, nil
}
