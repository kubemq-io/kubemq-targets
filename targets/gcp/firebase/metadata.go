package firestore

import (
	"fmt"
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
}

var retrieveMap = map[string]string{
	"by_uid":   "by_uid",
	"by_email": "by_email",
	"by_phone": "by_phone",
}

type metadata struct {
	method     string
	key        string
	item       string
	tokenID    string
	retrieveBy string
	uid        string
	email      string
	phone      string
}

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

	return m, nil
}
