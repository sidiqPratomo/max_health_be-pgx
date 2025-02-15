package database

const (
	GetOneDrugFormById = `
		SELECT form_id, form_name
		FROM drug_forms
		WHERE form_id = $1 
		AND deleted_at IS NULL
	`
)
