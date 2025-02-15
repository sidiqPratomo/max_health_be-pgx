package repository

import (
	"context"
	// "database/sql"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type StockChangeRepository interface {
	PostStockChangesByCartIds(ctx context.Context, cartItems []entity.CartItemChanges) error
	PostStockChangesFromMutation(ctx context.Context, stockChangesList []entity.StockChange) error
	PostStockChanges(ctx context.Context, stockChanges []entity.StockChange) error
	PostStockChangesFromUpdate(ctx context.Context, stockChanges []entity.StockChange) error
	GetStockChanges(ctx context.Context, managerId int64, pharmacyId *int64) ([]dto.StockChangeResponse, error)
}

type stockChangeRepositoryPostgres struct {
	db DBTX
}

func NewStockChangeRepositoryPostgres(db *pgxpool.Pool) stockChangeRepositoryPostgres {
	return stockChangeRepositoryPostgres{
		db: db,
	}
}

func (r *stockChangeRepositoryPostgres) PostStockChangesByCartIds(ctx context.Context, cartItems []entity.CartItemChanges) error {
	query := database.CreateStockChanges

	args := []interface{}{}
	for i, cartItem := range cartItems {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) +
			`, $` + strconv.Itoa(len(args)+4) + `)`
		args = append(args, cartItem.PharmacyDrugId)
		args = append(args, cartItem.Stock-cartItem.Quantity)
		args = append(args, -1*cartItem.Quantity)
		args = append(args, "bought by customer")
		if i != len(cartItems)-1 {
			query += `,`
		}
	}
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *stockChangeRepositoryPostgres) PostStockChangesFromMutation(ctx context.Context, stockChangesList []entity.StockChange) error {
	query := database.CreateStockChanges
	args := []interface{}{}
	for i, stockChange := range stockChangesList {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) +
			`, $` + strconv.Itoa(len(args)+4) + `)`
		args = append(args, stockChange.PharmacyDrugId)
		args = append(args, stockChange.FinalStock)
		args = append(args, stockChange.Amount)
		args = append(args, "transfer from stock mutation")
		if i != len(stockChangesList)-1 {
			query += `,`
		}
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *stockChangeRepositoryPostgres) PostStockChanges(ctx context.Context, stockChanges []entity.StockChange) error {
	query := database.CreateStockChanges
	args := []interface{}{}
	for i, stockChange := range stockChanges {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) + `, $` + strconv.Itoa(len(args)+4) + `)`
		if i != len(stockChanges)-1 {
			query += `,`
		}
		args = append(args, stockChange.PharmacyDrugId)
		args = append(args, stockChange.FinalStock)
		args = append(args, stockChange.Amount)
		args = append(args, "transfer from cancelled order")
	}
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *stockChangeRepositoryPostgres) PostStockChangesFromUpdate(ctx context.Context, stockChanges []entity.StockChange) error {
	query := database.CreateStockChanges
	args := []interface{}{}
	for i, stockChange := range stockChanges {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) + `, $` + strconv.Itoa(len(args)+4) + `)`
		if i != len(stockChanges)-1 {
			query += `,`
		}
		args = append(args, stockChange.PharmacyDrugId)
		args = append(args, stockChange.FinalStock)
		args = append(args, stockChange.Amount)
		args = append(args, stockChange.Description)
	}
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *stockChangeRepositoryPostgres) GetStockChanges(ctx context.Context, managerId int64, pharmacyId *int64) ([]dto.StockChangeResponse, error) {
	stockChanges := []dto.StockChangeResponse{}
	query := database.GetStockChanges
	args := []interface{}{}
	args = append(args, managerId)
	if pharmacyId != nil {
		query += `AND p.pharmacy_id = $2 `
		args = append(args, *pharmacyId)
	}
	query += `ORDER BY sc.created_at`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return stockChanges, err
	}
	defer rows.Close()
	for rows.Next() {
		var stockChange dto.StockChangeResponse
		err = rows.Scan(&stockChange.PharmacyName, &stockChange.PharmacyAddress, &stockChange.DrugName, &stockChange.DrugImage, &stockChange.FinalStock, &stockChange.Change, &stockChange.Description)
		if err != nil {
			return []dto.StockChangeResponse{}, err
		}
		stockChanges = append(stockChanges, stockChange)
	}
	return stockChanges, nil
}
