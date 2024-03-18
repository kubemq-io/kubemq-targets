package types

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-io/kubemq-go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	Metadata Metadata `json:"metadata,omitempty"`
	Data     []byte   `json:"data,omitempty"`
}

func NewRequest() *Request {
	return &Request{
		Metadata: NewMetadata(),
		Data:     nil,
	}
}

func (r *Request) SetMetadata(value Metadata) *Request {
	r.Metadata = value
	return r
}

func (r *Request) SetMetadataKeyValue(key, value string) *Request {
	r.Metadata.Set(key, value)
	return r
}

func (r *Request) SetData(value []byte) *Request {
	r.Data = value
	return r
}

func (r *Request) Size() float64 {
	return float64(len(r.Data))
}

func ParseRequest(body []byte) (*Request, error) {
	if body == nil {
		return nil, fmt.Errorf("empty request")
	}
	baseRequest := NewRequest()
	err := json.Unmarshal(body, baseRequest)
	if err == nil {
		return baseRequest, nil
	}
	req := &TransportRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		return nil, fmt.Errorf("invalid request format, %w", err)
	}
	return NewRequest().SetMetadata(req.Metadata).SetData(body), nil
}

func (r *Request) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}

func (r *Request) ToEvent() *kubemq.Event {
	return kubemq.NewEvent().
		SetBody(r.MarshalBinary())
}

func (r *Request) ToEventStore() *kubemq.EventStore {
	return kubemq.NewEventStore().
		SetBody(r.MarshalBinary())
}

func (r *Request) ToCommand() *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(r.MarshalBinary())
}

func (r *Request) ToQuery() *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(r.MarshalBinary())
}

func (r *Request) ToQueueMessage() *kubemq.QueueMessage {
	return kubemq.NewQueueMessage().
		SetBody(r.MarshalBinary())
}

func (r *Request) String() string {
	str, err := json.MarshalToString(r)
	if err != nil {
		return ""
	}
	return str
}

type TransportRequest struct {
	Metadata Metadata    `json:"metadata,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func NewTransportRequest() *TransportRequest {
	return &TransportRequest{
		Metadata: NewMetadata(),
		Data:     nil,
	}
}

func (r *TransportRequest) SetMetadata(value Metadata) *TransportRequest {
	r.Metadata = value
	return r
}

func (r *TransportRequest) SetMetadataKeyValue(key, value string) *TransportRequest {
	r.Metadata.Set(key, value)
	return r
}

func (r *TransportRequest) SetData(value interface{}) *TransportRequest {
	r.Data = value
	return r
}

func (r *TransportRequest) ToEvent() *kubemq.Event {
	return kubemq.NewEvent().
		SetBody(r.MarshalBinary())
}

func (r *TransportRequest) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}
