package services

import (
	"bigdata-api/internal/repository"
	"context"
)

func GetUserByAPIKey(ctx context.Context, keyHash string) (*repository.User, error) {
	return repository.GetUserByAPIHash(ctx, keyHash)
}
