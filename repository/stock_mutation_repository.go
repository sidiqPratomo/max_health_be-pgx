package repository

import (
	"context"
	// "database/sql"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type StockMutationRepository interface {
	GetPossibleStockMutation(ctx context.Context, cartItems []entity.CartItemForCheckout) ([]entity.PossibleStockMutation, error)
	PostStockMutations(ctx context.Context, stockMutationList []entity.PossibleStockMutation) error
}

type stockMutationRepositoryPostgres struct {
	db DBTX
}

func NewStockMutationRepositoryPostgres(db *pgxpool.Pool) stockMutationRepositoryPostgres {
	return stockMutationRepositoryPostgres{
		db: db,
	}
}

func (r *stockMutationRepositoryPostgres) GetPossibleStockMutation(ctx context.Context, cartItems []entity.CartItemForCheckout) ([]entity.PossibleStockMutation, error) {
	alternatives := []entity.PossibleStockMutation{}
	args := []interface{}{}
	query := database.GetTwoClosestAvailableStockBase1
	query += ` WHERE `
	for i, cartItem := range cartItems {
		query += `cart_item_id = $` + strconv.Itoa(len(args)+1)
		args = append(args, cartItem.Id)
		if i != len(cartItems)-1 {
			query += ` OR `
		}
	}

	query += `),`
	query += database.GetTwoClosestAvailableStockBase2
	lockQuery := query + database.GetTwoClosestAvailableStockLock
	listQuery := query + database.GetTwoClosestAvailableStockList

	_, err := r.db.Exec(ctx, lockQuery, args...)
	if err != nil {
		return alternatives, err
	}
	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return alternatives, err
	}
	defer rows.Close()

	for rows.Next() {
		var alternative entity.PossibleStockMutation
		err := rows.Scan(&alternative.CartItemId, &alternative.CartQuantity, &alternative.DrugId, &alternative.OriginalPharmacyDrug, &alternative.OriginalPharmacy,
			&alternative.OriginalStock, &alternative.AlternativePharmacyDrug, &alternative.AlternativePharmacy, &alternative.AlternativeStock)
		if err != nil {
			return []entity.PossibleStockMutation{}, err
		}
		if alternative.OriginalStock < 0 {
			alternatives = append(alternatives, alternative)
		}
	}
	return alternatives, nil
}

func (r *stockMutationRepositoryPostgres) PostStockMutations(ctx context.Context, stockMutationList []entity.PossibleStockMutation) error {
	query := database.CreateStockMutations
	args := []interface{}{}
	for i, stockMutation := range stockMutationList {
		query += `($` + strconv.Itoa(len(args)+1) + `, $` + strconv.Itoa(len(args)+2) + `, $` + strconv.Itoa(len(args)+3) + 
		`, $` + strconv.Itoa(len(args)+4) + `, $` + strconv.Itoa(len(args)+5) + `)` 
		args = append(args, stockMutation.OriginalPharmacy)
		args = append(args, stockMutation.AlternativePharmacy)
		args = append(args, stockMutation.DrugId)
		args = append(args, stockMutation.AlternativeStock)
		args = append(args, 2)
		if i != len(stockMutationList) - 1 {
			query += `,`
		}
	}
	
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}