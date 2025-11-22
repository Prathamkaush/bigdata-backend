package services

import (
	"bigdata-api/internal/repository"
	"context"
)

func SearchRecords(ctx context.Context, sql string, args []interface{}) ([]map[string]interface{}, error) {
	return repository.SearchRecords(ctx, sql, args)
}

func CountRecords(ctx context.Context, sql string, args []interface{}) (uint64, error) {
	return repository.CountRecords(ctx, sql, args)
}
