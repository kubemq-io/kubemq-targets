package amqp

import (
	"fmt"
	"github.com/Azure/go-amqp"
	"github.com/kubemq-io/kubemq-targets/types"
	"time"
)

type metadata struct {
	address        string
	durable        bool
	priority       int
	messageId      string
	to             string
	subject        string
	replyTo        string
	correlationId  string
	contentType    string
	expiryTime     int64
	groupId        string
	groupSequence  int
	replyToGroupId string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.address, err = meta.MustParseString("address")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing address name, %w", err)
	}
	m.durable = meta.ParseBool("durable", false)
	p := meta.ParseInt("priority", 0)
	if p < 0 {
		return metadata{}, fmt.Errorf("error parsing priority, value cannot be lower than 0")
	}
	m.priority = p
	m.messageId = meta.ParseString("message_id", "")
	m.to = meta.ParseString("to", "")
	m.subject = meta.ParseString("subject", "")
	m.replyTo = meta.ParseString("reply_to", "")
	m.correlationId = meta.ParseString("correlation_id", "")
	m.contentType = meta.ParseString("content_type", "")
	m.expiryTime = meta.ParseInt64("expiry_time", 0)
	m.groupId = meta.ParseString("group_id", "")
	m.groupSequence = meta.ParseInt("group_sequence", 0)
	m.replyToGroupId = meta.ParseString("reply_to_group_id", "")
	return m, nil
}

func (m metadata) amqpMessage(data []byte) *amqp.Message {
	msg := amqp.NewMessage(data)
	msg.Header = &amqp.MessageHeader{}
	msg.Properties = &amqp.MessageProperties{}
	if m.durable {
		msg.Header.Durable = true
	}
	if m.priority != 0 {
		msg.Header.Priority = uint8(m.priority)
	}

	if m.messageId != "" {
		msg.Properties.MessageID = m.messageId
	}
	if m.to != "" {
		msg.Properties.To = &m.to
	}
	if m.subject != "" {
		msg.Properties.Subject = &m.subject
	}
	if m.replyTo != "" {
		msg.Properties.ReplyTo = &m.replyTo
	}
	if m.correlationId != "" {
		msg.Properties.CorrelationID = &m.correlationId
	}
	if m.contentType != "" {
		msg.Properties.ContentType = &m.contentType
	}
	if m.expiryTime != 0 {
		t := time.Unix(m.expiryTime, 0)
		msg.Properties.AbsoluteExpiryTime = &t
	}
	if m.groupId != "" {
		msg.Properties.GroupID = &m.groupId

	}
	if m.groupSequence != 0 {
		val := uint32(m.groupSequence)
		msg.Properties.GroupSequence = &val
	}
	if m.replyToGroupId != "" {
		msg.Properties.ReplyToGroupID = &m.replyToGroupId
	}

	return msg
}
