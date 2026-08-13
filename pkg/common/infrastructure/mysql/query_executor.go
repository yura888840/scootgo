package mysql

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

const DefaultBatchSize = 1000

func BatchSave(ctx context.Context, db *sqlx.DB, tableName string, entities any, batchSize int) error {
	if batchSize < 1 {
		batchSize = 1
	}

	v := reflect.ValueOf(entities)
	if v.Kind() != reflect.Slice || v.Len() == 0 {
		return fmt.Errorf("entities must be a non-empty slice: %+v", entities)
	}

	for i := 0; i < v.Len(); i += batchSize {
		end := min(i+batchSize, v.Len())
		batch := v.Slice(i, end)

		query, args, err := generateBatchUpsertQuery(tableName, batch)
		if err != nil {
			return err
		}

		stmt, err := db.PreparexContext(ctx, query)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error while batch inserting into %s", tableName), err)
		}
	}

	return nil
}

func Update(ctx context.Context, db *sqlx.DB, query string, args []any) (uint, error) {
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return 0, fmt.Errorf("error while expanding IN clause: %w", err)
	}

	query = db.Rebind(query)
	truncatedQuery := truncateQueryForLogging(query)

	stmt, err := db.PreparexContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf(fmt.Sprintf("error while preparing update query: %s", truncatedQuery), err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf(fmt.Sprintf("error while executing update query: %s", truncatedQuery), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf(fmt.Sprintf("error while fetching rowsAffected from update query: %s", truncatedQuery), err)
	}

	return uint(rowsAffected), nil
}

func truncateQueryForLogging(query string) string {
	truncatedQuery := query
	if len(query) > 100 {
		truncatedQuery = query[:100] + "..."
	}

	return truncatedQuery
}
