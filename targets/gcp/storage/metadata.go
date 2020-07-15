package firestore

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

var methodsMap = map[string]string{
	"upload":   "upload",
	"download": "download",
	"delete":   "delete",
	"rename":   "rename",
	"copy":     "copy",
	"move":     "move",
	"list":     "list",
}

type metadata struct {
	method       string
	object       string
	renameObject string
	bucket       string

	dstBucket string
	path      string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	if m.method == "upload" || m.method == "download" || m.method == "delete" || m.method == "rename" || m.method == "copy" {
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

	return m, nil
}
