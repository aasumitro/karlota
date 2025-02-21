package chat

import (
	"context"

	"karlota.aasumitro.id/internal/model/entity"
	"karlota.aasumitro.id/internal/model/request"
)

type IContactService interface {
	List(ctx context.Context, requestBy string) ([]*entity.User, error)
}

type IConversationService interface {
	New(ctx context.Context, form request.NewGroupRequest) (uint, error)
	Chats(ctx context.Context, userID uint) (interface{}, error)
	Messages(ctx context.Context, conversationID uint) (interface{}, error)
	NewTextMessage(ctx context.Context, payload *request.WebsocketPayload) (*entity.Message, error)

	RequestDeleteGroup(ctx context.Context, payload *request.WebsocketPayload) error
	RequestLeaveGroup(ctx context.Context, payload *request.WebsocketPayload) (*entity.Message, error)
	RequestLeaveConversation(ctx context.Context, payload *request.WebsocketPayload) error

	AddMessageQueue(ctx context.Context, userID uint, msg *entity.Message)
	LoadQueue(ctx context.Context, userID uint) ([]*entity.Queue, error)
	DeleteQueue(ctx context.Context, q *entity.Queue)
}
