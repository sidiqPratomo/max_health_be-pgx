package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type UserAddressRepository interface {
	PostOneUserAddress(ctx context.Context, userAddress entity.UserAddress) error
	SetAllIsMainFalse(ctx context.Context, userId int64) error
	GetOneUserAddressByAddressId(ctx context.Context, addressId int64) (*entity.UserAddress, error)
	UpdateOneUserAddress(ctx context.Context, userAddress entity.UserAddress) error
	FindAllByUserId(ctx context.Context, userId int64) ([]entity.UserAddress, error)
	DeleteOneUserAddress(ctx context.Context, userAddressId int64) error
	FindOneUserAddressById(ctx context.Context, userAddressId int64) (*int64, error)
}

type userAddressRepositoryPostgres struct {
	db DBTX
}

func NewUserAddressRepositoryPostgres(db *pgxpool.Pool) userAddressRepositoryPostgres {
	return userAddressRepositoryPostgres{
		db: db,
	}
}

func (r *userAddressRepositoryPostgres) GetOneUserAddressByAddressId(ctx context.Context, addressId int64) (*entity.UserAddress, error) {
	var userAddress entity.UserAddress

	err := r.db.QueryRow(ctx, database.GetOneUserAddressByAddressIdQuery, addressId).Scan(&userAddress.UserId, &userAddress.Province.Id, &userAddress.City.Id, &userAddress.District.Id, &userAddress.Subdistrict.Id, &userAddress.Latitude, &userAddress.Longitude, &userAddress.Label, &userAddress.Address, &userAddress.IsActive, &userAddress.IsMain)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &userAddress, nil
}

func (r *userAddressRepositoryPostgres) SetAllIsMainFalse(ctx context.Context, userId int64) error {
	_, err := r.db.Exec(ctx, database.SetAllIsMainFalseQuery, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userAddressRepositoryPostgres) PostOneUserAddress(ctx context.Context, userAddress entity.UserAddress) error {
	floatLatitude, _ := strconv.ParseFloat(userAddress.Latitude, 32)
	floatLongitude, _ := strconv.ParseFloat(userAddress.Longitude, 32)
	_, err := r.db.Exec(ctx, database.PostOneUserAddressQuery, userAddress.UserId, userAddress.Province.Id, userAddress.City.Id, userAddress.District.Id, userAddress.Subdistrict.Id, userAddress.Latitude, userAddress.Longitude, userAddress.Label, userAddress.Address, userAddress.IsMain, floatLongitude, floatLatitude)
	if err != nil {
		return err
	}

	return nil
}

func (r *userAddressRepositoryPostgres) UpdateOneUserAddress(ctx context.Context, userAddress entity.UserAddress) error {
	floatLatitude, _ := strconv.ParseFloat(userAddress.Latitude, 32)
	floatLongitude, _ := strconv.ParseFloat(userAddress.Longitude, 32)

	_, err := r.db.Exec(ctx, database.UpdateUserAddressQuery, userAddress.Province.Id, userAddress.City.Id, userAddress.District.Id, userAddress.Subdistrict.Id, userAddress.Latitude, userAddress.Longitude, userAddress.Label, userAddress.Address, userAddress.IsActive, userAddress.IsMain, floatLongitude, floatLatitude, userAddress.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userAddressRepositoryPostgres) FindAllByUserId(ctx context.Context, userId int64) ([]entity.UserAddress, error) {
	query := database.FindAllUserAddressByUserIdQuery

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userAddress := []entity.UserAddress{}

	for rows.Next() {
		var address entity.UserAddress

		err := rows.Scan(
			&address.Id,
			&address.Province.Id,
			&address.Province.Code,
			&address.Province.Name,
			&address.City.Id,
			&address.City.Code,
			&address.City.Name,
			&address.District.Id,
			&address.District.Code,
			&address.District.Name,
			&address.Subdistrict.Id,
			&address.Subdistrict.Code,
			&address.Subdistrict.Name,
			&address.Latitude,
			&address.Longitude,
			&address.Label,
			&address.Address,
			&address.IsMain,
		)
		if err != nil {
			return nil, err
		}

		userAddress = append(userAddress, address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userAddress, nil
}

func (r *userAddressRepositoryPostgres) DeleteOneUserAddress(ctx context.Context, userAddressId int64) error {
	_, err := r.db.Exec(ctx, database.DeleteOneUserAddressQuery, userAddressId)
	if err != nil {
		return err
	}
	return nil
}

func (r *userAddressRepositoryPostgres) FindOneUserAddressById(ctx context.Context, userAddressId int64) (*int64, error) {
	var userId int64

	err := r.db.QueryRow(ctx, database.GetOneUserAddressById, userAddressId).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &userId, nil
}
