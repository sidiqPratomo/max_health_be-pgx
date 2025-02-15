package repository

import (
	"context"
	// "database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewRefreshTokenRepositoryPostgres(db *pgxpool.Pool) refreshTokenRepositoryPostgres {
	return refreshTokenRepositoryPostgres{
		db: db,
	}
}

func (r *refreshTokenRepositoryPostgres) InvalidateCodes(ctx context.Context, accountId int64) error {
	_, err := r.db.Exec(ctx, database.InvalidateRefreshTokensQuery, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepositoryPostgres) PostOneCode(ctx context.Context, accountId int64, code string) error {
	_, err := r.db.Exec(ctx, database.PostOneRefreshTokenQuery, accountId, code)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepositoryPostgres) FindOneCode(ctx context.Context, refreshToken string) (string, error) {
	err := r.db.QueryRow(ctx, database.FindOneRefreshTokenQuery, refreshToken).Scan(&refreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return refreshToken, nil
}
