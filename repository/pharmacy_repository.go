package repository

import (
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type PharmacyRepository interface {
	FindAllByManagerId(ctx context.Context, managerId int64, limit int, offset int, search string) ([]entity.Pharmacy, *entity.PageInfo, error)
	FindOneById(ctx context.Context, id int64) (*entity.Pharmacy, error)
	CreateOne(ctx context.Context, pharmacy *entity.Pharmacy) (*int64, error)
	UpdateOne(ctx context.Context, pharmacy entity.Pharmacy) error
	DeleteOneById(ctx context.Context, id int64) error
	GetAllCourierOptionsByPharmacyId(ctx context.Context, userAddressId, pharmacyId int64, weight float64) ([]entity.AvailableCourier, error)
	GetOnePharmacyByPharmacyId(ctx context.Context, pharmacyId int64) (*entity.Pharmacy, error)
}

type pharmacyRepositoryPostgres struct {
	db DBTX
}

func NewPharmacyRepositoryPostgres(db *sql.DB) pharmacyRepositoryPostgres {
	return pharmacyRepositoryPostgres{
		db: db,
	}
}

func (r *pharmacyRepositoryPostgres) CreateOne(ctx context.Context, pharmacy *entity.Pharmacy) (*int64, error) {
	var pharmacyId int64

	floatLatitude, _ := strconv.ParseFloat(pharmacy.Latitude, 32)
	floatLongitude, _ := strconv.ParseFloat(pharmacy.Longitude, 32)

	if err := r.db.QueryRowContext(ctx, database.CreateOnePharmacy, pharmacy.PharmacyManagerId, pharmacy.Name, pharmacy.PharmacistName, pharmacy.PharmacistLicenseNumber, pharmacy.PharmacistPhoneNumber, pharmacy.Address, pharmacy.City, pharmacy.Latitude, pharmacy.Longitude, floatLatitude, floatLongitude).Scan(&pharmacyId); err != nil {
		return nil, err
	}

	return &pharmacyId, nil
}

func (r *pharmacyRepositoryPostgres) FindOneById(ctx context.Context, id int64) (*entity.Pharmacy, error) {
	var pharmacy entity.Pharmacy

	if err := r.db.QueryRowContext(ctx, database.FindOnePharmacyById, id).Scan(&pharmacy.Id, &pharmacy.Name, &pharmacy.PharmacyManagerId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &pharmacy, nil
}

func (r *pharmacyRepositoryPostgres) UpdateOne(ctx context.Context, pharmacy entity.Pharmacy) error {
	floatLatitude, _ := strconv.ParseFloat(pharmacy.Latitude, 32)
	floatLongitude, _ := strconv.ParseFloat(pharmacy.Longitude, 32)

	_, err := r.db.ExecContext(ctx, database.UpdateOnePharmacy, pharmacy.Name, pharmacy.PharmacistName, pharmacy.PharmacistLicenseNumber, pharmacy.PharmacistPhoneNumber, pharmacy.Address, pharmacy.City, pharmacy.Latitude, pharmacy.Longitude, floatLatitude, floatLongitude, pharmacy.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *pharmacyRepositoryPostgres) DeleteOneById(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, database.DeleteOnePharmacyById, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *pharmacyRepositoryPostgres) GetOnePharmacyByPharmacyId(ctx context.Context, pharmacyId int64) (*entity.Pharmacy, error) {
	pharmacy := entity.Pharmacy{}

	err := r.db.QueryRowContext(ctx, database.GetOnePharmacyByPharmacyId, pharmacyId).
		Scan(
			&pharmacy.Id,
			&pharmacy.PharmacyManagerId,
			&pharmacy.Name,
			&pharmacy.PharmacistName,
			&pharmacy.PharmacistLicenseNumber,
			&pharmacy.PharmacistPhoneNumber,
			&pharmacy.City,
			&pharmacy.Address,
			&pharmacy.Latitude,
			&pharmacy.Longitude,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &pharmacy, nil
}

func (r *pharmacyRepositoryPostgres) FindAllByManagerId(ctx context.Context, managerId int64, limit int, offset int, search string) ([]entity.Pharmacy, *entity.PageInfo, error) {
	pharmacies := []entity.Pharmacy{}
	pageInfo := &entity.PageInfo{}

	rows, err := r.db.QueryContext(ctx, database.GetAllPharmacyByPharmacyManagerId, managerId, limit, offset, search)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pharmacy := entity.Pharmacy{}

		err := rows.Scan(
			&pharmacy.Id,
			&pharmacy.PharmacyManagerId,
			&pharmacy.Name,
			&pharmacy.PharmacistName,
			&pharmacy.PharmacistLicenseNumber,
			&pharmacy.PharmacistPhoneNumber,
			&pharmacy.City,
			&pharmacy.Address,
			&pharmacy.Longitude,
			&pharmacy.Latitude)
		if err != nil {
			return nil, nil, err
		}
		pharmacies = append(pharmacies, pharmacy)
	}
	defer rows.Close()

	countQuery := `
		SELECT COUNT(*) 
		FROM pharmacies
		where pharmacy_manager_id = $1 and pharmacy_name ILIKE '%' || $2 || '%' and deleted_at is null
	`
	countRow := r.db.QueryRowContext(ctx, countQuery, managerId, search)
	if err := countRow.Scan(&pageInfo.ItemCount); err != nil {
		return nil, nil, err
	}

	pageInfo.PageCount = int(math.Ceil(float64(pageInfo.ItemCount) / float64(limit)))
	pageInfo.Page = int(math.Ceil(float64(offset+1) / float64(limit)))

	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	return pharmacies, pageInfo, nil
}

func (r *pharmacyRepositoryPostgres) GetAllCourierOptionsByPharmacyId(ctx context.Context, userAddressId, pharmacyId int64, weight float64) ([]entity.AvailableCourier, error) {
	var availableCourierList []entity.AvailableCourier

	rows, err := r.db.QueryContext(ctx, database.GetAllCourierOptionsByPharmacyId, userAddressId, pharmacyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var availableCourier entity.AvailableCourier
		var courierOption entity.CourierOption
		var origin *int64
		var destination *int64

		err = rows.Scan(&availableCourier.PharmacyCourierId, &availableCourier.CourierName, &courierOption.Price, &origin, &destination)
		if err != nil {
			return nil, err
		}

		if availableCourier.CourierName == "Official Instant" {
			courierOption.Etd = "2-4 hours"
			availableCourier.CourierOptions = append(availableCourier.CourierOptions, courierOption)
		} else if availableCourier.CourierName == "Official Same Day" {
			courierOption.Etd = "1 day"
			availableCourier.CourierOptions = append(availableCourier.CourierOptions, courierOption)
		} else {
			if origin != nil && destination != nil {
				options, err := util.GetUnofficialDelivery(int64(*origin), int64(*destination), int64(weight), availableCourier.CourierName)
				if err != nil {
					return nil, err
				}
				availableCourier.CourierOptions = options
			}
		}

		availableCourierList = append(availableCourierList, availableCourier)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return availableCourierList, nil
}
