package types

import (
	b64 "encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-io/kubemq-go"
	"reflect"
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
	req := &TransportRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		fmt.Println("error", err.Error())
		return NewRequest().SetData(body), nil
	}
	switch v := req.Data.(type) {
	case nil:
		return &Request{
			Metadata: req.Metadata,
			Data:     nil,
		}, nil
	case []byte:
		return &Request{
			Metadata: req.Metadata,
			Data:     v,
		}, nil
	case string:
		sDec, err := b64.StdEncoding.DecodeString(v)
		if err != nil {
			sDec = []byte(v)
		}
		return &Request{
			Metadata: req.Metadata,
			Data:     sDec,
		}, nil
	case map[string]interface{}:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error during casting json data, %s", err.Error())
		}
		return &Request{
			Metadata: req.Metadata,
			Data:     data,
		}, nil
	default:
		return nil, fmt.Errorf("invalid data format, %s", reflect.TypeOf(v))
	}
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
