package services

import (
	"bigdata-api/internal/repository"
	"context"
)

func DeductUserCredits(userID int64, creditsUsed int64) error {
	return repository.DeductCredits(
		context.Background(),
		int(userID),
		int(creditsUsed),
	)
}
