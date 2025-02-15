package repository

import (
	"context"
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type VerificationCodeRepository interface {
	InvalidateCodes(ctx context.Context, accountId int64) error
	PostOneCode(ctx context.Context, accountId int64, code string) error
	FindOneByAccountIdAndCode(ctx context.Context, verificationCode *entity.VerificationCode) (*entity.VerificationCode, error)
	UpdateExpiredAtOne(ctx context.Context, verificationCode *entity.VerificationCode) error
}

type verificationCodeRepositoryPostgres struct {
	db DBTX
}

func NewVerificationCodeRepositoryPostgres(db *sql.DB) verificationCodeRepositoryPostgres {
	return verificationCodeRepositoryPostgres{
		db: db,
	}
}

func (r *verificationCodeRepositoryPostgres) InvalidateCodes(ctx context.Context, accountId int64) error {
	_, err := r.db.ExecContext(ctx, database.InvalidateCodesQuery, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *verificationCodeRepositoryPostgres) PostOneCode(ctx context.Context, accountId int64, code string) error {
	_, err := r.db.ExecContext(ctx, database.PostOneCodeQuery, accountId, code)
	if err != nil {
		return err
	}

	return nil
}

func (r *verificationCodeRepositoryPostgres) FindOneByAccountIdAndCode(ctx context.Context, verificationCode *entity.VerificationCode) (*entity.VerificationCode, error) {
	query := database.SelectOneVerificationCodeByCodeQuery
	if err := r.db.QueryRowContext(ctx, query, verificationCode.AccountId, verificationCode.Code).Scan(&verificationCode.Id, &verificationCode.ExpiredAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return verificationCode, nil
}

func (r *verificationCodeRepositoryPostgres) UpdateExpiredAtOne(ctx context.Context, verificationCode *entity.VerificationCode) error {
	query := database.UpdateExpiredAtOneVerificationCodeQuery

	_, err := r.db.ExecContext(ctx, query, verificationCode.Id)
	if err != nil {
		return err
	}

	return nil
}
