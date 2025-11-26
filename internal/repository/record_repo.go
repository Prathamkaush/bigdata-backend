package repository

import (
	"bigdata-api/internal/database"
	"context"
	"time"

	"github.com/google/uuid"
)

func SearchRecords(ctx context.Context, sql string, args []interface{}) ([]map[string]interface{}, error) {

	rows, err := database.ClickHouse.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := rows.Columns()
	colTypes := rows.ColumnTypes()

	results := make([]map[string]interface{}, 0)

	for rows.Next() {

		values := make([]interface{}, len(cols))

		for i, ct := range colTypes {

			switch ct.DatabaseTypeName() {

			case "UUID":
				values[i] = new(uuid.UUID)

			case "Date", "Date32", "DateTime", "DateTime64":
				values[i] = new(time.Time)

			case "String":
				values[i] = new(string)

			case "Int8", "Int16", "Int32", "Int64":
				values[i] = new(int64)

			// 🔥 FIXED: UInt16 must scan into *uint16
			case "UInt8", "UInt16":
				values[i] = new(uint16)

			case "UInt32":
				values[i] = new(uint32)

			case "UInt64":
				values[i] = new(uint64)

			case "Float32", "Float64":
				values[i] = new(float64)

			default:
				var v interface{}
				values[i] = &v
			}
		}

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			rowMap[col] = deref(values[i])
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func deref(v interface{}) interface{} {
	switch t := v.(type) {

	case *uuid.UUID:
		return t.String()

	case *time.Time:
		return t.Format(time.RFC3339)

	case *string:
		return *t

	case *int64:
		return *t

	case *uint16:
		return *t

	case *uint32:
		return *t

	case *uint64:
		return *t

	case *float64:
		return *t

	case *interface{}:
		return *t

	default:
		return t
	}
}

func CountRecords(ctx context.Context, sql string, args []interface{}) (uint64, error) {
	row := database.ClickHouse.QueryRow(ctx, sql, args...)
	var count uint64
	err := row.Scan(&count)
	return count, err
}
