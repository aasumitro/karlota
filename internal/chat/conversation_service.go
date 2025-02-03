package chat

import (
	"context"
	"fmt"
	"sort"
	"time"

	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/entity"
	"karlota.aasumitro.id/internal/model/request"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
)

type conversationService struct {
	conversationCommandRepository sqlRepo.ISQLCommandRepository[*entity.Conversation, entity.Conversation]
	participantQueryRepository    sqlRepo.ISQLQueryRepository[*entity.Participant, entity.Participant]
	participantCommandRepository  sqlRepo.ISQLCommandRepository[*entity.Participant, entity.Participant]
	messageCommandRepository      sqlRepo.ISQLCommandRepository[*entity.Message, entity.Message]
	messageQueryRepository        sqlRepo.ISQLQueryRepository[*entity.Message, entity.Message]
	chatQueryRepository           sqlRepo.ISQLQueryRepository[*entity.UserChat, entity.UserChat]
}

func (service *conversationService) New(
	ctx context.Context,
	form request.NewGroupRequest,
) (uint, error) {
	// create a new conversation
	newConversation := &entity.Conversation{}
	if err := newConversation.MakeNewGroupConversation(form); err != nil {
		return 0, err
	}
	if err := service.conversationCommandRepository.Insert(ctx, newConversation); err != nil {
		return 0, err
	}
	// add participants
	participants := []*entity.Participant{{
		ConvID: newConversation.ID,
		UserID: form.CreatorID,
		Role: func() string {
			if newConversation.Type == common.ConversationTypePrivate {
				return common.ParticipantRoleNone
			}
			return common.ParticipantRoleAdmin
		}(),
		CreatedAt: time.Now(),
	}}
	for _, member := range form.Members {
		participants = append(participants, &entity.Participant{
			ConvID: newConversation.ID,
			UserID: uint(member.ID),
			Role: func() string {
				if newConversation.Type == common.ConversationTypePrivate {
					return common.ParticipantRoleNone
				}
				return common.ParticipantRoleMember
			}(),
			CreatedAt: time.Now(),
		})
	}
	if err := service.participantCommandRepository.InsertMany(ctx, participants); err != nil {
		return 0, err
	}
	// add a message
	if form.Message != "" {
		if err := service.messageCommandRepository.Insert(ctx, &entity.Message{
			ConvID:    newConversation.ID,
			UserID:    form.CreatorID,
			Type:      common.MessageTypeText,
			Content:   form.Message,
			CreatedAt: time.Now(),
		}); err != nil {
			return 0, err
		}
	}
	return newConversation.ID, nil
}

func (service *conversationService) Chats(
	ctx context.Context,
	userID uint,
) (interface{}, error) {
	limit, offset := -1, -1 // we will load all
	subQuery := fmt.Sprintf("SELECT conversation_id FROM participants WHERE user_id = %d", userID)
	return service.chatQueryRepository.FindWithLimit(ctx, limit, offset,
		sqlRepo.WithInSubquery("conversation_id", subQuery))
}

func (service *conversationService) Messages(
	ctx context.Context,
	conversationID uint,
) (interface{}, error) {
	limit, offset := -1, -1
	messages, err := service.messageQueryRepository.FindWithLimit(
		ctx, limit, offset, sqlRepo.Equal("conversation_id", conversationID))
	if err != nil {
		return nil, err
	}
	// Ascending order: oldest first
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})
	return messages, nil
}

func (service *conversationService) NewTextMessage(
	ctx context.Context,
	payload *request.WebsocketPayload,
) (*entity.Message, error) {
	data := &entity.Message{
		ConvID:    payload.ConversationID,
		UserID:    payload.UserID,
		Type:      common.MessageTypeText,
		Content:   payload.TextMessage,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		ReadAt:    nil,
	}
	if err := service.messageCommandRepository.Insert(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (service *conversationService) RequestDeleteGroup(
	ctx context.Context,
	payload *request.WebsocketPayload,
) error {
	users, err := service.participantQueryRepository.FindWithLimit(ctx, 1, 0,
		sqlRepo.And(sqlRepo.Equal("user_id", payload.UserID),
			sqlRepo.Equal("conversation_id", payload.ConversationID)))
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return common.ErrAccountNotFound
	}
	user := users[0]
	if user.Role != common.ParticipantRoleAdmin {
		return common.ErrDeleteGroup
	}
	return service.cleanupGroup(ctx, payload)
}

func (service *conversationService) RequestLeaveGroup(
	ctx context.Context,
	payload *request.WebsocketPayload,
) (*entity.Message, error) {
	// get total group member
	totalMember, err := service.participantQueryRepository.Count(ctx,
		sqlRepo.Equal("conversation_id", payload.ConversationID))
	if err != nil {
		return nil, err
	}
	// if participant only 1 then cleanup all data
	if totalMember <= 1 {
		if err := service.cleanupGroup(ctx, payload); err != nil {
			return nil, err
		}
		return nil, nil
	}
	// otherwise just make the user leave the group
	return service.leaveGroup(ctx, payload)
}

func (service *conversationService) cleanupGroup(
	ctx context.Context,
	payload *request.WebsocketPayload,
) error {
	if err := service.messageCommandRepository.DeleteWithSpec(ctx,
		&entity.Message{}, sqlRepo.Equal("conversation_id", payload.ConversationID),
	); err != nil {
		return err
	}
	if err := service.participantCommandRepository.DeleteWithSpec(ctx,
		&entity.Participant{}, sqlRepo.Equal("conversation_id", payload.ConversationID),
	); err != nil {
		return err
	}
	return service.conversationCommandRepository.DeleteByID(ctx, payload.ConversationID)
}

func (service *conversationService) leaveGroup(
	ctx context.Context,
	payload *request.WebsocketPayload,
) (*entity.Message, error) {
	if err := service.participantCommandRepository.DeleteWithSpec(ctx,
		&entity.Participant{}, sqlRepo.And(
			sqlRepo.Equal("conversation_id", payload.ConversationID),
			sqlRepo.Equal("user_id", payload.UserID)),
	); err != nil {
		return nil, err
	}
	msg := &entity.Message{
		ConvID:    payload.ConversationID,
		UserID:    payload.UserID,
		Type:      common.MessageTypeText,
		Content:   fmt.Sprintf("%s Leave Group", payload.UserDN),
		CreatedAt: time.Now(),
	}
	if err := service.messageCommandRepository.Insert(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (service *conversationService) RequestLeaveConversation(
	ctx context.Context,
	payload *request.WebsocketPayload,
) error {
	//TODO implement me
	panic("implement me")
}

func (service *conversationService) cleanupConversation() error {
	return nil
}

func (service *conversationService) leaveConversation() error {
	return nil
}

func NewConversationService(
	conversationCommandRepository sqlRepo.ISQLCommandRepository[*entity.Conversation, entity.Conversation],
	participantQueryRepository sqlRepo.ISQLQueryRepository[*entity.Participant, entity.Participant],
	participantCommandRepository sqlRepo.ISQLCommandRepository[*entity.Participant, entity.Participant],
	messageCommandRepository sqlRepo.ISQLCommandRepository[*entity.Message, entity.Message],
	messageQueryRepository sqlRepo.ISQLQueryRepository[*entity.Message, entity.Message],
	chatQueryRepository sqlRepo.ISQLQueryRepository[*entity.UserChat, entity.UserChat],
) IConversationService {
	return &conversationService{
		conversationCommandRepository: conversationCommandRepository,
		participantQueryRepository:    participantQueryRepository,
		participantCommandRepository:  participantCommandRepository,
		messageCommandRepository:      messageCommandRepository,
		messageQueryRepository:        messageQueryRepository,
		chatQueryRepository:           chatQueryRepository,
	}
}
