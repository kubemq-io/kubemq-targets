package s3

import (
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string
}

var methodsMap = map[string]string{
	"list_tags":                     "list_tags",
	"list_snapshots":                "list_snapshots",
	"list_snapshots_by_tags_keys":   "list_snapshots_by_tags_keys",
	"list_snapshots_by_tags_values": "list_snapshots_by_tags_values",
	"list_clusters":                 "list_clusters",
	"list_clusters_by_tags_keys":    "list_clusters_by_tags_keys",
	"list_clusters_by_tags_values":  "list_clusters_by_tags_values",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}

	return m, nil
}
