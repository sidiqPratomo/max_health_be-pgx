package database

const (
	PostOneUserQuery = `
		INSERT 
		INTO users (account_id)
		VALUES ($1)
	`

	FindUserByAccountIdQuery = `
		SELECT u.user_id, u.gender_id, g.gender_name, u.date_of_birth
		FROM users u
		JOIN accounts a
		ON u.account_id = a.account_id
		LEFT JOIN genders g
		ON u.gender_id = g.gender_id
		WHERE u.account_id = $1
		AND a.verified_at IS NOT NULL
		AND (u.deleted_at IS NULL AND a.deleted_at IS NULL)
	`
	
	UpdateDataOneUser = `
		UPDATE users 
		SET gender_id = $1,
		date_of_birth = $2,
		updated_at = NOW()
		WHERE account_id = $3
	`
)
