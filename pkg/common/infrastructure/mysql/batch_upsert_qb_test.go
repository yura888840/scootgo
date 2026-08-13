package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var time1 = time.Date(2025, 2, 15, 16, 52, 20, 0, time.UTC)
var time2 = time.Date(2024, 9, 10, 16, 52, 20, 0, time.UTC)

func TestGenerateBatchUpsertQuery(t *testing.T) {
	type Entity struct {
		ID    int       `db:"id" pk:"true"`
		Name  string    `db:"name"`
		Value string    `db:"value"`
		Date  time.Time `db:"updated_at"`
	}

	type AnotherEntity struct {
		Key   int       `db:"key" pk:"true"`
		Data  string    `db:"data"`
		Extra string    `db:"extra"`
		Date  time.Time `db:"updated_at"`
	}

	tests := []struct {
		name          string
		tableName     string
		entities      any
		expectedQuery string
		expectedArgs  []any
		expectError   bool
	}{
		{
			name:      "Basic functionality",
			tableName: "your_table_name",
			entities: []Entity{
				{ID: 1, Name: "Entity1", Value: "Value1", Date: time1},
				{ID: 2, Name: "Entity2", Value: "Value2", Date: time2},
			},
			//nolint lll
			expectedQuery: "INSERT INTO your_table_name (id, name, value, updated_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name), value=VALUES(value), updated_at=VALUES(updated_at)",
			expectedArgs:  []any{1, "Entity1", "Value1", time1, 2, "Entity2", "Value2", time2},
			expectError:   false,
		},
		{
			name:          "Empty slice",
			tableName:     "your_table_name",
			entities:      []Entity{},
			expectedQuery: "",
			expectedArgs:  []any{},
			expectError:   true,
		},
		{
			name:      "Different struct",
			tableName: "another_table",
			entities: []AnotherEntity{
				{Key: 1, Data: "Data1", Extra: "Extra1", Date: time1},
				{Key: 2, Data: "Data2", Extra: "Extra2", Date: time2},
			},
			//nolint lll
			expectedQuery: "INSERT INTO another_table (key, data, extra, updated_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?) ON DUPLICATE KEY UPDATE data=VALUES(data), extra=VALUES(extra), updated_at=VALUES(updated_at)",
			expectedArgs:  []any{1, "Data1", "Extra1", time1, 2, "Data2", "Extra2", time2},
			expectError:   false,
		},
		{
			name:      "Field missing db tag",
			tableName: "your_table_name",
			entities: []struct {
				ID    int    `db:"id" pk:"true"`
				Name  string `db:"name"`
				Value string
			}{
				{ID: 1, Name: "Entity1", Value: "Value1"},
				{ID: 2, Name: "Entity2", Value: "Value2"},
			},
			expectedQuery: "INSERT INTO your_table_name (id, name) VALUES (?, ?), (?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name)",
			expectedArgs:  []any{1, "Entity1", 2, "Entity2"},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.entities)
			query, args, err := generateBatchUpsertQuery(tt.tableName, v)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedQuery, query)
				assert.Equal(t, tt.expectedArgs, args)
			}
		})
	}
}
