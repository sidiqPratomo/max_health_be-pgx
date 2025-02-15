package repository

import (
	"context"
	// "database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type AddressRepository interface {
	FindAllProvinces(ctx context.Context) ([]entity.Province, error)
	FindOneProvinceByName(ctx context.Context, name string) (*int64, error)
	FindAllCitiesByProvinceCode(ctx context.Context, provinceCode string) ([]entity.City, error)
	FindOneCityByName(ctx context.Context, name string) (*int64, error)
	FindAllDistrictsByCityCode(ctx context.Context, cityCode string) ([]entity.District, error)
	FindOneDistrictByName(ctx context.Context, name string) (*int64, error)
	FindAllSubdistrictsByDistrictCode(ctx context.Context, districtCode string) ([]entity.Subdistrict, error)
	FindOneSubdistrictByName(ctx context.Context, name string) (*int64, error)
}

type addressRepositoryPostgres struct {
	db DBTX
}

func NewAddressRepositoryPostgres(db *pgxpool.Pool) addressRepositoryPostgres {
	return addressRepositoryPostgres{
		db: db,
	}
}

func (r *addressRepositoryPostgres) FindAllProvinces(ctx context.Context) ([]entity.Province, error) {
	query := database.FindAllProvinces

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	provinces := []entity.Province{}

	for rows.Next() {
		var province entity.Province

		err := rows.Scan(
			&province.Id,
			&province.Code,
			&province.Name,
		)
		if err != nil {
			return nil, err
		}

		provinces = append(provinces, province)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return provinces, nil
}

func (r *addressRepositoryPostgres) FindOneProvinceByName(ctx context.Context, name string) (*int64, error) {
	var provinceId int64

	if err := r.db.QueryRow(ctx, database.FindOneProvinceByName, name).Scan(&provinceId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &provinceId, nil
}

func (r *addressRepositoryPostgres) FindAllCitiesByProvinceCode(ctx context.Context, provinceCode string) ([]entity.City, error) {
	query := database.FindAllCitiesByProvinceCode

	rows, err := r.db.Query(ctx, query, provinceCode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cities := []entity.City{}

	for rows.Next() {
		var city entity.City

		err := rows.Scan(
			&city.Id,
			&city.Code,
			&city.ProvinceCode,
			&city.Name,
		)
		if err != nil {
			return nil, err
		}

		cities = append(cities, city)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *addressRepositoryPostgres) FindOneCityByName(ctx context.Context, name string) (*int64, error) {
	var cityId int64

	if err := r.db.QueryRow(ctx, database.FindOneCityByName, name).Scan(&cityId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &cityId, nil
}

func (r *addressRepositoryPostgres) FindAllDistrictsByCityCode(ctx context.Context, cityCode string) ([]entity.District, error) {
	query := database.FindAllDistrictsByCityCode

	rows, err := r.db.Query(ctx, query, cityCode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	districts := []entity.District{}

	for rows.Next() {
		var district entity.District

		err := rows.Scan(
			&district.Id,
			&district.Code,
			&district.CityCode,
			&district.Name,
		)
		if err != nil {
			return nil, err
		}

		districts = append(districts, district)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return districts, nil
}

func (r *addressRepositoryPostgres) FindOneDistrictByName(ctx context.Context, name string) (*int64, error) {
	var districtId int64

	if err := r.db.QueryRow(ctx, database.FindOneDistrictByName, name).Scan(&districtId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &districtId, nil
}

func (r *addressRepositoryPostgres) FindAllSubdistrictsByDistrictCode(ctx context.Context, districtCode string) ([]entity.Subdistrict, error) {
	query := database.FindAllSubdistrictsByDistrictCode

	rows, err := r.db.Query(ctx, query, districtCode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subdistricts := []entity.Subdistrict{}

	for rows.Next() {
		var subdistrict entity.Subdistrict

		err := rows.Scan(
			&subdistrict.Id,
			&subdistrict.Code,
			&subdistrict.DistrictCode,
			&subdistrict.Name,
		)
		if err != nil {
			return nil, err
		}

		subdistricts = append(subdistricts, subdistrict)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subdistricts, nil
}

func (r *addressRepositoryPostgres) FindOneSubdistrictByName(ctx context.Context, name string) (*int64, error) {
	var subdistrictId int64

	if err := r.db.QueryRow(ctx, database.FindOneSubdistrictByName, name).Scan(&subdistrictId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &subdistrictId, nil
}
