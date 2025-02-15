package repository

import (
	"context"
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type ResetPasswordTokenRepository interface {
	InvalidateTokens(ctx context.Context, accountId int64) error
	PostOneToken(ctx context.Context, accountId int64, token string) error
	FindOneByAccountIdAndToken(ctx context.Context, resetPasswordToken *entity.ResetPasswordToken) (*entity.ResetPasswordToken, error)
	UpdateExpiredAtOne(ctx context.Context, resetPasswordToken *entity.ResetPasswordToken) error
}

type resetPasswordTokenRepositoryPostgres struct {
	db DBTX
}

func NewResetPasswordTokenRepositoryPostgres(db *sql.DB) resetPasswordTokenRepositoryPostgres {
	return resetPasswordTokenRepositoryPostgres{
		db: db,
	}
}

func (r *resetPasswordTokenRepositoryPostgres) InvalidateTokens(ctx context.Context, accountId int64) error {
	_, err := r.db.ExecContext(ctx, database.InvalidateTokensQuery, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *resetPasswordTokenRepositoryPostgres) PostOneToken(ctx context.Context, accountId int64, token string) error {
	_, err := r.db.ExecContext(ctx, database.PostOneResetPasswordTokenQuery, accountId, token)
	if err != nil {
		return err
	}

	return nil
}

func (r *resetPasswordTokenRepositoryPostgres) FindOneByAccountIdAndToken(ctx context.Context, resetPasswordToken *entity.ResetPasswordToken) (*entity.ResetPasswordToken, error) {
	query := database.SelectOneResetPasswordTokenByAccountIdAndTokenQuery
	if err := r.db.QueryRowContext(ctx, query, resetPasswordToken.AccountId, resetPasswordToken.Token).Scan(&resetPasswordToken.Id, &resetPasswordToken.ExpiredAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return resetPasswordToken, nil
}

func (r *resetPasswordTokenRepositoryPostgres) UpdateExpiredAtOne(ctx context.Context, resetPasswordToken *entity.ResetPasswordToken) error {
	query := database.UpdateExpiredAtOneResetPasswordTokenQuery

	_, err := r.db.ExecContext(ctx, query, resetPasswordToken.Id)
	if err != nil {
		return err
	}

	return nil
}
