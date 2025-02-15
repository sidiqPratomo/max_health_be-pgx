package repository

import (
	"context"
	// "database/sql"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/util"
	"golang.org/x/exp/maps"
)

type OrderRepository interface {
	PostOneOrder(ctx context.Context, userId int64, address string, amount int) (int64, error)
	FindAllPendingByUserId(ctx context.Context, userId int64, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error)
	FindAllPendingWithDetailsByUserId(ctx context.Context, userId int64, orderIds []int64) ([]*entity.Order, error)
	FindAll(ctx context.Context, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error)
	FindAllWithDetails(ctx context.Context, orderIds []int64) ([]*entity.Order, error)
	UpdatePaymentProofOne(ctx context.Context, order *entity.Order) error
	FindOneOrderByOrderId(ctx context.Context, orderId int64) (*entity.Order, error)
}

type orderRepositoryPostgres struct {
	db DBTX
}

func NewOrderRepositoryPostgres(db *pgxpool.Pool) orderRepositoryPostgres {
	return orderRepositoryPostgres{
		db: db,
	}
}

func (r *orderRepositoryPostgres) PostOneOrder(ctx context.Context, userId int64, address string, amount int) (int64, error) {
	query := database.CreateOneOrder
	var orderId int64
	err := r.db.QueryRow(ctx, query, userId, address, amount).Scan(&orderId)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (r *orderRepositoryPostgres) FindAllPendingByUserId(ctx context.Context, userId int64, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error) {
	rows, err := r.db.Query(ctx, database.FindAllPendingOrdersByUserId, userId, validatedGetOrderQuery.Limit, validatedGetOrderQuery.Limit*(validatedGetOrderQuery.Page-1))
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

func (r *orderRepositoryPostgres) FindAllPendingWithDetailsByUserId(ctx context.Context, userId int64, orderIds []int64) ([]*entity.Order, error) {
	query := database.FindAllPendingOrdersWithDetailsByUserId
	args := []interface{}{}
	args = append(args, userId)

	if len(orderIds) > 0 {
		query += "AND o.order_id IN ("
		for i := 0; i < len(orderIds); i++ {
			args = append(args, orderIds[i])
			if i == len(orderIds)-1 {
				query += "$" + strconv.Itoa(len(args)) + ")"
			} else {
				query += "$" + strconv.Itoa(len(args)) + ","
			}
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []*entity.Order{}, err
	}
	defer rows.Close()

	ordersMap := map[int64]*entity.Order{}

	for rows.Next() {
		var order entity.Order
		var orderPharmacy entity.OrderPharmacy

		err := rows.Scan(
			&order.Id,
			&order.TotalAmount,
			&orderPharmacy.Id,
			&orderPharmacy.OrderStatusId,
			&orderPharmacy.SubtotalAmount,
			&orderPharmacy.DeliveryFee,
			&orderPharmacy.PharmacyName,
			&orderPharmacy.ProfilePicture,
			&orderPharmacy.OrderItemsCount,
		)
		if err != nil {
			return []*entity.Order{}, err
		}

		if existingOrder, exist := ordersMap[order.Id]; !exist {
			order.OrderPharmacies = append(order.OrderPharmacies, orderPharmacy)
			ordersMap[order.Id] = &order
		} else {
			existingOrder.OrderPharmacies = append(existingOrder.OrderPharmacies, orderPharmacy)
		}
	}

	if err = rows.Err(); err != nil {
		return []*entity.Order{}, err
	}

	return maps.Values(ordersMap), nil
}

func (r *orderRepositoryPostgres) FindAll(ctx context.Context, validatedGetOrderQuery util.ValidatedGetOrderQuery) ([]int64, *entity.PageInfo, error) {
	query := database.FindAllOrders
	args := []interface{}{}

	if validatedGetOrderQuery.StatusId != nil {
		query += `AND op.order_status_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, *validatedGetOrderQuery.StatusId)
	}

	query += " GROUP BY o.order_id"
	query += " ORDER BY o.created_at DESC"

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

func (r *orderRepositoryPostgres) FindAllWithDetails(ctx context.Context, orderIds []int64) ([]*entity.Order, error) {
	query := database.FindAllOrdersWithDetails
	args := []interface{}{}

	if len(orderIds) > 0 {
		query += "AND o.order_id IN ("
		for i := 0; i < len(orderIds); i++ {
			args = append(args, orderIds[i])
			if i == len(orderIds)-1 {
				query += "$" + strconv.Itoa(len(args)) + ")"
			} else {
				query += "$" + strconv.Itoa(len(args)) + ","
			}
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []*entity.Order{}, err
	}
	defer rows.Close()

	ordersMap := map[int64]*entity.Order{}

	for rows.Next() {
		var order entity.Order
		var orderPharmacy entity.OrderPharmacy

		err := rows.Scan(
			&order.Id,
			&order.TotalAmount,
			&order.PaymentProof,
			&order.Address,
			&order.UpdatedAt,
			&orderPharmacy.Id,
			&orderPharmacy.OrderStatusId,
			&orderPharmacy.SubtotalAmount,
			&orderPharmacy.DeliveryFee,
			&orderPharmacy.PharmacyName,
			&orderPharmacy.ProfilePicture,
			&orderPharmacy.CourierName,
			&orderPharmacy.OrderItemsCount,
		)
		if err != nil {
			return []*entity.Order{}, err
		}

		if existingOrder, exist := ordersMap[order.Id]; !exist {
			order.OrderPharmacies = append(order.OrderPharmacies, orderPharmacy)
			ordersMap[order.Id] = &order
		} else {
			existingOrder.OrderPharmacies = append(existingOrder.OrderPharmacies, orderPharmacy)
		}
	}

	if err = rows.Err(); err != nil {
		return []*entity.Order{}, err
	}

	return maps.Values(ordersMap), nil
}

func (r *orderRepositoryPostgres) UpdatePaymentProofOne(ctx context.Context, order *entity.Order) error {
	_, err := r.db.Exec(ctx, database.UpdatePaymentProofOneOrder, order.PaymentProof, order.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepositoryPostgres) FindOneOrderByOrderId(ctx context.Context, orderId int64) (*entity.Order, error) {
	var order entity.Order
	err := r.db.QueryRow(ctx, database.GetOneOrderByOrderId, orderId).Scan(&order.Id, &order.UserId, &order.Address,
		&order.PaymentProof, &order.TotalAmount, &order.ExpiredAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
