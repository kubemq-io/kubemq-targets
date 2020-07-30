package s3

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string

	waitForCompletion bool
	bucketName        string
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
		if m.method == "create_bucket" {
			m.waitForCompletion = meta.ParseBool("wait_for_completion", false)
		}
	}
	return m, nil
}
