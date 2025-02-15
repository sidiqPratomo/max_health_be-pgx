package database

const (
	GetAllDrugFormQuery = `
		SELECT form_id, form_name
		FROM drug_forms
		WHERE deleted_at IS NULL
	`
)