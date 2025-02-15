package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type PharmacyOperationalRepository interface {
	CreateBulk(ctx context.Context, pharmacyId int64, days []string) error
	UpdateOneById(ctx context.Context, pharmacyOperational entity.PharmacyOperational) error
	DeleteBulkByPharmacyId(ctx context.Context, pharmacyId int64) error
}

type pharmacyOperationalRepositoryPostgres struct {
	db DBTX
}

func NewPharmacyOperationalRepositoryPostgres(db *sql.DB) pharmacyOperationalRepositoryPostgres {
	return pharmacyOperationalRepositoryPostgres{
		db: db,
	}
}

func (r *pharmacyOperationalRepositoryPostgres) CreateBulk(ctx context.Context, pharmacyId int64, days []string) error {
	query := database.CreatePharmacyOperational

	args := []interface{}{}
	for i, day := range days {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `)`
		args = append(args, pharmacyId)
		args = append(args, day)
		if i != len(days)-1 {
			query += `,`
		}
	}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *pharmacyOperationalRepositoryPostgres) UpdateOneById(ctx context.Context, pharmacyOperational entity.PharmacyOperational) error {
	_, err := r.db.ExecContext(ctx, database.UpdateOnePharmacyOperational, pharmacyOperational.OperationalDay, pharmacyOperational.OpenHour, pharmacyOperational.CloseHour, pharmacyOperational.IsOpen, pharmacyOperational.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *pharmacyOperationalRepositoryPostgres) DeleteBulkByPharmacyId(ctx context.Context, pharmacyId int64) error {
	_, err := r.db.ExecContext(ctx, database.DeleteBulkPharmacyOperationalByPharmacyId, pharmacyId)
	if err != nil {
		return err
	}
	return nil
}
