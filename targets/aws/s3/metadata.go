package s3

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string

	waitForCompletion bool
	bucketName        string
	copySource        string

	itemName string
}

var methodsMap = map[string]string{
	"list_buckets":                 "list_buckets",
	"list_bucket_items":            "list_bucket_items",
	"create_bucket":                "create_bucket",
	"delete_bucket":                "delete_bucket",
	"delete_item_from_bucket":      "delete_item_from_bucket",
	"delete_all_items_from_bucket": "delete_all_items_from_bucket",
	"upload_item":                  "upload_item",
	"copy_item":                    "copy_item",
	"get_item":                     "get_item",
}


func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidHttpMethodTypes(methodsMap)
	}
	if m.method != "list_buckets" {
		m.bucketName, err = meta.MustParseString("bucket_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing bucket_name, %w", err)
		}
		m.waitForCompletion = meta.ParseBool("wait_for_completion", false)
		if m.method == "upload_item" || m.method == "copy" || m.method == "delete_item_from_bucket" || m.method == "copy_item" || m.method == "get_item" {
			m.itemName, err = meta.MustParseString("item_name")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing item_name, %w", err)
			}
			if m.method == "copy_item" {
				m.copySource, err = meta.MustParseString("copy_source")
				if err != nil {
					return metadata{}, fmt.Errorf("error parsing copy_source, %w", err)
				}
			}
		}
	}
	return m, nil
}
