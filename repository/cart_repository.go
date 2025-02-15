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

type CartRepository interface {
	PostOneCart(ctx context.Context, accountID int64, pharmacyDrugId int64, quantity int) (*int64, error)
	UpdateOneCart(ctx context.Context, accountID int64, cartItemID int64, quantity int) error
	DeleteOneCart(ctx context.Context, accountID int64, cartItemID int64) error
	GetAllCart(ctx context.Context, accountID int64, Limit string, offset int) ([]entity.CartItemData, *entity.PageInfo, error)
	GetPharmacyDeliveryFeeForCart(ctx context.Context, cartItemsId []int64, userAddressId int64) ([]entity.PharmacyDeliveryFee, error)
	GetCartsByIds(ctx context.Context, cartItemsIds []int64) ([]entity.CartItem, error)
	GetStockByCartId(ctx context.Context, cartItemId int64) (*int, error)
	GetAllCartDetailByIds(ctx context.Context, cartItemIds []int64) ([]entity.CartItemForCheckout, error)
	GetAllCartsForChangesByCartIds(ctx context.Context, cartItems []entity.CartItemForCheckout) ([]entity.CartItemChanges, error)
	DeleteCarts(ctx context.Context, cartItems []entity.CartItemForCheckout) error
}

type cartRepositoryPostgres struct {
	db DBTX
}

func NewCartRepositoryPostgres(db *sql.DB) cartRepositoryPostgres {
	return cartRepositoryPostgres{
		db: db,
	}
}

func (r *cartRepositoryPostgres) PostOneCart(ctx context.Context, accountID int64, pharmacyDrugId int64, quantity int) (*int64, error) {
	var userID string
	err := r.db.QueryRowContext(ctx, database.CheckUserQuery, accountID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var cartItemId *int64
	err = r.db.QueryRowContext(ctx, database.PostOneCartQuery, userID, pharmacyDrugId, quantity).Scan(&cartItemId)
	if err != nil {
		return nil, err
	}

	return cartItemId, nil
}

func (r *cartRepositoryPostgres) UpdateOneCart(ctx context.Context, accountID int64, cartItemID int64, quantity int) error {
	var userID string
	err := r.db.QueryRowContext(ctx, database.CheckUserQuery, accountID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	_, err = r.db.ExecContext(ctx, database.UpdateOneCartQuery, cartItemID, userID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepositoryPostgres) DeleteOneCart(ctx context.Context, accountID int64, cartItemID int64) error {
	var userID string
	err := r.db.QueryRowContext(ctx, database.CheckUserQuery, accountID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	_, err = r.db.ExecContext(ctx, database.DeleteOneCartQuery, cartItemID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepositoryPostgres) GetStockByCartId(ctx context.Context, cartItemId int64) (*int, error) {
	var stock int
	err := r.db.QueryRowContext(ctx, database.GetStockPharmacyDrug, cartItemId).Scan(&stock)
	if err != nil {
		return nil, err
	}

	return &stock, err
}

func (r *cartRepositoryPostgres) GetAllCart(ctx context.Context, accountID int64, Limit string, offset int) ([]entity.CartItemData, *entity.PageInfo, error) {
	var userID string
	carts := []entity.CartItemData{}
	pageInfo := &entity.PageInfo{}

	intLimit, err := strconv.Atoi(Limit)
	if err != nil {
		return nil, nil, err
	}

	err = r.db.QueryRowContext(ctx, database.CheckUserQuery, accountID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	rows, err := r.db.QueryContext(ctx, database.GetAllCartQuery, userID, intLimit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}

		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cart := entity.CartItemData{}

		err := rows.Scan(&cart.CartItemId, &cart.UserId, &cart.PharmacyDrugId, &cart.Quantity, &cart.PharmacyId, &cart.DrugId, &cart.Price, &cart.Stock, &cart.DrugName, &cart.Image, &cart.PharmacyName)
		if err != nil {
			return nil, nil, err
		}
		carts = append(carts, cart)
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM cart_items WHERE deleted_at ISNULL
	`
	countRow := r.db.QueryRowContext(ctx, countQuery)
	if err := countRow.Scan(&pageInfo.ItemCount); err != nil {
		return nil, nil, err
	}

	pageInfo.PageCount = int(math.Ceil(float64(pageInfo.ItemCount) / float64(intLimit)))
	pageInfo.Page = int(math.Ceil(float64(offset+1) / float64(intLimit)))

	return carts, pageInfo, nil
}

func (r *cartRepositoryPostgres) GetPharmacyDeliveryFeeForCart(ctx context.Context, cartItemsId []int64, userAddressId int64) ([]entity.PharmacyDeliveryFee, error) {
	query := database.GetAllDeliveryFee1
	args := []interface{}{}
	for i, cart := range cartItemsId {
		if i != len(cartItemsId)-1 {
			query += `ci.cart_item_id = $` + strconv.Itoa(len(args)+1) + ` OR `
		} else {
			query += `ci.cart_item_id = $` + strconv.Itoa(len(args)+1)
		}
		args = append(args, cart)
	}
	query += `)),` + database.GetAllDeliveryFee2 + strconv.Itoa(len(args)+1) + `),` + database.GetAllDeliveryFee3
	args = append(args, userAddressId)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryFees := []entity.PharmacyDeliveryFee{}
	deliveryFee := entity.PharmacyDeliveryFee{}
	couriers := []entity.AvailableCourier{}

	for rows.Next() {
		pharmacy := ""
		pharmacyId := 0
		var origin *int
		var destination *int
		weight := 0
		distance := 0
		isActive := true
		isOfficial := true
		courier := entity.AvailableCourier{}
		courierOption := entity.CourierOption{}

		err := rows.Scan(
			&pharmacyId,
			&pharmacy,
			&distance,
			&courier.PharmacyCourierId,
			&courier.CourierName,
			&origin,
			&destination,
			&weight,
			&courierOption.Price,
			&isActive,
			&isOfficial,
		)
		if err != nil {
			return nil, err
		}

		if deliveryFee.PharmacyName == "" {
			deliveryFee.Id = int64(pharmacyId)
			deliveryFee.PharmacyName = pharmacy
			deliveryFee.Distance = distance
		}

		if courier.CourierName == "Official Instant" {
			courierOption.Etd = "2-4 hours"
			courier.CourierOptions = append(courier.CourierOptions, courierOption)
		} else if courier.CourierName == "Official Same Day" {
			courierOption.Etd = "1 day"
			courier.CourierOptions = append(courier.CourierOptions, courierOption)
		} else {
			if origin != nil && destination != nil {
				options, err := util.GetUnofficialDelivery(int64(*origin), int64(*destination), int64(weight), courier.CourierName)
				if err != nil {
					return nil, err
				}
				courier.CourierOptions = options
			}
		}

		if pharmacyId == int(deliveryFee.Id) {
			couriers = append(couriers, courier)
		}

		if pharmacyId != int(deliveryFee.Id) {
			deliveryFee.Couriers = couriers
			deliveryFees = append(deliveryFees, deliveryFee)
			deliveryFee = entity.PharmacyDeliveryFee{Id: int64(pharmacyId), PharmacyName: pharmacy, Distance: distance}
			couriers = []entity.AvailableCourier{courier}
		}
	}

	deliveryFee.Couriers = couriers
	deliveryFees = append(deliveryFees, deliveryFee)

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deliveryFees, nil
}

func (r *cartRepositoryPostgres) GetCartsByIds(ctx context.Context, cartItemsIds []int64) ([]entity.CartItem, error) {
	query := database.GetCartsByIds
	carts := []entity.CartItem{}
	args := []interface{}{}
	if len(cartItemsIds) > 0 {
		query += ` AND (`
		for i, cartId := range cartItemsIds {
			if i == 0 {
				query += `cart_item_id = $` + strconv.Itoa(len(args)+1)
			} else {
				query += ` OR cart_item_id = $` + strconv.Itoa(len(args)+1)
			}
			args = append(args, cartId)
		}
		query += `)`
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cart := entity.CartItem{}

		err := rows.Scan(&cart.Id, &cart.UserId)
		if err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func (r *cartRepositoryPostgres) GetAllCartDetailByIds(ctx context.Context, cartItemIds []int64) ([]entity.CartItemForCheckout, error) {
	query := database.GetAllDetailedCartItems
	cartItems := []entity.CartItemForCheckout{}
	args := []interface{}{}
	if len(cartItemIds) > 0 {
		query += `WHERE `
		for i := range cartItemIds {
			if i != len(cartItemIds)-1 {
				query += `cart_item_id = $` + strconv.Itoa(i+1) + ` OR `
			} else {
				query += `cart_item_id = $` + strconv.Itoa(i+1)
			}
			args = append(args, cartItemIds[i])
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return cartItems, err
	}
	defer rows.Close()

	for rows.Next() {
		cartItem := entity.CartItemForCheckout{}
		err := rows.Scan(&cartItem.Id, &cartItem.DrugId, &cartItem.DrugName, &cartItem.PharmacyDrugId, &cartItem.Price, &cartItem.Unit, &cartItem.Quantity)
		if err != nil {
			return []entity.CartItemForCheckout{}, err
		}
		cartItems = append(cartItems, cartItem)
	}
	return cartItems, nil
}

func (r *cartRepositoryPostgres) GetAllCartsForChangesByCartIds(ctx context.Context, cartItems []entity.CartItemForCheckout) ([]entity.CartItemChanges, error) {
	cartItemChanges := []entity.CartItemChanges{}
	query := database.GetAllCartsForChangesByCartIds
	args := []interface{}{}
	query += `WHERE `
	for i, cartItem := range cartItems {
		query += `ci.cart_item_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, cartItem.Id)
		if i != len(cartItems)-1 {
			query += ` OR `
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return cartItemChanges, err
	}
	defer rows.Close()

	for rows.Next() {
		var cartItemChange entity.CartItemChanges
		err := rows.Scan(&cartItemChange.PharmacyDrugId, &cartItemChange.Stock, &cartItemChange.Quantity)
		if err != nil {
			return []entity.CartItemChanges{}, err
		}
		cartItemChanges = append(cartItemChanges, cartItemChange)
	}
	return cartItemChanges, nil
}

func (r *cartRepositoryPostgres) DeleteCarts(ctx context.Context, cartItems []entity.CartItemForCheckout) error {
	query := database.DeleteAllCarts
	args := []interface{}{}
	query += ` WHERE `
	for i, cartItem := range cartItems {
		query += `cart_item_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, cartItem.Id)
		if i != len(cartItems)-1 {
			query += ` OR `
		}
	}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
