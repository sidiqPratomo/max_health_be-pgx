package database

const (
	PostOneChatQuery = `
		INSERT INTO chats (chat_room_id, sender_account_id, chat_message, attachment_format, attachment_url, prescription_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING chat_id, created_at
	`

	GetAllChatQuery = `
		SELECT c.chat_id, c.sender_account_id, c.chat_message, c.attachment_format, c.attachment_url, c.prescription_id, p.redeemed_at, c.created_at
		FROM chats c
		LEFT JOIN prescriptions p ON p.prescription_id = c.prescription_id
		WHERE chat_room_id = $1
		AND c.deleted_at IS NULL
	`
)
