package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type PharmacyCourierRepository interface {
	CreateBulk(ctx context.Context, pharmacyId int64, courierIds []int64) error
	UpdateOneById(ctx context.Context, pharmacyCourier entity.PharmacyCourier) error
	DeleteBulkByPharmacyId(ctx context.Context, pharmacyId int64) error
}

type pharmacyCourierRepositoryPostgres struct {
	db DBTX
}

func NewPharmacyCourierRepositoryPostgres(db *sql.DB) pharmacyCourierRepositoryPostgres {
	return pharmacyCourierRepositoryPostgres{
		db: db,
	}
}

func (r *pharmacyCourierRepositoryPostgres) CreateBulk(ctx context.Context, pharmacyId int64, courierIds []int64) error {
	query := database.CreatePharmacyCourier

	args := []interface{}{}
	for i, courierId := range courierIds {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) + `)`
		args = append(args, pharmacyId)
		args = append(args, courierId)
		args = append(args, false)
		if i != len(courierIds)-1 {
			query += `,`
		}
	}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *pharmacyCourierRepositoryPostgres) UpdateOneById(ctx context.Context, pharmacyCourier entity.PharmacyCourier) error {
	_, err := r.db.ExecContext(ctx, database.UpdateOnePharmacyCourier, pharmacyCourier.IsActive, pharmacyCourier.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *pharmacyCourierRepositoryPostgres) DeleteBulkByPharmacyId(ctx context.Context, pharmacyId int64) error {
	_, err := r.db.ExecContext(ctx, database.DeleteBulkPharmacyCourierByPharmacyId, pharmacyId)
	if err != nil {
		return err
	}
	return nil
}
