package database

const (
	PostOnePrescriptionDrugQuery = `
		INSERT INTO prescription_drugs (prescription_id, drug_id, quantity, note)
		VALUES ($1, $2, $3, $4)
	`

	GetAllPrescriptionDrugQuery = `
		SELECT pd.prescription_drug_id, d.drug_id, d.drug_name, d.image, d.is_active, pd.quantity, pd.note
		FROM prescription_drugs pd 
		JOIN drugs d ON d.drug_id = pd.drug_id
		WHERE pd.prescription_id = $1
	`
)
