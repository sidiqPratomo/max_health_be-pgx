package repository

import (
	"context"
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/database"
)

type RefreshTokenRepository interface {
	InvalidateCodes(ctx context.Context, accountId int64) error
	PostOneCode(ctx context.Context, accountId int64, code string) error
	FindOneCode(ctx context.Context, refreshToken string) (string, error)
}

type refreshTokenRepositoryPostgres struct {
	db DBTX
}

func NewRefreshTokenRepositoryPostgres(db *sql.DB) refreshTokenRepositoryPostgres {
	return refreshTokenRepositoryPostgres{
		db: db,
	}
}

func (r *refreshTokenRepositoryPostgres) InvalidateCodes(ctx context.Context, accountId int64) error {
	_, err := r.db.ExecContext(ctx, database.InvalidateRefreshTokensQuery, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepositoryPostgres) PostOneCode(ctx context.Context, accountId int64, code string) error {
	_, err := r.db.ExecContext(ctx, database.PostOneRefreshTokenQuery, accountId, code)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepositoryPostgres) FindOneCode(ctx context.Context, refreshToken string) (string, error) {
	err := r.db.QueryRowContext(ctx, database.FindOneRefreshTokenQuery, refreshToken).Scan(&refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return refreshToken, nil
}
