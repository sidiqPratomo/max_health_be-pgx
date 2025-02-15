package repository

import (
	"context"
	// "database/sql"
	"time"

	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/util"
	"golang.org/x/exp/maps"
)

type OrderPharmacyRepository interface {
	PostOrderPharmacies(ctx context.Context, orderId int64, orderCheckoutRequest dto.OrderCheckoutRequest) ([]entity.OrderPharmacyForCheckout, error)
	FindAllByOrderId(ctx context.Context, orderId int64) ([]entity.OrderPharmacy, error)
	UpdateStatusBulkByOrderId(ctx context.Context, orderId int64, newOrderStatusId int64) error
	FindAllOngoingIdsByPharmacyId(ctx context.Context, pharmacyId int64) ([]int64, error)
	FindOneById(ctx context.Context, id int64) (*entity.OrderPharmacy, error)
	FindAllByOrderUserId(ctx context.Context, userId int64, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]entity.OrderPharmacy, *entity.PageInfo, error)
	FindAllIdsByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error)
	FindAllWithDetailsByIds(ctx context.Context, orderPharmacyIds []int64) ([]*entity.OrderPharmacy, error)
	FindAllIds(ctx context.Context, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error)
	FindOneByOrderPharmacyId(ctx context.Context, orderPharmacyId int64) (*entity.OrderPharmacy, error)
	FindCountGroupedByOrderStatusIdByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64) (*entity.OrderPharmacySummary, error)
	UpdateOneStatusById(ctx context.Context, orderPharmacyId int64, newOrderStatusId int64) error
}

type orderPharmacyRepositoryPostgres struct {
	db DBTX
}

func NewOrderPharmacyRepositoryPostgres(db *pgxpool.Pool) orderPharmacyRepositoryPostgres {
	return orderPharmacyRepositoryPostgres{
		db: db,
	}
}

func (r *orderPharmacyRepositoryPostgres) PostOrderPharmacies(ctx context.Context, orderId int64, orderCheckoutRequest dto.OrderCheckoutRequest) ([]entity.OrderPharmacyForCheckout, error) {
	var orderPharmacyIds []entity.OrderPharmacyForCheckout
	query := database.CreateOrderPharmacies
	args := []interface{}{}
	args = append(args, orderId)
	for i, pharmacy := range orderCheckoutRequest.Pharmacies {
		query += `($1, 1, $` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) + `)`
		args = append(args, pharmacy.PharmacyCourierId, pharmacy.Subtotal, pharmacy.DeliveryFee)
		if i != len(orderCheckoutRequest.Pharmacies)-1 {
			query += `,`
		}
	}
	query += `RETURNING order_pharmacy_id`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return orderPharmacyIds, err
	}
	defer rows.Close()

	for rows.Next() {
		orderPharmacyId := entity.OrderPharmacyForCheckout{}

		err := rows.Scan(&orderPharmacyId.Id)
		if err != nil {
			return []entity.OrderPharmacyForCheckout{}, err
		}
		orderPharmacyIds = append(orderPharmacyIds, orderPharmacyId)
	}

	return orderPharmacyIds, nil
}

func (r *orderPharmacyRepositoryPostgres) FindAllByOrderId(ctx context.Context, orderId int64) ([]entity.OrderPharmacy, error) {
	query := database.FindAllOrderPharmaciesByOrderId

	rows, err := r.db.Query(ctx, query, orderId)
	if err != nil {
		return []entity.OrderPharmacy{}, err
	}

	defer rows.Close()

	orderPharmacies := []entity.OrderPharmacy{}

	for rows.Next() {
		var orderPharmacy entity.OrderPharmacy

		err := rows.Scan(
			&orderPharmacy.Id,
			&orderPharmacy.UserId,
			&orderPharmacy.OrderStatusId,
		)
		if err != nil {
			return []entity.OrderPharmacy{}, err
		}

		orderPharmacies = append(orderPharmacies, orderPharmacy)
	}

	if err = rows.Err(); err != nil {
		return []entity.OrderPharmacy{}, err
	}

	return orderPharmacies, nil
}

func (r *orderPharmacyRepositoryPostgres) FindAllOngoingIdsByPharmacyId(ctx context.Context, pharmacyId int64) ([]int64, error) {
	rows, err := r.db.Query(ctx, database.FindAllOngoingOrderPharmacyIdsByPharmacyId, pharmacyId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orderPharmacyIds := []int64{}

	for rows.Next() {
		var orderPharmacyId int64

		err := rows.Scan(
			&orderPharmacyId,
		)
		if err != nil {
			return nil, err
		}

		orderPharmacyIds = append(orderPharmacyIds, orderPharmacyId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orderPharmacyIds, nil
}

func (r *orderPharmacyRepositoryPostgres) UpdateStatusBulkByOrderId(ctx context.Context, orderId int64, newOrderStatusId int64) error {
	_, err := r.db.Exec(ctx, database.UpdateStatusBulkOrderPharmaciesByOrderId, newOrderStatusId, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderPharmacyRepositoryPostgres) FindOneById(ctx context.Context, id int64) (*entity.OrderPharmacy, error) {
	var orderPharmacy entity.OrderPharmacy

	if err := r.db.QueryRow(ctx, database.FindOneOrderPharmacyById, id).Scan(&orderPharmacy.Id, &orderPharmacy.OrderStatusId, &orderPharmacy.SubtotalAmount, &orderPharmacy.DeliveryFee, &orderPharmacy.CreatedAt, &orderPharmacy.UpdatedAt, &orderPharmacy.PharmacyName, &orderPharmacy.CourierName, &orderPharmacy.ProfilePicture, &orderPharmacy.Address); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &orderPharmacy, nil
}

func (r *orderPharmacyRepositoryPostgres) FindAllByOrderUserId(ctx context.Context, userId int64, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]entity.OrderPharmacy, *entity.PageInfo, error) {
	query := database.FindAllOrderPharmaciesByOrderUserId
	args := []interface{}{}
	args = append(args, userId)

	if validatedGetOrderQuery.StatusId != nil {
		query += `AND op.order_status_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, *validatedGetOrderQuery.StatusId)
	}

	query += ` GROUP BY op.order_pharmacy_id, p.pharmacy_name, a.profile_picture, first_order_item.drug_name, first_order_item.drug_price, first_order_item.drug_unit, first_order_item.quantity, first_order_item.image, first_order_item.total_order_items`
	query += ` ORDER BY op.updated_at DESC`

	query += ` LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, validatedGetOrderQuery.Limit)
	query += ` OFFSET $` + strconv.Itoa(len(args)+1)
	args = append(args, (validatedGetOrderQuery.Limit * (validatedGetOrderQuery.Page - 1)))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []entity.OrderPharmacy{}, nil, err
	}
	defer rows.Close()

	orderPharmacies := []entity.OrderPharmacy{}
	var pageInfo entity.PageInfo

	for rows.Next() {
		var orderPharmacy entity.OrderPharmacy
		err := rows.Scan(
			&orderPharmacy.Id,
			&orderPharmacy.OrderStatusId,
			&orderPharmacy.SubtotalAmount,
			&orderPharmacy.DeliveryFee,
			&orderPharmacy.PharmacyName,
			&orderPharmacy.ProfilePicture,
			&orderPharmacy.FirstOrderItem.DrugName,
			&orderPharmacy.FirstOrderItem.DrugPrice,
			&orderPharmacy.FirstOrderItem.DrugUnit,
			&orderPharmacy.FirstOrderItem.Quantity,
			&orderPharmacy.FirstOrderItem.DrugImage,
			&orderPharmacy.OrderItemsCount,
			&pageInfo.ItemCount,
		)
		if err != nil {
			return []entity.OrderPharmacy{}, nil, err
		}
		orderPharmacies = append(orderPharmacies, orderPharmacy)
	}

	if err = rows.Err(); err != nil {
		return []entity.OrderPharmacy{}, nil, err
	}

	pageInfo.Page = validatedGetOrderQuery.Page
	pageInfo.PageCount = pageInfo.ItemCount / validatedGetOrderQuery.Limit
	if pageInfo.ItemCount%validatedGetOrderQuery.Limit != 0 {
		pageInfo.PageCount += 1
	}

	return orderPharmacies, &pageInfo, nil
}

func (r *orderPharmacyRepositoryPostgres) FindAllIdsByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error) {
	query := database.FindAllOrderPharmaciesByPharmacyManagerId
	args := []interface{}{}
	args = append(args, pharmacyManagerId)

	if validatedGetOrderQuery.StatusId != nil {
		query += `AND op.order_status_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, *validatedGetOrderQuery.StatusId)
	}

	if validatedGetOrderQuery.PharmacyName != nil {
		query += `AND p.pharmacy_name ILIKE `
		query += `$` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+*validatedGetOrderQuery.PharmacyName+"%")
	}

	query += ` ORDER BY op.updated_at DESC`

	query += ` LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, validatedGetOrderQuery.Limit)
	query += ` OFFSET $` + strconv.Itoa(len(args)+1)
	args = append(args, (validatedGetOrderQuery.Limit * (validatedGetOrderQuery.Page - 1)))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []int64{}, nil, err
	}
	defer rows.Close()

	orderIds := []int64{}
	var pageInfo entity.PageInfo

	for rows.Next() {
		var orderId int64

		err := rows.Scan(
			&orderId,
			&time.Time{},
			&pageInfo.ItemCount,
		)
		if err != nil {
			return []int64{}, nil, err
		}

		orderIds = append(orderIds, orderId)
	}

	if err = rows.Err(); err != nil {
		return []int64{}, nil, err
	}

	pageInfo.Page = validatedGetOrderQuery.Page
	pageInfo.PageCount = pageInfo.ItemCount / validatedGetOrderQuery.Limit
	if pageInfo.ItemCount%validatedGetOrderQuery.Limit != 0 {
		pageInfo.PageCount += 1
	}

	return orderIds, &pageInfo, nil
}

func (r *orderPharmacyRepositoryPostgres) FindAllWithDetailsByIds(ctx context.Context, orderPharmacyIds []int64) ([]*entity.OrderPharmacy, error) {
	query := database.FindAllOrderPharmaciesWithDetails
	args := []interface{}{}

	if len(orderPharmacyIds) > 0 {
		query += "AND op.order_pharmacy_id IN ("
		for i := 0; i < len(orderPharmacyIds); i++ {
			args = append(args, orderPharmacyIds[i])
			if i == len(orderPharmacyIds)-1 {
				query += "$" + strconv.Itoa(len(args)) + ")"
			} else {
				query += "$" + strconv.Itoa(len(args)) + ","
			}
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []*entity.OrderPharmacy{}, err
	}
	defer rows.Close()

	orderPharmaciesMap := map[int64]*entity.OrderPharmacy{}

	for rows.Next() {
		var orderPharmacy entity.OrderPharmacy
		var orderItem entity.OrderItem

		err := rows.Scan(
			&orderPharmacy.Id,
			&orderPharmacy.OrderStatusId,
			&orderPharmacy.SubtotalAmount,
			&orderPharmacy.DeliveryFee,
			&orderPharmacy.UpdatedAt,
			&orderPharmacy.PharmacyName,
			&orderPharmacy.PharmacistPhoneNumber,
			&orderPharmacy.CourierName,
			&orderPharmacy.PharmacyManagerEmail,
			&orderItem.DrugName,
			&orderItem.DrugPrice,
			&orderItem.DrugUnit,
			&orderItem.DrugImage,
			&orderItem.Quantity,
		)
		if err != nil {
			return []*entity.OrderPharmacy{}, err
		}

		if existingOrderPharmacy, exist := orderPharmaciesMap[orderPharmacy.Id]; !exist {
			orderPharmacy.OrderItems = append(orderPharmacy.OrderItems, orderItem)
			orderPharmaciesMap[orderPharmacy.Id] = &orderPharmacy
		} else {
			existingOrderPharmacy.OrderItems = append(existingOrderPharmacy.OrderItems, orderItem)
		}
	}

	if err = rows.Err(); err != nil {
		return []*entity.OrderPharmacy{}, err
	}

	return maps.Values(orderPharmaciesMap), nil
}

func (r *orderPharmacyRepositoryPostgres) FindAllIds(ctx context.Context, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error) {
	query := database.FindAllOrderPharmacies
	args := []interface{}{}

	if validatedGetOrderQuery.StatusId != nil {
		query += `AND op.order_status_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, *validatedGetOrderQuery.StatusId)
	}

	if validatedGetOrderQuery.PharmacyName != nil {
		query += `AND p.pharmacy_name ILIKE `
		query += `$` + strconv.Itoa(len(args)+1)
		args = append(args, "%"+*validatedGetOrderQuery.PharmacyName+"%")
	}

	query += ` ORDER BY op.updated_at DESC`

	query += ` LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, validatedGetOrderQuery.Limit)
	query += ` OFFSET $` + strconv.Itoa(len(args)+1)
	args = append(args, (validatedGetOrderQuery.Limit * (validatedGetOrderQuery.Page - 1)))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []int64{}, nil, err
	}
	defer rows.Close()

	orderIds := []int64{}
	var pageInfo entity.PageInfo

	for rows.Next() {
		var orderId int64

		err := rows.Scan(
			&orderId,
			&time.Time{},
			&pageInfo.ItemCount,
		)
		if err != nil {
			return []int64{}, nil, err
		}

		orderIds = append(orderIds, orderId)
	}

	if err = rows.Err(); err != nil {
		return []int64{}, nil, err
	}

	pageInfo.Page = validatedGetOrderQuery.Page
	pageInfo.PageCount = pageInfo.ItemCount / validatedGetOrderQuery.Limit
	if pageInfo.ItemCount%validatedGetOrderQuery.Limit != 0 {
		pageInfo.PageCount += 1
	}

	return orderIds, &pageInfo, nil
}

func (r *orderPharmacyRepositoryPostgres) FindOneByOrderPharmacyId(ctx context.Context, orderPharmacyId int64) (*entity.OrderPharmacy, error) {
	var orderPharmacy entity.OrderPharmacy
	err := r.db.QueryRow(ctx, database.FindOrderPharmacyByOrderPharmacyId, orderPharmacyId).Scan(&orderPharmacy.Id, &orderPharmacy.UserId, &orderPharmacy.OrderId,
		&orderPharmacy.OrderStatusId, &orderPharmacy.PharmacyCourierId, &orderPharmacy.SubtotalAmount, &orderPharmacy.DeliveryFee)
	if err != nil {
		return nil, err
	}

	return &orderPharmacy, err
}

func (r *orderPharmacyRepositoryPostgres) FindCountGroupedByOrderStatusIdByPharmacyManagerId(ctx context.Context, pharmacyManagerId int64) (*entity.OrderPharmacySummary, error) {
	rows, err := r.db.Query(ctx, database.FindOrderPharmacyCountGroupedByOrderStatusIdByPharmacyManagerId, pharmacyManagerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderPharmacySummaryMap := map[int64]int64{}

	for rows.Next() {
		var orderStatusId int64
		var count int64

		if err := rows.Scan(
			&orderStatusId,
			&count,
		); err != nil {
			return nil, err
		}

		orderPharmacySummaryMap[orderStatusId] = count
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var allCount int64
	for _, count := range maps.Values(orderPharmacySummaryMap) {
		allCount += count
	}

	return &entity.OrderPharmacySummary{
		AllCount:       int64(allCount),
		UnpaidCount:    orderPharmacySummaryMap[1],
		ApprovalCount:  orderPharmacySummaryMap[2],
		PendingCount:   orderPharmacySummaryMap[3],
		SentCount:      orderPharmacySummaryMap[4],
		ConfirmedCount: orderPharmacySummaryMap[5],
		CanceledCount:  orderPharmacySummaryMap[6],
	}, nil
}

func (r *orderPharmacyRepositoryPostgres) UpdateOneStatusById(ctx context.Context, orderPharmacyId int64, newOrderStatusId int64) error {
	_, err := r.db.Exec(ctx, database.UpdateOneStatusById, newOrderStatusId, orderPharmacyId)
	if err != nil {
		return err
	}

	return nil
}
