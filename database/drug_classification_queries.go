package database

const (
	GetAllDrugClassificationQuery = `
		SELECT classification_id, classification_name
		FROM drug_classifications
		WHERE deleted_at IS NULL
	`
)