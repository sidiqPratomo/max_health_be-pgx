package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type DoctorSpecializationRepository interface {
	GetAllDoctorSpecialization(ctx context.Context) ([]entity.DoctorSpecialization, error)
}

type doctorSpecializationRepositoryPostgres struct {
	db DBTX
}

func NewDoctorSpecializationRepositoryPostgres(db *pgxpool.Pool) doctorSpecializationRepositoryPostgres {
	return doctorSpecializationRepositoryPostgres{
		db: db,
	}
}

func (r *doctorSpecializationRepositoryPostgres) GetAllDoctorSpecialization(ctx context.Context) ([]entity.DoctorSpecialization, error) {
	rows, err := r.db.Query(ctx, database.GetAllDoctorSpecializationQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specializationList []entity.DoctorSpecialization

	for rows.Next() {
		var specialization entity.DoctorSpecialization

		err := rows.Scan(&specialization.Id, &specialization.Name)
		if err != nil {
			return nil, err
		}

		specializationList = append(specializationList, specialization)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return specializationList, nil
}
