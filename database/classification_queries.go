package database

const (
	GetOneDrugClassficationById = `
		SELECT classification_id, classification_name
		FROM drug_classifications
		WHERE classification_id = $1 
		AND deleted_at IS NULL
	`
)
