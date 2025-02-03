package chat

import (
	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/model/entity"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
)

func NewChatProvider(router *gin.RouterGroup, infra *config.Infra) {
	routerGroup := router.Group("chats")
	// data source/repository
	userQueryRepository := sqlRepo.NewSQLQueryRepository[*entity.User, entity.User](infra.GormPool)
	conversationCommandRepository := sqlRepo.NewSQLCommandRepository[*entity.Conversation, entity.Conversation](infra.GormPool)
	participantQueryRepository := sqlRepo.NewSQLQueryRepository[*entity.Participant, entity.Participant](infra.GormPool)
	participantCommandRepository := sqlRepo.NewSQLCommandRepository[*entity.Participant, entity.Participant](infra.GormPool)
	messageCommandRepository := sqlRepo.NewSQLCommandRepository[*entity.Message, entity.Message](infra.GormPool)
	messageQueryRepository := sqlRepo.NewSQLQueryRepository[*entity.Message, entity.Message](infra.GormPool)
	chatQueryRepository := sqlRepo.NewSQLQueryRepository[*entity.UserChat, entity.UserChat](infra.GormPool)
	// contact handler
	contactSrv := NewContactService(userQueryRepository)
	NewContactHandler(routerGroup, contactSrv)
	// chat (conversation) handler
	convSrv := NewConversationService(conversationCommandRepository, participantQueryRepository,
		participantCommandRepository, messageCommandRepository, messageQueryRepository, chatQueryRepository)
	NewConversationHandler(routerGroup, convSrv, infra.RabbitMQConsumer)
}
