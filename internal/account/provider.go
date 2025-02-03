package account

import (
	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/model/entity"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
	"karlota.aasumitro.id/internal/utils/http/middleware"
	"karlota.aasumitro.id/internal/utils/security"
)

func NewAccountProvider(
	router *gin.RouterGroup,
	infra *config.Infra,
) {
	// data source/repository
	sqlQueryRepository := sqlRepo.NewSQLQueryRepository[*entity.User, entity.User](infra.GormPool)
	sqlCommandRepository := sqlRepo.NewSQLCommandRepository[*entity.User, entity.User](infra.GormPool)
	// auth handler
	authSrv := NewAuthService(sqlQueryRepository, sqlCommandRepository, infra.RabbitMQPublisher,
		&security.JSONWebToken{Issuer: config.GetServerInfo()})
	authRouterGroup := router.Group("/auth")
	NewAuthHandler(authRouterGroup, authSrv)
	// account handler
	accountSrv := NewAccountService(sqlQueryRepository, sqlCommandRepository)
	accountRouterGroup := router.Group("/account").
		Use(middleware.Auth())
	NewAccountHandler(accountRouterGroup, accountSrv)
}
