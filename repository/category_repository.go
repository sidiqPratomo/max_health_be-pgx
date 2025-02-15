package repository

import (
	"context"
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type CategoryRepository interface {
	FindAllCategories(ctx context.Context) ([]entity.DrugCategory, error)
	DeleteOneCategoryById(ctx context.Context, categoryId int64) error
	FindOneCategoryById(ctx context.Context, categoryId int64) (*entity.DrugCategory, error)
	FindOneCategoryByName(ctx context.Context, name string) (*entity.DrugCategory, error)
	PostOneCategory(ctx context.Context, category entity.DrugCategory) error
	UpdateOneCategoryById(ctx context.Context, category entity.DrugCategory) error
	FindSimilarCategory(ctx context.Context, category entity.DrugCategory) (*entity.DrugCategory, error)
}

type categoryRepositoryPostgres struct {
	db DBTX
}

func NewCategoryRepositoryPostgres(db *sql.DB) categoryRepositoryPostgres {
	return categoryRepositoryPostgres{
		db: db,
	}
}

func (r *categoryRepositoryPostgres) FindAllCategories(ctx context.Context) ([]entity.DrugCategory, error) {
	query := database.FindAllCategories

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []entity.DrugCategory{}
	for rows.Next() {
		var category entity.DrugCategory
		err := rows.Scan(
			&category.Id,
			&category.Url,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
		
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepositoryPostgres) DeleteOneCategoryById(ctx context.Context, categoryId int64) error {
	_, err := r.db.ExecContext(ctx, database.DeleteOneCategoryById, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepositoryPostgres) FindOneCategoryById(ctx context.Context, categoryId int64) (*entity.DrugCategory, error) {
	var category entity.DrugCategory

	err := r.db.QueryRowContext(ctx, database.GetOneCategoryById, categoryId).Scan(&category.Id, &category.Url, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &category, nil
}

func (r *categoryRepositoryPostgres) FindOneCategoryByName(ctx context.Context, name string) (*entity.DrugCategory, error) {
	var category entity.DrugCategory

	err := r.db.QueryRowContext(ctx, database.GetOneCategoryByName, name).Scan(&category.Id, &category.Name, &category.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &category, nil
}

func (r *categoryRepositoryPostgres) PostOneCategory(ctx context.Context, category entity.DrugCategory) error {
	_, err := r.db.ExecContext(ctx, database.PostOneCategoryQuery, category.Name, category.Url)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryPostgres) UpdateOneCategoryById(ctx context.Context, category entity.DrugCategory) error {
	_, err := r.db.ExecContext(ctx, database.UpdateOneCategoryQuery, category.Name, category.Url, category.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryPostgres) FindSimilarCategory(ctx context.Context, category entity.DrugCategory) (*entity.DrugCategory, error) {
	var oldCategory entity.DrugCategory

	err := r.db.QueryRowContext(ctx, database.GetSimilarCategory, category.Name, category.Id).Scan(&oldCategory.Id, &oldCategory.Name, &oldCategory.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &oldCategory, nil
}