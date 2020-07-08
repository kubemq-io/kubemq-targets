package sqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

type metadata struct {
	delay int
	tags  map[string]*sqs.MessageAttributeValue
	queueURL string
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.tags = make(map[string]*sqs.MessageAttributeValue)
	var err error
	m.queueURL, err = meta.MustParseString("queue")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing queue, %w", err)
	}
	m.delay = meta.ParseInt("delay", opts.defaultDelay)
	tags, err := meta.MustParseJsonMap("tags")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing tags, %w", err)
	}
	for k, v := range tags {
		attributeValue := &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
		m.tags[k] = attributeValue
	}
	return m, nil
}
