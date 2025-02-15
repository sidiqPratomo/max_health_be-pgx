package repository

import (
	"context"
	// "database/sql"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type OrderItemRepository interface {
	FindAllByOrderPharmacyId(ctx context.Context, orderPharmacyId int64) ([]entity.OrderItem, error)
	FindPharmacyDrugCategorySalesVolumeRevenueByPharmacyId(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) ([]entity.DrugCategorySalesVolumeRevenue, error)
	FindPharmacyDrugSalesVolumeRevenueByPharmacyId(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) ([]entity.DrugSalesVolumeRevenue, error)
	PostOrderItems(ctx context.Context, orderPharmacies []entity.OrderPharmacyForCheckout) error
}

type orderItemRepositoryPostgres struct {
	db DBTX
}

func NewOrderItemRepositoryPostgres(db *pgxpool.Pool) orderItemRepositoryPostgres {
	return orderItemRepositoryPostgres{
		db: db,
	}
}

func (r *orderItemRepositoryPostgres) FindAllByOrderPharmacyId(ctx context.Context, orderPharmacyId int64) ([]entity.OrderItem, error) {
	rows, err := r.db.Query(ctx, database.FindAllOrderItemByOrderPharmacyId, orderPharmacyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []entity.OrderItem

	for rows.Next() {
		var orderItem entity.OrderItem

		err := rows.Scan(
			&orderItem.Id,
			&orderItem.DrugName,
			&orderItem.DrugPrice,
			&orderItem.DrugUnit,
			&orderItem.Quantity,
			&orderItem.DrugImage,
		)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (r *orderItemRepositoryPostgres) FindPharmacyDrugCategorySalesVolumeRevenueByPharmacyId(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) ([]entity.DrugCategorySalesVolumeRevenue, error) {
	sql := database.FindPharmacyDrugCategorySalesVolumeRevenue

	sql += ` ORDER BY sales_volume`
	if validatedGetReportQuery.Sort != nil {
		if *validatedGetReportQuery.Sort == "desc" {
			sql += ` DESC`
		} else {
			sql += ` ASC`
		}
	} else {
		sql += ` DESC`
	}

	rows, err := r.db.Query(ctx, sql, validatedGetReportQuery.PharmacyId, validatedGetReportQuery.MaxDate, validatedGetReportQuery.MinDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drugCategorySalesVolumeRevenues []entity.DrugCategorySalesVolumeRevenue

	for rows.Next() {
		var drugCategorySalesVolumeRevenue entity.DrugCategorySalesVolumeRevenue

		err := rows.Scan(
			&drugCategorySalesVolumeRevenue.DrugCategoryId,
			&drugCategorySalesVolumeRevenue.DrugCategoryName,
			&drugCategorySalesVolumeRevenue.SalesVolume,
			&drugCategorySalesVolumeRevenue.Revenue,
		)
		if err != nil {
			return nil, err
		}

		drugCategorySalesVolumeRevenues = append(drugCategorySalesVolumeRevenues, drugCategorySalesVolumeRevenue)
	}

	return drugCategorySalesVolumeRevenues, nil
}

func (r *orderItemRepositoryPostgres) FindPharmacyDrugSalesVolumeRevenueByPharmacyId(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) ([]entity.DrugSalesVolumeRevenue, error) {
	sql := database.FindPharmacyDrugSalesVolumeRevenue

	sql += ` ORDER BY sales_volume`
	if validatedGetReportQuery.Sort != nil {
		if *validatedGetReportQuery.Sort == "desc" {
			sql += ` DESC`
		} else {
			sql += ` ASC`
		}
	} else {
		sql += ` DESC`
	}

	rows, err := r.db.Query(ctx, sql, validatedGetReportQuery.PharmacyId, validatedGetReportQuery.MaxDate, validatedGetReportQuery.MinDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drugSalesVolumeRevenues []entity.DrugSalesVolumeRevenue

	for rows.Next() {
		var drugSalesVolumeRevenue entity.DrugSalesVolumeRevenue

		err := rows.Scan(
			&drugSalesVolumeRevenue.DrugId,
			&drugSalesVolumeRevenue.DrugName,
			&drugSalesVolumeRevenue.SalesVolume,
			&drugSalesVolumeRevenue.Revenue,
		)
		if err != nil {
			return nil, err
		}

		drugSalesVolumeRevenues = append(drugSalesVolumeRevenues, drugSalesVolumeRevenue)
	}

	return drugSalesVolumeRevenues, nil
}

func (r *orderItemRepositoryPostgres) PostOrderItems(ctx context.Context, orderPharmacies []entity.OrderPharmacyForCheckout) error {
	query := database.CreateOrderItems
	args := []interface{}{}
	for i, orderPharmacy := range orderPharmacies {
		args = append(args, orderPharmacy.Id)
		index := len(args)
		for j, cartItem := range orderPharmacy.CartItems {
			query += `($` + strconv.Itoa(index) + `, $` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) +
				`, $` + strconv.Itoa(len(args)+3) + `, $` + strconv.Itoa(len(args)+4) + `, $` + strconv.Itoa(len(args)+5) +
				`, $` + strconv.Itoa(len(args)+6) + `)`
			args = append(args, cartItem.DrugId)
			args = append(args, cartItem.DrugName)
			args = append(args, cartItem.PharmacyDrugId)
			args = append(args, cartItem.Price)
			args = append(args, cartItem.Unit)
			args = append(args, cartItem.Quantity)
			if !(i == len(orderPharmacies)-1 && j == len(orderPharmacy.CartItems)-1) {
				query += `,`
			}
		}
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
