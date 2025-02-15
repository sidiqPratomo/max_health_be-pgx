package repository

import (
	"context"
	// "database/pgx"
	"math"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type DrugRepository interface {
	GetDrugByName(ctx context.Context, drugName string) (*entity.Drug, error)
	GetDrugIdByName(ctx context.Context, drugName string) (*int64, error)
	GetOneActiveDrugById(ctx context.Context, drugId int64) (*entity.Drug, error)
	GetOneDrugById(ctx context.Context, drugId int64) (*entity.Drug, error)
	GetDrugById(ctx context.Context, drugId int64) (*entity.DrugDetail, error)
	GetAllDrugs(ctx context.Context, validatedGetProductAdminQuery util.ValidatedGetDrugAdminQuery) ([]entity.Drug, *entity.PageInfo, error)
	UpdateOneDrug(ctx context.Context, drug entity.Drug) error
	CreateOneDrug(ctx context.Context, drug entity.Drug) error
	DeleteOneDrug(ctx context.Context, drugId int64) error
	GetDrugsByPharmacyId(ctx context.Context, pharmacyId int64, Limit string, offset int, search string) ([]entity.PharmacyDrugByPharmacyId, *entity.PageInfo, error)
}

type drugRepositoryPostgres struct {
	db DBTX
}

func NewDrugRepositoryPostgres(db *pgxpool.Pool) drugRepositoryPostgres {
	return drugRepositoryPostgres{
		db: db,
	}
}

func (r *drugRepositoryPostgres) GetOneActiveDrugById(ctx context.Context, drugId int64) (*entity.Drug, error) {
	var drug entity.Drug

	err := r.db.QueryRow(ctx, database.GetOneActiveDrugByIdQuery, drugId).Scan(
		&drug.Id,
		&drug.Name,
		&drug.GenericName,
		&drug.Content,
		&drug.Manufacture,
		&drug.Description,
		&drug.Classification.Id,
		&drug.Classification.Name,
		&drug.Form.Id,
		&drug.Form.Name,
		&drug.UnitInPack,
		&drug.SellingUnit,
		&drug.Weight,
		&drug.Height,
		&drug.Length,
		&drug.Width,
		&drug.Image,
		&drug.Category.Id,
		&drug.Category.Name,
		&drug.Category.Url,
		&drug.IsPrescriptionRequired,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &drug, nil
}

func (r *drugRepositoryPostgres) GetDrugById(ctx context.Context, drugId int64) (*entity.DrugDetail, error){
	drug := entity.DrugDetail{}
	err := r.db.QueryRow(ctx, database.GetDrugById, drugId).Scan(
		&drug.Id, &drug.Name, &drug.GenericName, &drug.Content, &drug.Manufacture, &drug.Description, 
		&drug.ClassificationId, &drug.FormId, &drug.UnitInPack, &drug.SellingUnit, 
		&drug.Weight, &drug.Height, &drug.Length, &drug.Width, &drug.Image,&drug.DrugCategoryId, &drug.IsPrescriptionRequired, &drug.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &drug, nil
}
		
func (r *drugRepositoryPostgres) GetDrugByName(ctx context.Context, drugName string) (*entity.Drug, error) {
	var drug entity.Drug

	err := r.db.QueryRow(ctx, database.GetDrugByNameQuery, drugName).Scan(&drug.Id, &drug.Name, &drug.GenericName, &drug.Content, &drug.Manufacture)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &drug, nil
}

func (r *drugRepositoryPostgres) UpdateOneDrug(ctx context.Context, drug entity.Drug) error {
	_, err := r.db.Exec(ctx, database.UpdateOneDrugQuery,
		drug.Id,
		drug.Name,
		drug.GenericName,
		drug.Content,
		drug.Manufacture,
		drug.Description,
		drug.Classification.Id,
		drug.Form.Id,
		drug.Category.Id,
		drug.UnitInPack,
		drug.SellingUnit,
		drug.Weight,
		drug.Height,
		drug.Length,
		drug.Width,
		drug.Image,
		drug.IsActive,
		drug.IsPrescriptionRequired,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *drugRepositoryPostgres) GetOneDrugById(ctx context.Context, drugId int64) (*entity.Drug, error) {
	var drug entity.Drug

	err := r.db.QueryRow(ctx, database.GetOneDrugByIdQuery, drugId).Scan(
		&drug.Id,
		&drug.Name,
		&drug.GenericName,
		&drug.Content,
		&drug.Manufacture,
		&drug.Description,
		&drug.Classification.Id,
		&drug.Classification.Name,
		&drug.Form.Id,
		&drug.Form.Name,
		&drug.UnitInPack,
		&drug.SellingUnit,
		&drug.Weight,
		&drug.Height,
		&drug.Length,
		&drug.Width,
		&drug.Image,
		&drug.Category.Id,
		&drug.Category.Name,
		&drug.Category.Url,
		&drug.IsActive,
		&drug.IsPrescriptionRequired,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &drug, nil
}

func (r *drugRepositoryPostgres) GetAllDrugs(ctx context.Context, validatedGetProductAdminQuery util.ValidatedGetDrugAdminQuery) ([]entity.Drug, *entity.PageInfo, error) {
	query := database.GetAllDrugsByIdQuery
	args := []interface{}{}

	if validatedGetProductAdminQuery.Search != nil {
		query += `$` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+*validatedGetProductAdminQuery.Search+"%")
	} else {
		query += `'%%'`
	}

	query += ` LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, validatedGetProductAdminQuery.Limit)
	query += ` OFFSET $` + strconv.Itoa(len(args)+1)
	args = append(args, (validatedGetProductAdminQuery.Limit * (validatedGetProductAdminQuery.Page - 1)))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	drugs := []entity.Drug{}
	var pageInfo entity.PageInfo

	for rows.Next() {
		var drug entity.Drug
		err := rows.Scan(
			&drug.Id,
			&drug.Name,
			&drug.GenericName,
			&drug.Content,
			&drug.Manufacture,
			&drug.Description,
			&drug.Classification.Id,
			&drug.Classification.Name,
			&drug.Form.Id,
			&drug.Form.Name,
			&drug.UnitInPack,
			&drug.SellingUnit,
			&drug.Weight,
			&drug.Height,
			&drug.Length,
			&drug.Width,
			&drug.Image,
			&drug.Category.Id,
			&drug.Category.Name,
			&drug.Category.Url,
			&drug.IsActive,
			&drug.IsPrescriptionRequired,
			&pageInfo.ItemCount,
		)
		if err != nil {
			return nil, nil, err
		}
		drugs = append(drugs, drug)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	pageInfo.Page = validatedGetProductAdminQuery.Page
	pageInfo.PageCount = pageInfo.ItemCount / validatedGetProductAdminQuery.Limit
	if pageInfo.ItemCount%validatedGetProductAdminQuery.Limit != 0 {
		pageInfo.PageCount += 1
	}

	return drugs, &pageInfo, nil
}

func (r *drugRepositoryPostgres) GetDrugIdByName(ctx context.Context, drugName string) (*int64, error) {
	var drugId int64

	err := r.db.QueryRow(ctx, database.GetDrugIdByNameQuery, drugName).Scan(&drugId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &drugId, nil
}

func (r *drugRepositoryPostgres) CreateOneDrug(ctx context.Context, drug entity.Drug) error {
	_, err := r.db.Exec(ctx, database.CreateOneDrugQuery,
		drug.Name,
		drug.GenericName,
		drug.Content,
		drug.Manufacture,
		drug.Description,
		drug.Classification.Id,
		drug.Form.Id,
		drug.Category.Id,
		drug.UnitInPack,
		drug.SellingUnit,
		drug.Weight,
		drug.Height,
		drug.Length,
		drug.Width,
		drug.Image,
		drug.IsPrescriptionRequired,
		drug.IsActive,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *drugRepositoryPostgres) DeleteOneDrug(ctx context.Context, drugId int64) error {
	_, err := r.db.Exec(ctx, database.DeleteOneDrugQuery, drugId)
	if err != nil {
		return err
	}

	return nil
}

func (r *drugRepositoryPostgres) GetDrugsByPharmacyId(ctx context.Context, pharmacyId int64, Limit string, offset int, search string) ([]entity.PharmacyDrugByPharmacyId, *entity.PageInfo, error) {
	query := database.GetDrugsByPharmacyId
	drugs := []entity.PharmacyDrugByPharmacyId{}
	pageInfo := &entity.PageInfo{}

	intLimit, err := strconv.Atoi(Limit)
	if err != nil {
		return nil, nil, err
	}

	rows, err := r.db.Query(ctx, query, pharmacyId, intLimit, offset, search)
	if err != nil {
		return nil,nil, err
	}
	defer rows.Close()

	for rows.Next() {
		drug := entity.PharmacyDrugByPharmacyId{}
		err := rows.Scan(
			&drug.Id,
			&drug.Price,
			&drug.Stock,
			&drug.Drug.Id,
			&drug.Drug.Name,
			&drug.Drug.GenericName,
			&drug.Drug.Content,
			&drug.Drug.Manufacture,
			&drug.Drug.Description,
			&drug.Drug.Classification.Id,
			&drug.Drug.Classification.Name,
			&drug.Drug.Form.Id,
			&drug.Drug.Form.Name,
			&drug.Drug.UnitInPack,
			&drug.Drug.SellingUnit,
			&drug.Drug.Weight,
			&drug.Drug.Height,
			&drug.Drug.Length,
			&drug.Drug.Width,
			&drug.Drug.Image,
			&drug.Drug.Category.Id,
			&drug.Drug.Category.Name,
			&drug.Drug.Category.Url,
			&drug.Drug.IsActive,
			&drug.Drug.IsPrescriptionRequired,
		)
		if err != nil {
			return nil, nil, err
		}
		drugs = append(drugs, drug)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM pharmacy_drugs pd
		JOIN drugs d
		ON pd.drug_id = d.drug_id
		WHERE pd.pharmacy_id= $1 and pd.deleted_at IS NULL AND d.deleted_at IS NULL and d.drug_name ILIKE '%' || $2 || '%'
	`
	countRow := r.db.QueryRow(ctx, countQuery, pharmacyId, search)
	if err := countRow.Scan(&pageInfo.ItemCount); err != nil {
		return nil, nil, err
	}

	pageInfo.PageCount = int(math.Ceil(float64(pageInfo.ItemCount) / float64(intLimit)))
	pageInfo.Page = int(math.Ceil(float64(offset+1) / float64(intLimit)))

	return drugs, pageInfo, nil
}
