package repository

import (
	"bigdata-api/internal/database"
	"context"
	"reflect"
	"time"
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
			kind := ct.ScanType().Kind()

			switch kind {

			case reflect.String:
				var v string
				values[i] = &v

			case reflect.Int, reflect.Int32, reflect.Int64:
				var v int64
				values[i] = &v

			case reflect.Uint, reflect.Uint32, reflect.Uint64:
				var v uint64
				values[i] = &v

			case reflect.Float32, reflect.Float64:
				var v float64
				values[i] = &v

			default:
				// Special case for ClickHouse DateTime, Date, Date32
				if ct.ScanType() == reflect.TypeOf(time.Time{}) {
					var v time.Time
					values[i] = &v
				} else {
					// Fallback: string
					var v string
					values[i] = &v
				}
			}
		}

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{}, len(cols))
		for idx, col := range cols {
			rowMap[col] = deref(values[idx])
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func CountRecords(ctx context.Context, sql string, args []interface{}) (uint64, error) {
	row := database.ClickHouse.QueryRow(ctx, sql, args...)

	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func deref(ptr interface{}) interface{} {
	switch v := ptr.(type) {
	case *string:
		return *v
	case *int32:
		return *v
	case *int64:
		return *v
	case *uint32:
		return *v
	case *uint64:
		return *v
	case *float32:
		return *v
	case *float64:
		return *v
	case *time.Time:
		return v.Format(time.RFC3339) // return as ISO string
	default:
		return v
	}
}
