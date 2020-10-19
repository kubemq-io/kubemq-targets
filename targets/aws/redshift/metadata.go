package redshift

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method       string
	resourceName string
	resourceARN  string
}

var methodsMap = map[string]string{
	"create_tags":                   "create_tags",
	"delete_tags":                   "delete_tags",
	"list_tags":                     "list_tags",
	"list_snapshots":                "list_snapshots",
	"list_snapshots_by_tags_keys":   "list_snapshots_by_tags_keys",
	"list_snapshots_by_tags_values": "list_snapshots_by_tags_values",
	"describe_cluster":              "describe_cluster",
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

	if m.method == "create_tags" || m.method == "delete_tags" {
		m.resourceARN, err = meta.MustParseString("resource_arn")
		if err != nil {
			return metadata{}, fmt.Errorf("resource_arn is required when using :%s,error parsing resource_arn, %w", m.method, err)
		}
	} else if m.method == "describe_cluster" {
		m.resourceName, err = meta.MustParseString("resource_name")
		if err != nil {
			return metadata{}, fmt.Errorf("resource_name is required when using :%s ,error parsing resource_name, %w", m.method, err)
		}
	}

	return m, nil
}
