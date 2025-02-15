package database

const (
	PostOneCodeQuery = `
		INSERT 
		INTO verification_codes (account_id, code, expired_at)
		VALUES ($1, $2, NOW() + INTERVAL '60 minutes')
	`

	InvalidateCodesQuery = `
		UPDATE verification_codes
		SET expired_at = NOW(), updated_at = NOW()
		WHERE account_id = $1
	`

	SelectOneVerificationCodeByCodeQuery = `
		SELECT 
		verification_code_id,
		expired_at
		FROM verification_codes
		WHERE account_id = $1 AND code = $2
		AND deleted_at IS NULL
	`

	UpdateExpiredAtOneVerificationCodeQuery = `
		UPDATE verification_codes 
		SET expired_at = NOW(),
		updated_at = NOW()
		WHERE verification_code_id = $1
	`
)
