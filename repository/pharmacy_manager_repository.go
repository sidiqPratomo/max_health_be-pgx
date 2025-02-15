package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type PharmacyManagerRepository interface {
	PostOne(ctx context.Context, accountId int64) error
	FindAll(ctx context.Context) ([]entity.PharmacyManager, error)
	FindOneById(ctx context.Context, pharmacyManagerId int64) (*entity.PharmacyManager, error)
	DeleteOneById(ctx context.Context, pharmacyManagerId int64) error
	FindOneByAccountId(ctx context.Context, accountId int64) (*entity.PharmacyManager, error)
	FindOneByPharmacyCourierId(ctx context.Context, pharmacyCourierId int64) (*entity.PharmacyManager, error)
}

type pharmacyManagerRepositoryPostgres struct {
	db DBTX
}

func NewpharmacyManagerRepositoryPostgres(db *pgxpool.Pool) pharmacyManagerRepositoryPostgres {
	return pharmacyManagerRepositoryPostgres{
		db: db,
	}
}

func (r *pharmacyManagerRepositoryPostgres) PostOne(ctx context.Context, accountId int64) error {
	_, err := r.db.Exec(ctx, database.PostOnePharmacyManagerQuery, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *pharmacyManagerRepositoryPostgres) FindAll(ctx context.Context) ([]entity.PharmacyManager, error) {
	query := database.FindAllPharmacyManagers

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pharmacyManagers := []entity.PharmacyManager{}
	for rows.Next() {
		var pharmacyManager entity.PharmacyManager
		err := rows.Scan(
			&pharmacyManager.Id,
			&pharmacyManager.Account.Id,
			&pharmacyManager.Account.Email,
			&pharmacyManager.Account.Name,
			&pharmacyManager.Account.ProfilePicture,
		)
		if err != nil {
			return nil, err
		}
		pharmacyManagers = append(pharmacyManagers, pharmacyManager)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pharmacyManagers, nil
}

func (r *pharmacyManagerRepositoryPostgres) FindOneById(ctx context.Context, pharmacyManagerId int64) (*entity.PharmacyManager, error) {
	var pharmacyManager entity.PharmacyManager

	if err := r.db.QueryRow(ctx, database.GetOnePharmacyManagerByIdQuery, pharmacyManagerId).Scan(&pharmacyManager.Id, &pharmacyManager.Account.Id, &pharmacyManager.Account.ProfilePicture); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &pharmacyManager, nil
}

func (r *pharmacyManagerRepositoryPostgres) DeleteOneById(ctx context.Context, pharmacyManagerId int64) error {
	_, err := r.db.Exec(ctx, database.DeleteOnePharmacyManagerByIdQuery, pharmacyManagerId)
	if err != nil {
		return err
	}
	return nil
}

func (r *pharmacyManagerRepositoryPostgres) FindOneByAccountId(ctx context.Context, accountId int64) (*entity.PharmacyManager, error) {
	var pharmacyManager entity.PharmacyManager

	if err := r.db.QueryRow(ctx, database.GetOnePharmacyManagerByAccountIdQuery, accountId).Scan(&pharmacyManager.Id, &pharmacyManager.Account.Id, &pharmacyManager.Account.ProfilePicture); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &pharmacyManager, nil
}

func (r *pharmacyManagerRepositoryPostgres) FindOneByPharmacyCourierId(ctx context.Context, pharmacyCourierId int64) (*entity.PharmacyManager, error) {
	var pharmacyManager entity.PharmacyManager

	if err := r.db.QueryRow(ctx, database.GetOnePharmacyManagerByPharmacyCourierIdQuery, pharmacyCourierId).Scan(&pharmacyManager.Id, &pharmacyManager.Account.Id, &pharmacyManager.Account.ProfilePicture); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &pharmacyManager, nil
}
