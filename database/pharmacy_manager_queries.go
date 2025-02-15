package database

const (
	PostOnePharmacyManagerQuery = `
		INSERT INTO pharmacy_managers (account_id)
		VALUES ($1)
	`

	FindAllPharmacyManagers = `
		SELECT pm.pharmacy_manager_id, a.account_id, a.email, a.account_name, a.profile_picture
		FROM pharmacy_managers pm
		JOIN accounts a ON a.account_id = pm.account_id
		WHERE pm.deleted_at IS NULL
	`

	GetOnePharmacyManagerByIdQuery = `
		SELECT pm.pharmacy_manager_id, a.account_id, a.profile_picture
		FROM pharmacy_managers pm
		JOIN accounts a ON a.account_id = pm.account_id
		WHERE pm.pharmacy_manager_id = $1 
		AND pm.deleted_at IS NULL
	`

	DeleteOnePharmacyManagerByIdQuery = `
		UPDATE pharmacy_managers
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE pharmacy_manager_id = $1
	`

	GetOnePharmacyManagerByAccountIdQuery = `
		SELECT pm.pharmacy_manager_id, a.account_id, a.profile_picture
		FROM pharmacy_managers pm
		JOIN accounts a ON a.account_id = pm.account_id
		WHERE a.account_id = $1 
		AND pm.deleted_at IS NULL
	`

	GetOnePharmacyManagerByPharmacyCourierIdQuery = `
		SELECT pm.pharmacy_manager_id, a.account_id, a.profile_picture
		FROM pharmacy_managers pm
		JOIN accounts a ON a.account_id = pm.account_id
		JOIN pharmacies p ON p.pharmacy_manager_id = pm.pharmacy_manager_id
		JOIN pharmacy_couriers pc ON pc.pharmacy_id = p.pharmacy_id
		WHERE pc.pharmacy_courier_id = $1 
		AND pm.deleted_at IS NULL
	`
)
