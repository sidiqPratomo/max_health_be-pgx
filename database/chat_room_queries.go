package database

const (
	CreateOneRoomQuery = `
		INSERT INTO chat_rooms (user_account_id, doctor_account_id)
		VALUES ($1, $2)
		RETURNING chat_room_id
	`

	StartChatQuery = `
		UPDATE chat_rooms
		SET expired_at = NOW() + INTERVAL '30 minutes'
		WHERE chat_room_id = $1
		AND doctor_account_id = $2
		AND expired_at IS NULL
		AND deleted_at IS NULL
	`

	FindActiveChatRoomQuery = `
		SELECT chat_room_id, user_account_id, doctor_account_id, expired_at
		FROM chat_rooms
		WHERE
		user_account_id = $1 AND doctor_account_id = $2
		AND (expired_at > NOW() OR expired_at IS NULL)
		AND deleted_at IS NULL
	`

	FindChatRoomByIdQuery = `
		SELECT user_account_id, doctor_account_id, expired_at
		FROM chat_rooms
		WHERE chat_room_id = $1 
		AND deleted_at IS NULL
	`

	GetAllChatRoomPreviewQuery = `
		SELECT cr.chat_room_id, a.account_name, a.profile_picture, c.chat_message, attachment_format, c.attachment_url, cr.expired_at,
			CASE 
				WHEN c.chat_message IS NULL THEN cr.created_at
				ELSE c.created_at
			END
		FROM chat_rooms cr 
		JOIN accounts a  ON a.account_id  =
			CASE
				WHEN cr.user_account_id = $1 THEN cr.doctor_account_id
				WHEN cr.doctor_account_id = $1 THEN cr.user_account_id
			END
		LEFT JOIN chats c ON c.chat_id = 
			(SELECT c2.chat_id
			FROM chats c2
			WHERE c2.chat_room_id = cr.chat_room_id 
			ORDER BY created_at DESC limit 1)
		WHERE cr.deleted_at IS NULL AND 
			CASE
				WHEN cr.user_account_id = $1 THEN cr.user_account_id = $1
				WHEN cr.doctor_account_id = $1 THEN cr.doctor_account_id = $1
			END
		ORDER BY created_at DESC
	`

	DoctorGetChatRequestQuery = `
		SELECT cr.chat_room_id, a.account_name, a.profile_picture, c.chat_message, c.attachment_format, c.attachment_url ,
			CASE 
				WHEN c.chat_message IS NULL THEN cr.created_at
				ELSE c.created_at
			END
		FROM chat_rooms cr 
		JOIN accounts a  ON a.account_id  = cr.user_account_id
		LEFT JOIN chats c ON c.chat_id = 
			(SELECT c2.chat_id
			FROM chats c2
			WHERE c2.chat_room_id = cr.chat_room_id 
			ORDER BY created_at DESC limit 1)
		WHERE cr.deleted_at IS NULL AND expired_at IS NULL AND doctor_account_id = $1
		ORDER BY created_at DESC
	`

	CloseChatRoomQuery = `
		UPDATE chat_rooms
		SET expired_at = NOW(), updated_at = NOW()
		WHERE chat_room_id = $1
	`
)
