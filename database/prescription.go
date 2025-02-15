package database

const (
	CreateOnePrescriptionQuery = `
		INSERT INTO prescriptions (user_account_id, doctor_account_id)
		VALUES ($1, $2)
		RETURNING prescription_id
	`

	GetPrescriptionByIdQuery = `
		SELECT user_account_id, doctor_account_id, redeemed_at, ordered_at
		FROM prescriptions
		WHERE prescription_id = $1
		AND deleted_at IS NULL
	`

	SetPrescriptionRedeemedNowQuery = `
		UPDATE prescriptions
		SET redeemed_at = NOW(), updated_at = NOW()
		WHERE prescription_id = $1 AND deleted_at IS NULL
	`

	GetPrescriptionListByUserAccountIdQuery = `
		SELECT p.prescription_id, p.user_account_id, a1.account_name, p.doctor_account_id, a2.account_name, p.redeemed_at, p.ordered_at, p.created_at
		FROM prescriptions p
		JOIN accounts a1 ON a1.account_id = p.user_account_id
		JOIN accounts a2 ON a2.account_id = p.doctor_account_id
		WHERE p.user_account_id = $1 AND p.deleted_at IS NULL AND redeemed_at IS NOT NULL
		ORDER BY p.created_at DESC
		LIMIT $2
		OFFSET $3
	`

	GetPrescriptionListByUserAccountIdTotalPageQuery = `
		SELECT COUNT (*)
		FROM prescriptions p
		WHERE p.user_account_id = $1 AND p.deleted_at IS NULL
	`

	SetPrescriptionOrderedAtNowQuery = `
		UPDATE prescriptions
		SET ordered_at = NOW(), updated_at = NOW()
		WHERE prescription_id = $1
	`
)
