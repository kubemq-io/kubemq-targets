package bigquery

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cloud.google.com/go/bigquery"
)

type insertRecords struct {
	records []record
}

func newInsertRecord(payload []byte) (*insertRecords, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("no insert payload found")
	}
	ir := &insertRecords{
		records: nil,
	}
	var recs []record
	err := json.Unmarshal(payload, &recs)
	if err != nil {
		return nil, err
	}
	for _, rec := range recs {
		if rec != nil {
			ir.records = append(ir.records, rec)
		}
	}
	return ir, err
}

func hash(data []byte) string {
	if data == nil {
		return ""
	}
	h := sha256.New()
	_, _ = h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

type record map[string]bigquery.Value

func (rec record) Save() (map[string]bigquery.Value, string, error) {
	data, err := json.Marshal(&rec)
	if err != nil {
		return nil, "", err
	}
	st := hash(data)
	return rec, st, nil
}
