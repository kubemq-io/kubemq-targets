package spanner

import (
	"cloud.google.com/go/spanner"
	"encoding/base64"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"time"
)

type Column struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Row struct {
	Columns []*Column `json:"columns"`
}

func extractDataByType(r *spanner.Row) (*Row, error) {
	row:= &Row{}
	for i := 0; i < r.Size(); i++ {
		column := &Column{}
		var col spanner.GenericColumnValue
		err := r.Column(i,&col)
		if err != nil{
			return row, err
		}
		column.Name = r.ColumnName(i)
		switch col.Type.Code {
		case sppb.TypeCode_INT64:
			var v spanner.NullInt64
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			column.Value = v.Int64
		case sppb.TypeCode_FLOAT64:
			var v spanner.NullFloat64
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			column.Value = v.Float64
		case sppb.TypeCode_STRING:
			var v spanner.NullString
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			column.Value = v.StringVal
		case sppb.TypeCode_BYTES:
			var v spanner.NullString
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			if v.IsNull() {
				column.Value = []byte(nil)
			} else {
				b, err := base64.StdEncoding.DecodeString(v.StringVal)
				if err != nil {
					return row, err
				}
				column.Value = b
			}
		case sppb.TypeCode_BOOL:
			var v spanner.NullBool
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			column.Value = v.Bool
		case sppb.TypeCode_DATE:
			var v spanner.NullDate
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			if v.IsNull() {
				column.Value = v.Date // typed nil
			} else {
				column.Value = v.Date.In(time.Local)
			}
		case sppb.TypeCode_TIMESTAMP:
			var v spanner.NullTime
			if err := col.Decode(&v); err != nil {
				return row, err
			}
			column.Value = v.Time

		}
		row.Columns = append(row.Columns, column)

	}

	return row, nil
}
