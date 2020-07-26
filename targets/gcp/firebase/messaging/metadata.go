package messaging

import (
	"encoding/json"

	"firebase.google.com/go/messaging"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	message          *messaging.Message
	sendMulticast    bool
	multicastMessage *messaging.MulticastMessage
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{
		message: opts.defult,
	}

	m.sendMulticast = meta.ParseBool("sendMulticast", false)

	if m.sendMulticast {
		n := meta.ParseString("multicastMessage", "")
		if n != "" {
			err := json.Unmarshal([]byte(n), &m.multicastMessage)
			if err != nil {
				return m, err
			}
		}
	} else {

		n := meta.ParseString("message", "")
		if n != "" {
			err := json.Unmarshal([]byte(n), &m.message)
			if err != nil {
				return m, err
			}
		}
	}
	return m, nil
}
