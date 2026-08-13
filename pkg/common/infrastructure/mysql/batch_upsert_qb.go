package mysql

import (
	"fmt"
	"reflect"
	"strings"
)

func generateBatchUpsertQuery(tableName string, batch reflect.Value) (string, []any, error) {
	if batch.Len() == 0 {
		return "", nil, fmt.Errorf("batch must be a non-empty slice: %+v", batch)
	}
	valueArgs, columns, valueStrings, primaryKeys := prepInsertQueryComponents(batch)

	updateColumns := prepUpdateQueryComponents(columns, primaryKeys)

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(valueStrings, ", "),
		strings.Join(updateColumns, ", "),
	)

	return query, valueArgs, nil
}

func prepInsertQueryComponents(v reflect.Value) (valueArgs []any, columns, valueStrings []string, primaryKeys map[string]bool) {
	primaryKeys = make(map[string]bool)
	for i := 0; i < v.Len(); i++ {
		entity := v.Index(i)

		var valuePlaceholders []string
		for j := 0; j < entity.NumField(); j++ {
			field := entity.Type().Field(j)
			columnName := field.Tag.Get("db")
			if columnName == "" {
				continue
			}
			isPrimaryKey := field.Tag.Get("pk")
			if i == 0 {
				columns = append(columns, columnName)
				if isPrimaryKey == "true" {
					primaryKeys[columnName] = true
				}
			}
			valuePlaceholders = append(valuePlaceholders, "?")
			valueArgs = append(valueArgs, entity.Field(j).Interface())
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(valuePlaceholders, ", ")))
	}

	return valueArgs, columns, valueStrings, primaryKeys
}

func prepUpdateQueryComponents(columns []string, primaryKeys map[string]bool) []string {
	updateColumns := make([]string, 0, len(columns))
	for _, col := range columns {
		if _, ok := primaryKeys[col]; !ok { // Exclude primary key columns
			updateColumns = append(updateColumns, fmt.Sprintf("%s=VALUES(%s)", col, col))
		}
	}

	return updateColumns
}
