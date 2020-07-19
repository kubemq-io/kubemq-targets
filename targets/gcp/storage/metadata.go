package storage

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
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

func getValidMethodTypes() string {
	s := "invalid method type, method type should be one of the following:"
	for k := range methodsMap {
		s = fmt.Sprintf("%s :%s,", s, k)
	}
	return s
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	if m.method == "create_bucket" {
		m.bucket, err = meta.MustParseString("bucket")
		if err != nil {
			return metadata{}, fmt.Errorf("error on parsing bucket, %w", err)
		}
		m.projectID, err = meta.MustParseString("project_id")
		if err != nil {
			return metadata{}, fmt.Errorf("error on project_id, %w", err)
		}
		m.storageClass, err = meta.MustParseString("storage_class")
		if err != nil {
			return metadata{}, fmt.Errorf("error on storage_class, %w", err)
		}
		m.location, err = meta.MustParseString("location")
		if err != nil {
			return metadata{}, fmt.Errorf("error on location, %w", err)
		}
	}
	if m.method == "upload" || m.method == "download" || m.method == "delete" || m.method == "rename" || m.method == "copy" || m.method == "move" {
		m.object, err = meta.MustParseString("object")
		if err != nil {
			return metadata{}, fmt.Errorf("error on parsing upload, %w", err)
		}
		m.bucket, err = meta.MustParseString("bucket")
		if err != nil {
			return metadata{}, fmt.Errorf("error on parsing bucket, %w", err)
		}
		if m.method == "upload" {
			m.path, err = meta.MustParseString("path")
			if err != nil {
				return metadata{}, fmt.Errorf("error on path, %w", err)
			}
		} else if m.method == "rename" || m.method == "copy" || m.method == "move" {
			m.renameObject, err = meta.MustParseString("rename_object")
			if err != nil {
				return metadata{}, fmt.Errorf("error on rename_object, %w", err)
			}
			if m.method == "copy" || m.method == "move" {
				m.dstBucket, err = meta.MustParseString("dst_bucket")
				if err != nil {
					return metadata{}, fmt.Errorf("error on dst_bucket, %w", err)
				}
			}
		}
	}
	if m.method == "list" {
		m.bucket, err = meta.MustParseString("bucket")
		if err != nil {
			return metadata{}, fmt.Errorf("error on parsing bucket, %w", err)
		}
	}

	return m, nil
}
