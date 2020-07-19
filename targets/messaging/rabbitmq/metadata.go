package rabbitmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/streadway/amqp"
	"math"
	"strconv"
	"time"
)

type metadata struct {
	queue         string
	exchange      string
	confirm       bool
	mandatory     bool
	immediate     bool
	deliveryMode  int
	priority      int
	correlationId string
	replyTo       string
	expiration    time.Duration
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.queue, err = meta.MustParseString("queue")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing queue name, %w", err)
	}
	m.exchange = meta.ParseString("exchange", "")
	m.mandatory = meta.ParseBool("mandatory", false)
	m.confirm = meta.ParseBool("confirm", false)
	m.immediate = meta.ParseBool("immediate", false)
	m.deliveryMode, err = meta.ParseIntWithRange("delivery_mode", 1, 1, 2)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing delivery mode, %w", err)
	}
	m.priority, err = meta.ParseIntWithRange("priority", 0, 0, 9)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing priority, %w", err)
	}
	m.correlationId = meta.ParseString("correlation_id", "")
	m.correlationId = meta.ParseString("reply_to", "")
	expirySeconds, err := meta.ParseIntWithRange("expiry_seconds", math.MaxInt32, 0, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing expiry_seconds, %w", err)
	}
	m.expiration = time.Duration(expirySeconds) * time.Second
	return m, nil
}

func (m metadata) amqpMessage(data []byte) amqp.Publishing {
	return amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    uint8(m.deliveryMode),
		Priority:        uint8(m.priority),
		CorrelationId:   m.correlationId,
		ReplyTo:         m.replyTo,
		Expiration:      strconv.FormatInt(m.expiration.Milliseconds(), 10),
		Body:            data,
	}
}
