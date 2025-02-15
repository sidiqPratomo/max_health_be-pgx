package repository

import (
	"context"
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type PrescriptionDrugRepository interface {
	PostOnePrescriptionDrug(ctx context.Context, prescriptionId int64, prescriptionDrug entity.PrescriptionDrug) error
	GetAllPrescriptionDrug(ctx context.Context, prescriptionId int64) ([]entity.PrescriptionDrug, error)
}

type prescriptionDrugRepositoryPostgres struct {
	db DBTX
}

func NewPrescriptionDrugRepositoryPostgres(db *sql.DB) prescriptionDrugRepositoryPostgres {
	return prescriptionDrugRepositoryPostgres{
		db: db,
	}
}

func (r *prescriptionDrugRepositoryPostgres) PostOnePrescriptionDrug(ctx context.Context, prescriptionId int64, prescriptionDrug entity.PrescriptionDrug) error {
	_, err := r.db.ExecContext(ctx, database.PostOnePrescriptionDrugQuery, prescriptionId, prescriptionDrug.Drug.Id, prescriptionDrug.Quantity, prescriptionDrug.Note)
	if err != nil {
		return err
	}

	return nil
}

func (r *prescriptionDrugRepositoryPostgres) GetAllPrescriptionDrug(ctx context.Context, prescriptionId int64) ([]entity.PrescriptionDrug, error) {
	rows, err := r.db.QueryContext(ctx, database.GetAllPrescriptionDrugQuery, prescriptionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptionDrugList []entity.PrescriptionDrug

	for rows.Next() {
		var prescriptionDrug entity.PrescriptionDrug

		err := rows.Scan(&prescriptionDrug.Id, &prescriptionDrug.Drug.Id, &prescriptionDrug.Drug.Name, &prescriptionDrug.Drug.Image, &prescriptionDrug.Drug.IsActive, &prescriptionDrug.Quantity, &prescriptionDrug.Note)
		if err != nil {
			return nil, err
		}

		prescriptionDrugList = append(prescriptionDrugList, prescriptionDrug)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return prescriptionDrugList, nil
}

func (r *prescriptionDrugRepositoryPostgres) GetPrescriptionDrugByCartItemId(ctx context.Context, cartItemId int64) error {
	return nil
}
