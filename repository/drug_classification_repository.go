package repository

import (
	"context"
	// "database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type DrugClassificationRepository interface {
	GetAllDrugClassification(ctx context.Context) ([]entity.DrugClassification, error)
	FindOneById(ctx context.Context, id int64) (*entity.DrugClassification, error)
}

type drugClassificationRepositoryPostgres struct {
	db DBTX
}

func NewDrugClassificationRepositoryPostgres(db *pgxpool.Pool) drugClassificationRepositoryPostgres {
	return drugClassificationRepositoryPostgres{
		db: db,
	}
}

func (r *drugClassificationRepositoryPostgres) GetAllDrugClassification(ctx context.Context) ([]entity.DrugClassification, error) {
	rows, err := r.db.Query(ctx, database.GetAllDrugClassificationQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classificationList []entity.DrugClassification

	for rows.Next() {
		var drugClassification entity.DrugClassification

		err := rows.Scan(&drugClassification.Id, &drugClassification.Name)
		if err != nil {
			return nil, err
		}

		classificationList = append(classificationList, drugClassification)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classificationList, nil
}

func (r *drugClassificationRepositoryPostgres) FindOneById(ctx context.Context, id int64) (*entity.DrugClassification, error) {
	var classification entity.DrugClassification

	if err := r.db.QueryRow(ctx, database.GetOneDrugClassficationById, id).Scan(&classification.Id, &classification.Name); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &classification, nil
}
