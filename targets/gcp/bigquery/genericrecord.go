package bigquery

import (
	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
)

type genericRecord map[string]bigquery.Value

func (rec genericRecord) Save() (map[string]bigquery.Value, string, error) {
	insertID := uuid.New().String()
	return rec, insertID, nil
}
