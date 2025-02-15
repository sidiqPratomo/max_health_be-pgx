package database

const (
	PostOneRefreshTokenQuery = `
		INSERT 
		INTO refresh_tokens (account_id, refresh_token, expired_at)
		VALUES ($1, $2, NOW() + INTERVAL '1 day')
	`

	InvalidateRefreshTokensQuery = `
		UPDATE refresh_tokens
		SET expired_at = NOW(), updated_at = NOW()
		WHERE account_id = $1
	`

	FindOneRefreshTokenQuery = `
		SELECT refresh_token
		FROM refresh_tokens
		WHERE refresh_token = $1 
		AND expired_at > NOW()
	`
)