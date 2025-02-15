package repository

import (
	"context"
	"database/sql"

	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type ChatRepository interface {
	PostOneChat(ctx context.Context, chatRequest entity.Chat) (*int64, string, error)
	GetAllChat(ctx context.Context, roomId int64) ([]entity.Chat, error)
}

type chatRepositoryPostgres struct {
	db DBTX
}

func NewChatRepositoryPostgres(db *sql.DB) chatRepositoryPostgres {
	return chatRepositoryPostgres{
		db: db,
	}
}

func (r *chatRepositoryPostgres) PostOneChat(ctx context.Context, chatRequest entity.Chat) (*int64, string, error) {
	var chatId int64
	var createdAt string

	err := r.db.QueryRowContext(ctx, database.PostOneChatQuery, chatRequest.RoomId, chatRequest.SenderAccountId, chatRequest.Message, chatRequest.Attachment.Format, chatRequest.Attachment.Url, chatRequest.Prescription.Id).Scan(&chatId, &createdAt)
	if err != nil {
		return nil, "", err
	}

	return &chatId, createdAt, err
}

func (r *chatRepositoryPostgres) GetAllChat(ctx context.Context, roomId int64) ([]entity.Chat, error) {
	rows, err := r.db.QueryContext(ctx, database.GetAllChatQuery, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatList []entity.Chat

	for rows.Next() {
		var chat entity.Chat

		err := rows.Scan(
			&chat.Id,
			&chat.SenderAccountId,
			&chat.Message,
			&chat.Attachment.Format,
			&chat.Attachment.Url,
			&chat.Prescription.Id,
			&chat.Prescription.RedeemedAt,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		chat.RoomId = roomId

		chatList = append(chatList, chat)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return chatList, nil
}
