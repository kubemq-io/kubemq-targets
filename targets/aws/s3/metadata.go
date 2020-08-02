package s3

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string

	waitForCompletion bool
	bucketName        string

	itemName string
}

var methodsMap = map[string]string{
	"list_buckets": "list_buckets",
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
		return metadata{}, fmt.Errorf(getValidMethodTypes())
	}
	if m.method != "list_buckets" {
		m.bucketName, err = meta.MustParseString("bucket_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing bucket_name, %w", err)
		}
		m.waitForCompletion = meta.ParseBool("wait_for_completion", false)
	}
	if m.method == "upload_item" {
		m.itemName, err = meta.MustParseString("item_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing item_name, %w", err)
		}
	}
	return m, nil
}
