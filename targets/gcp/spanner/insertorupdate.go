package spanner

import "fmt"

type InsertOrUpdate struct {
	TableName   string        `json:"table_name"`
	ColumnName  []string      `json:"column_names"`
	ColumnValue []interface{} `json:"column_values"`
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
		return fmt.Errorf("please verify column_values and column_names are the same size")
	}
	return nil
}


