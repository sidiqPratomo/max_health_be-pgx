package database

const (
	PostOneResetPasswordTokenQuery = `
		INSERT 
		INTO reset_password_tokens (account_id, reset_token, expired_at)
		VALUES ($1, $2, NOW() + INTERVAL '60 minutes')
	`

	InvalidateTokensQuery = `
		UPDATE reset_password_tokens
		SET expired_at = NOW(), updated_at = NOW()
		WHERE account_id = $1
		AND expired_at > NOW()
	`

	SelectOneResetPasswordTokenByAccountIdAndTokenQuery = `
		SELECT 
		reset_password_token_id,
		expired_at
		FROM reset_password_tokens
		WHERE account_id = $1 AND reset_token = $2
		AND deleted_at IS NULL
	`

	UpdateExpiredAtOneResetPasswordTokenQuery = `
		UPDATE reset_password_tokens 
		SET expired_at = NOW(),
		updated_at = NOW()
		WHERE reset_password_token_id = $1
	`
)
