package repository

import (
	"context"
	// "database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type UserRepository interface {
	PostOneUser(ctx context.Context, accountId int) error
	FindUserByAccountId(ctx context.Context, accountId int64) (*entity.User, error)
	UpdateDataOne(ctx context.Context, user *entity.DetailedUser) error
}

type userRepositoryPostgres struct {
	db DBTX
}

func NewUserRepositoryPostgres(db *pgxpool.Pool) userRepositoryPostgres {
	return userRepositoryPostgres{
		db: db,
	}
}

func (r *userRepositoryPostgres) PostOneUser(ctx context.Context, accountId int) error {
	_, err := r.db.Exec(ctx, database.PostOneUserQuery, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryPostgres) FindUserByAccountId(ctx context.Context, accountId int64) (*entity.User, error) {
	var user entity.User

	if err := r.db.QueryRow(ctx, database.FindUserByAccountIdQuery, accountId).Scan(&user.Id, &user.GenderId, &user.GenderName, &user.DateOfBirth); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r* userRepositoryPostgres) UpdateDataOne(ctx context.Context, user *entity.DetailedUser) error {
	_, err := r.db.Exec(ctx, database.UpdateDataOneUser, user.GenderId, user.DateOfBirth, user.Id)
	if err != nil {
		return err
	}

	return nil
}
