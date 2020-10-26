package spanner

import (
	"fmt"
	"strconv"
)

type InsertOrUpdate struct {
	TableName   string        `json:"table_name"`
	ColumnName  []string      `json:"column_names"`
	ColumnValue []interface{} `json:"column_values"`
	ColumnType  []string      `json:"column_type"`
}

var allowedTypes = []string{
	"INT64",
	"FLOAT64",
	"BOOL",
	"STRING",
	"BYTES",
	"DATE",
	"TIMESTAMP",
}

func (i InsertOrUpdate) validate() error {
	if len(i.ColumnName) == 0 {
		return fmt.Errorf("please verify column_names got atleast one value")
	}
	if len(i.ColumnValue) == 0 {
		return fmt.Errorf("please verify column_values got atleast one value")
	}
	if len(i.TableName) == 0 {
		return fmt.Errorf("please verify table_name is not empty")
	}
	if len(i.ColumnName) != len(i.ColumnValue) {
		return fmt.Errorf("please verify that column_values and column_names are the same size")
	}
	if len(i.ColumnValue) != len(i.ColumnType) {
		return fmt.Errorf("please verify that column_type and column_values are the same size")
	}
	for r := 0; r < len(i.ColumnType); r++ {
		switch i.ColumnType[r] {
		case "INT64":
			v, err := strconv.ParseInt(fmt.Sprintf("%v", i.ColumnValue[r]), 10, 64)
			if err != nil {
				return err
			}
			i.ColumnValue[r] = v
		case "FLOAT64":
			v, err := strconv.ParseFloat(fmt.Sprintf("%v", i.ColumnValue[r]), 32)
			if err != nil {
				return err
			}
			i.ColumnValue[r] = v
		case "BOOL":
			v, err := strconv.ParseBool(fmt.Sprintf("%v", i.ColumnValue[r]))
			if err != nil {
				return err
			}
			i.ColumnValue[r] = v
		default:
			_, found := find(allowedTypes, i.ColumnType[r])
			if !found {
				return fmt.Errorf(getValidValueTypes())
			}
		}
	}
	return nil
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func getValidValueTypes() string {
	s := "invalid method type, method type should be one of the following:"
	for _, k := range allowedTypes {
		s = fmt.Sprintf("%s :%s,", s, k)
	}
	return s
}
