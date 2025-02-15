package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sidiqPratomo/max-health-backend/database"
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type ChatRoomRepository interface {
	CreateOneRoom(ctx context.Context, userAccountId, doctorAccountId int64) (*int64, error)
	StartChat(ctx context.Context, roomId, doctorAccountId int64) error
	FindActiveChatRoom(ctx context.Context, userAccountId, doctorAccountId int64) (*entity.ChatRoom, error)
	FindChatRoomById(ctx context.Context, chatRoomId int64) (*entity.ChatRoom, error)
	GetAllChatRoomPreview(ctx context.Context, accountId int64, role string) ([]entity.ChatRoomPreview, error)
	DoctorGetChatRequest(ctx context.Context, accountId int64) ([]entity.ChatRoomPreview, error)
	CloseChatRoom(ctx context.Context, roomId int64) error
}

type chatRoomRepositoryPostgres struct {
	db DBTX
}

func NewChatRoomRepositoryPostgres(db *pgxpool.Pool) chatRoomRepositoryPostgres {
	return chatRoomRepositoryPostgres{
		db: db,
	}
}

func (r *chatRoomRepositoryPostgres) CreateOneRoom(ctx context.Context, userAccountId, doctorAccountId int64) (*int64, error) {
	var chatRoomId int64

	err := r.db.QueryRow(ctx, database.CreateOneRoomQuery, userAccountId, doctorAccountId).Scan(&chatRoomId)
	if err != nil {
		return nil, err
	}

	return &chatRoomId, nil
}

func (r *chatRoomRepositoryPostgres) StartChat(ctx context.Context, roomId, doctorAccountId int64) error {
	_, err := r.db.Exec(ctx, database.StartChatQuery, roomId, doctorAccountId)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRoomRepositoryPostgres) FindActiveChatRoom(ctx context.Context, userAccountId, doctorAccountId int64) (*entity.ChatRoom, error) {
	var chatRoom entity.ChatRoom

	err := r.db.QueryRow(ctx, database.FindActiveChatRoomQuery, userAccountId, doctorAccountId).Scan(&chatRoom.Id, &chatRoom.UserAccountId, &chatRoom.DoctorAccountId, &chatRoom.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &chatRoom, nil
}

func (r *chatRoomRepositoryPostgres) FindChatRoomById(ctx context.Context, chatRoomId int64) (*entity.ChatRoom, error) {
	var chatRoom entity.ChatRoom

	err := r.db.QueryRow(ctx, database.FindChatRoomByIdQuery, chatRoomId).Scan(&chatRoom.UserAccountId, &chatRoom.DoctorAccountId, &chatRoom.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	chatRoom.Id = chatRoomId

	return &chatRoom, nil
}

func (r *chatRoomRepositoryPostgres) GetAllChatRoomPreview(ctx context.Context, accountId int64, role string) ([]entity.ChatRoomPreview, error) {
	rows, err := r.db.Query(ctx, database.GetAllChatRoomPreviewQuery, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatRoomPreviewList []entity.ChatRoomPreview

	for rows.Next() {
		var chatRoomPreview entity.ChatRoomPreview

		err := rows.Scan(
			&chatRoomPreview.Id,
			&chatRoomPreview.ParticipantName,
			&chatRoomPreview.ParticipantPictureUrl,
			&chatRoomPreview.LastChat.Message,
			&chatRoomPreview.LastChat.Attachment.Format,
			&chatRoomPreview.LastChat.Attachment.Url,
			&chatRoomPreview.ExpiredAt,
			&chatRoomPreview.LastChat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		chatRoomPreviewList = append(chatRoomPreviewList, chatRoomPreview)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return chatRoomPreviewList, nil
}

func (r *chatRoomRepositoryPostgres) DoctorGetChatRequest(ctx context.Context, accountId int64) ([]entity.ChatRoomPreview, error) {
	rows, err := r.db.Query(ctx, database.DoctorGetChatRequestQuery, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatRoomPreviewList []entity.ChatRoomPreview

	for rows.Next() {
		var chatRoomPreview entity.ChatRoomPreview

		err := rows.Scan(
			&chatRoomPreview.Id,
			&chatRoomPreview.ParticipantName,
			&chatRoomPreview.ParticipantPictureUrl,
			&chatRoomPreview.LastChat.Message,
			&chatRoomPreview.LastChat.Attachment.Format,
			&chatRoomPreview.LastChat.Attachment.Url,
			&chatRoomPreview.LastChat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		chatRoomPreviewList = append(chatRoomPreviewList, chatRoomPreview)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return chatRoomPreviewList, nil
}

func (r *chatRoomRepositoryPostgres) CloseChatRoom(ctx context.Context, roomId int64) error {
	_, err := r.db.Exec(ctx, database.CloseChatRoomQuery, roomId)
	if err != nil {
		return err
	}

	return nil
}
