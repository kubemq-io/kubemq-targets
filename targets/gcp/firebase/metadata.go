package firebase

import (
	"encoding/json"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var methodsMap = map[string]string{
	//------------Token-----------//
	"custom_token": "custom_token",
	"verify_token": "verify_token",
	//------------User------------//
	"retrieve_user":         "retrieve_user",
	"create_user":           "create_user",
	"update_user":           "update_user",
	"delete_user":           "delete_user",
	"delete_multiple_users": "delete_multiple_users",
	"list_users":            "list_users",

	//------------DB-------------//
	"get_db":    "get_db",
	"update_db": "update_db",
	"delete_db": "delete_db",
	"set_db":    "set_db",
	//--------Messaging---------//
	"SendMessage":      "SendMessage",
	"SendMessageBatch": "SendMessageBatch",
}

var retrieveMap = map[string]string{
	"by_uid":   "by_uid",
	"by_email": "by_email",
	"by_phone": "by_phone",
}

type metadata struct {
	method     string
	tokenID    string
	retrieveBy string
	uid        string
	email      string
	phone      string

	refPath      string
	childRefPath string
}

const ( // iota is reset to 0
	Unassigned  = iota // c0 == 0
	SendMessage = iota // c1 == 1
	SendBatch   = iota // c2 == 2
)

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	if m.method == "retrieve_user" {
		m.retrieveBy, err = meta.ParseStringMap("retrieve_by", retrieveMap)
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing retrieve_by, %w", err)
		}

		switch m.retrieveBy {
		case "by_uid":
			m.uid, err = meta.MustParseString("uid")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing uid, %w", err)
			}
		case "by_email":
			m.email, err = meta.MustParseString("email")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing email, %w", err)
			}
		case "by_phone":
			m.phone, err = meta.MustParseString("phone")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing phone, %w", err)
			}
		}
	}

	if m.method == "update_user" || m.method == "delete_user" {
		m.uid, err = meta.MustParseString("uid")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing uid, %w", err)
		}
	}
	if m.method == "verify_token" || m.method == "custom_token" {
		m.tokenID, err = meta.MustParseString("token_id")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing token_id, %w", err)
		}
	}
	if m.method == "get_db" || m.method == "update_db" || m.method == "set_db" || m.method == "delete_db" {
		m.refPath, err = meta.MustParseString("ref_path")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing refPath, %w", err)
		}
		m.childRefPath = meta.ParseString("child_ref", "")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing child_ref, %w", err)
		}
	}

	return m, nil
}

func parseMetadataMessages(data []byte, opts options, metaDatatype int) (*messages, error) {

	// DeepCopy deepcopies a to b using json marshaling

	switch metaDatatype {
	//messaging single
	case 1:

		ms := messaging.Message{}
		DeepCopy(&opts.defaultMessaging.single, &ms)
		err := json.Unmarshal([]byte(data), &ms)
		if err != nil {
			return opts.defaultMessaging, err
		}
		return &messages{single: &ms}, nil
	case 2:
		mm := messaging.MulticastMessage{}
		DeepCopy(&opts.defaultMessaging.multicast, &mm)
		err := json.Unmarshal([]byte(data), &mm)
		if err != nil {
			return opts.defaultMessaging, err
		}
		return &messages{multicast: &mm}, nil
	}
	return opts.defaultMessaging, nil
}

// DeepCopy deepcopies a to b using json marshaling
func DeepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}
