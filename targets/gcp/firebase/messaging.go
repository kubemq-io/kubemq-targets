package firebase

import (
	"context"
	"fmt"
	"strconv"

	"firebase.google.com/go/v4/messaging"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type messages struct {
	single    *messaging.Message
	multicast *messaging.MulticastMessage
}

func (c *Client) SendMessage(ctx context.Context, req *types.Request, opts options) (*types.Response, error) {
	m, err := parseMetadataMessages(req.Data, opts, SendMessage)
	if err != nil {
		return nil, err
	}

	r, err := c.messagingClient.Send(ctx, m.single)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().SetMetadataKeyValue("result", r), nil

}

func (c *Client) SendMessageBatch(ctx context.Context, req *types.Request, opts options) (*types.Response, error) {
	m, err := parseMetadataMessages(req.Data, opts, SendBatch)
	if err != nil {
		return nil, err
	}

	b, err := c.messagingClient.SendMulticast(ctx, m.multicast)
	if err != nil {
		return nil, err
	}
	r := types.NewResponse().
		SetMetadataKeyValue("SuccessCount", strconv.Itoa(b.SuccessCount)).
		SetMetadataKeyValue("FailureCount", strconv.Itoa(b.FailureCount))

	for _, res := range b.Responses {
		msg := fmt.Sprintf("MessageID:%s, Success:%t, Error:%s", res.MessageID, res.Success, res.Error.Error())
		r.SetMetadataKeyValue(fmt.Sprintf("mesage_%s", res.MessageID), msg)
	}
	return r, nil

}
