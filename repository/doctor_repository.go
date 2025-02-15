package repository

import (
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type DoctorRepository interface {
	PostOneDoctor(ctx context.Context, accountId int, specializationId int64, certificateName string) error
	UpdateDataOne(ctx context.Context, doctor *entity.DetailedDoctor) error
	FindSpecializationById(ctx context.Context, specializationId int64) (*string, error)
	GetAllDoctor(ctx context.Context, Sort []string, SortBy []string, Limit string, offset int, specialization_id string) ([]entity.Doctor, *entity.PageInfo, error)
	FindDoctorByAccountId(ctx context.Context, accountId int64) (*entity.Doctor, error)
	FindDoctorByDoctorId(ctx context.Context, doctorId int64) (*entity.DetailedDoctor, error)
	UpdateDoctorStatus(ctx context.Context, doctorAccountId int64, isOnline bool) error
	GetDoctorIsOnline(ctx context.Context, doctorAccountId int64) (*bool, error)
}

type doctorRepositoryPostgres struct {
	db DBTX
}

func NewDoctorRepositoryPostgres(db *sql.DB) doctorRepositoryPostgres {
	return doctorRepositoryPostgres{
		db: db,
	}
}

func (r *doctorRepositoryPostgres) PostOneDoctor(ctx context.Context, accountId int, specializationId int64, certificateName string) error {
	_, err := r.db.ExecContext(ctx, database.PostOneDoctorQuery, accountId, specializationId, certificateName)
	if err != nil {
		return err
	}

	return nil
}

func (r *doctorRepositoryPostgres) UpdateDataOne(ctx context.Context, doctor *entity.DetailedDoctor) error {
	_, err := r.db.ExecContext(ctx, database.UpdateOneDoctorQuery, doctor.FeePerPatient, doctor.Experience, doctor.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *doctorRepositoryPostgres) FindSpecializationById(ctx context.Context, specializationId int64) (*string, error) {
	var specializationName string
	err := r.db.QueryRowContext(ctx, database.FindSpecializationById, specializationId).Scan(&specializationName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &specializationName, nil
}

func (r *doctorRepositoryPostgres) GetAllDoctor(ctx context.Context, Sort []string, SortBy []string, Limit string, offset int, specialization_id string) ([]entity.Doctor, *entity.PageInfo, error) {
	doctors := []entity.Doctor{}
	pageInfo := &entity.PageInfo{}

	intLimit, err := strconv.Atoi(Limit)
	if err != nil {
		return nil, nil, err
	}

	var whereClause string
	if specialization_id == "" {
		whereClause = "WHERE d.deleted_at IS NULL AND a.verified_at IS NOT NULL"
	} else {
		whereClause = "where d.deleted_at IS NULL AND a.verified_at IS NOT NULL AND d.specialization_id = " + specialization_id
	}

	orderByClause := "ORDER BY d.is_online DESC"

	if Sort[0] != "" {
		for i := 0; i < len(Sort); i++ {
			if i == 0 {
				orderByClause += ","
			}

			orderByClause += " " + SortBy[i] + " " + Sort[i]

			if i < len(Sort)-1 {
				orderByClause += ","
			}
		}
	}

	queries := `
		SELECT d.doctor_id, d.account_id, d.certificate, d.fee_per_patient, d.is_online, d.experience, d.specialization_id 
		FROM doctors d
		JOIN accounts a ON a.account_id = d.account_id
		` + whereClause + `
		` + orderByClause + `
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, queries, Limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, err
		}

		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		doctor := entity.Doctor{}

		err := rows.Scan(&doctor.Id, &doctor.AccountId, &doctor.Certificate, &doctor.FeePerPatient, &doctor.IsOnline, &doctor.Experience, &doctor.SpecializationId)
		if err != nil {
			return nil, nil, err
		}
		doctors = append(doctors, doctor)
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM doctors d
		JOIN accounts a ON a.account_id = d.account_id
	` + whereClause
	countRow := r.db.QueryRowContext(ctx, countQuery)
	if err := countRow.Scan(&pageInfo.ItemCount); err != nil {
		return nil, nil, err
	}

	pageInfo.PageCount = int(math.Ceil(float64(pageInfo.ItemCount) / float64(intLimit)))
	pageInfo.Page = int(math.Ceil(float64(offset+1) / float64(intLimit)))

	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	return doctors, pageInfo, nil
}

func (r *doctorRepositoryPostgres) FindDoctorByAccountId(ctx context.Context, accountId int64) (*entity.Doctor, error) {
	var doctor entity.Doctor

	if err := r.db.QueryRowContext(ctx, database.FindDoctorByAccountIdQuery, accountId).Scan(&doctor.Id, &doctor.Experience, &doctor.SpecializationId,
		&doctor.SpecializationName, &doctor.FeePerPatient, &doctor.Certificate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepositoryPostgres) FindDoctorByDoctorId(ctx context.Context, doctorId int64) (*entity.DetailedDoctor, error) {
	var doctor entity.DetailedDoctor

	if err := r.db.QueryRowContext(ctx, database.FindDoctorByDoctorIdQuery, doctorId).Scan(&doctor.Id, &doctor.Email,
		&doctor.Name, &doctor.ProfilePicture, &doctor.Experience, &doctor.SpecializationId,
		&doctor.SpecializationName, &doctor.FeePerPatient); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepositoryPostgres) UpdateDoctorStatus(ctx context.Context, doctorAccountId int64, isOnline bool) error {
	_, err := r.db.ExecContext(ctx, database.UpdateDoctorStatusQuery, isOnline, doctorAccountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *doctorRepositoryPostgres) GetDoctorIsOnline(ctx context.Context, doctorAccountId int64) (*bool, error) {
	var isOnline bool
	err := r.db.QueryRowContext(ctx, database.GetDoctorIsOnlineQuery, doctorAccountId).Scan(&isOnline)
	if err != nil {
		return nil, err
	}

	return &isOnline, nil
}
