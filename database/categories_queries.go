package database

const (
	FindAllCategories = `
		SELECT drug_category_id, drug_category_url, drug_category_name
		FROM drug_categories
		WHERE deleted_at IS NULL
		ORDER BY drug_category_name ASC
	`

	DeleteOneCategoryById = `
		UPDATE drug_categories
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE drug_category_id = $1
	`

	GetOneCategoryById = `
		SELECT drug_category_id, drug_category_url, drug_category_name
		FROM drug_categories
		WHERE drug_category_id = $1 
		AND deleted_at IS NULL
	`

	GetOneCategoryByName = `
		SELECT drug_category_id, drug_category_url, drug_category_name
		FROM drug_categories
		WHERE drug_category_name ILIKE $1 
		AND deleted_at IS NULL
	`

	PostOneCategoryQuery = `
		INSERT 
		INTO drug_categories (drug_category_name, drug_category_url)
		VALUES ($1, $2)
	`

	UpdateOneCategoryQuery = `
		UPDATE drug_categories
		SET
		drug_category_name = $1,
		drug_category_url = $2,
		updated_at = NOW()
		WHERE drug_category_id = $3
		AND deleted_at IS NULL
	`

	GetSimilarCategory = `
		SELECT drug_category_id, drug_category_url, drug_category_name
		FROM drug_categories
		WHERE drug_category_name ILIKE $1 
		AND drug_category_id != $2
		AND deleted_at IS NULL
	`
)	