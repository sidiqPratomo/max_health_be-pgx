package database

const (
	FindAccountByEmailQuery = `
		SELECT a.account_id, a.email, a.password, a.role_id, r.role_name, a.account_name, a.profile_picture, a.verified_at
		FROM accounts a
		JOIN roles r ON a.role_id =  r.role_id
		WHERE a.email ILIKE $1
		AND a.deleted_at 
		IS 
		NULL
	`

	FindOneAccountPasswordByIdQuery = `
		SELECT password
		FROM accounts
		WHERE account_id = $1
		AND deleted_at IS NULL
	`

	PostOneAccountQuery = `
		INSERT 
		INTO accounts (email, password, role_id, account_name)
		VALUES ($1, $2, $3, $4)
		RETURNING account_id
	`

	PostOneVerifiedAccountQuery = `
		INSERT 
		INTO accounts (email, password, role_id, account_name, profile_picture, verified_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING account_id
	`

	QueryUpdatePasswordOneAccount = `
		UPDATE accounts 
		SET password = $1,
		verified_at = NOW(),
		updated_at = NOW()
		WHERE account_id = $2
	`

	UpdateDataOneAccount = `
		UPDATE accounts 
		SET account_name = $1,
		password = $2,
		profile_picture = $3,
		updated_at = NOW()
		WHERE account_id = $4
	`

	UpdateNameAndProfilePictureOneAccount = `
		UPDATE accounts 
		SET account_name = $1,
		profile_picture = $2,
		updated_at = NOW()
		WHERE account_id = $3
	`

	FindOneAccountByIdQuery = `
		SELECT account_name, password, email, profile_picture
		FROM accounts
		WHERE account_id = $1
		AND deleted_at IS NULL
	`

	DeleteOneAccountByIdQuery = `
		UPDATE accounts
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE account_id = $1
	`
)
