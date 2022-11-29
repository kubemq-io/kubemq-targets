package binding

import "github.com/kubemq-io/kubemq-targets/types"

type Response struct {
	Metadata map[string]string `json:"metadata"`
	Data     interface{}       `json:"data"`
	IsError  bool              `json:"is_error"`
	Error    string            `json:"error"`
}

func toResponse(r *types.Response) *Response {
	return &Response{
		Metadata: r.Metadata,
		Data:     types.ParseResponseBody(r.Data),
		IsError:  r.IsError,
		Error:    r.Error,
	}
}
