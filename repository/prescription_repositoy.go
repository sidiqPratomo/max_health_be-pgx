package repository

import (
	"context"
	// "database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type PrescriptionRepository interface {
	CreateOnePrescription(ctx context.Context, userAccountId, doctorAccountId int64) (*int64, error)
	GetPrescriptionById(ctx context.Context, prescriptionId int64) (*entity.Prescription, error)
	SetPrescriptionRedeemedNow(ctx context.Context, prescriptionId int64) error
	GetPrescriptionListByUserAccountId(ctx context.Context, accountId int64, limit, offset int) ([]entity.Prescription, error)
	GetPrescriptionListByUserAccountIdTotalItem(ctx context.Context, accountId int64) (int, error)
	SetPrescriptionOrderedAtNow(ctx context.Context, prescriptionId int64) error
}

type prescriptionRepositoryPostgres struct {
	db DBTX
}

func NewPrescriptionRepositoryPostgres(db *pgxpool.Pool) prescriptionRepositoryPostgres {
	return prescriptionRepositoryPostgres{
		db: db,
	}
}

func (r *prescriptionRepositoryPostgres) CreateOnePrescription(ctx context.Context, userAccountId, doctorAccountId int64) (*int64, error) {
	var prescriptionId int64

	err := r.db.QueryRow(ctx, database.CreateOnePrescriptionQuery, userAccountId, doctorAccountId).Scan(&prescriptionId)
	if err != nil {
		return nil, err
	}

	return &prescriptionId, nil
}

func (r *prescriptionRepositoryPostgres) GetPrescriptionById(ctx context.Context, prescriptionId int64) (*entity.Prescription, error) {
	var prescription entity.Prescription

	err := r.db.QueryRow(ctx, database.GetPrescriptionByIdQuery, prescriptionId).Scan(&prescription.UserAccountId, &prescription.DoctorAccountId, &prescription.RedeemedAt, &prescription.OrderedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	prescription.Id = &prescriptionId

	return &prescription, nil
}

func (r *prescriptionRepositoryPostgres) SetPrescriptionRedeemedNow(ctx context.Context, prescriptionId int64) error {
	_, err := r.db.Exec(ctx, database.SetPrescriptionRedeemedNowQuery, prescriptionId)
	if err != nil {
		return err
	}

	return nil
}

func (r *prescriptionRepositoryPostgres) GetPrescriptionListByUserAccountId(ctx context.Context, accountId int64, limit, offset int) ([]entity.Prescription, error) {
	rows, err := r.db.Query(ctx, database.GetPrescriptionListByUserAccountIdQuery, accountId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptionList []entity.Prescription

	for rows.Next() {
		var prescription entity.Prescription

		err = rows.Scan(
			&prescription.Id,
			&prescription.UserAccountId,
			&prescription.UserName,
			&prescription.DoctorAccountId,
			&prescription.DoctorName,
			&prescription.RedeemedAt,
			&prescription.OrderedAt,
			&prescription.CreatedAt,
		)
		if err != nil {
			return nil, nil
		}

		prescriptionList = append(prescriptionList, prescription)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil
	}

	return prescriptionList, nil
}

func (r *prescriptionRepositoryPostgres) GetPrescriptionListByUserAccountIdTotalItem(ctx context.Context, accountId int64) (int, error) {
	var totalPage int

	err := r.db.QueryRow(ctx, database.GetPrescriptionListByUserAccountIdTotalPageQuery, accountId).Scan(&totalPage)
	if err != nil {
		return totalPage, err
	}

	return totalPage, nil
}

func (r *prescriptionRepositoryPostgres) SetPrescriptionOrderedAtNow(ctx context.Context, prescriptionId int64) error {
	_, err := r.db.Exec(ctx, database.SetPrescriptionOrderedAtNowQuery, prescriptionId)
	if err != nil {
		return err
	}

	return nil
}
