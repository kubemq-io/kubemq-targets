package sns

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string
	topic  string

	endPoint           string
	protocol           string
	returnSubscription bool

	message           string
	phoneNumber       string
	subject           string
	targetArn         string
}

var methodsMap = map[string]string{
	"list_topics":                 "list_topics",
	"list_subscriptions":          "list_subscriptions",
	"list_subscriptions_by_topic": "list_subscriptions_by_topic",
	"create_topic":                "create_topic",
	"subscribe":                   "subscribe",
	"send_message":                "send_message",
	"delete_topic":                "delete_topic",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method != "list_topics" && m.method != "list_subscriptions" && m.method != "send_message" {
		m.topic, err = meta.MustParseString("topic")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing topic, %w", err)
		}
		if m.method == "subscribe" {
			m.endPoint = meta.ParseString("end_point","")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing end_point, %w", err)
			}
			m.protocol, err = meta.MustParseString("protocol")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing protocol, %w", err)
			}
			m.returnSubscription, err = meta.MustParseBool("return_subscription")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing return_subscription, %w", err)
			}
		}
	} else if m.method == "send_message" {
		m.targetArn, err = meta.MustParseString("target_arn")
		if err != nil {
			m.topic, err = meta.MustParseString("topic")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing topic or target_arn , one of them must be set , %w", err)
			}
		} else {
			m.topic, err = meta.MustNotParseString("topic", "target_arn")
			if err == nil {
				return metadata{}, err
			}
		}
		m.message, err = meta.MustParseString("message")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing message, %w", err)
		}
		m.phoneNumber = meta.ParseString("phone_number", "")
		m.subject = meta.ParseString("subject", "")
	}
	return m, nil
}
