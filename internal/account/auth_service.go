package account

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/entity"
	"karlota.aasumitro.id/internal/model/request"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
	"karlota.aasumitro.id/internal/utils/security"
)

type authService struct {
	sqlQueryRepository   sqlRepo.ISQLQueryRepository[*entity.User, entity.User]
	sqlCommandRepository sqlRepo.ISQLCommandRepository[*entity.User, entity.User]
	rabbitMQClient       *amqp.Connection
	jwtToken             *security.JSONWebToken
}

func (service *authService) SignIn(
	ctx context.Context,
	body request.LoginRequest,
) (map[string]interface{}, error) {
	// Check Email Exist
	user, err := service.userByEmail(ctx, body.Email)
	if err != nil {
		return nil, err
	}
	// validate password
	if err := user.ValidatePassword(body.Password); err != nil {
		return nil, err
	}
	// Generate Tokens (AccessToken & RefreshToken)
	tokens, err := user.GenerateAccessToken(service.jwtToken,
		config.GetAuthSecret(), false)
	if err != nil {
		return nil, err
	}
	tokens["user"] = user
	return tokens, nil
}

func (service *authService) RequestForgotLink(
	ctx context.Context,
	body request.ForgotPasswordRequest,
) (interface{}, error) {
	// Check Email Exist
	user, err := service.userByEmail(ctx, body.Email)
	if err != nil {
		return nil, err
	}
	// generate unique exchange token
	exchangeToken, err := user.GenerateExchangeToken(
		service.jwtToken, config.GetForgotPwdSecret())
	if err != nil {
		return nil, common.ErrGenerateExchangeToken
	}
	// Send Event to RabbitMQ to Send Email
	notifyPayload, err := user.MakePasswordLinkNotificationPayload(exchangeToken)
	if err != nil {
		return nil, common.ErrEncodeNotifyPayload
	}
	if err := service.publishEvent(ctx, common.EmailEventTopic, notifyPayload); err != nil {
		return nil, err
	}
	// Return Message
	return "we have sent you an email to reset your password", nil
}

func (service *authService) SetNewPassword(
	ctx context.Context,
	body request.ResetPasswordRequest,
) error {
	// validate exchange token
	data, err := security.ExtractAndValidateJWT(
		config.GetForgotPwdSecret(), body.ExchangeToken)
	if err != nil {
		return common.ErrInvalidExchangeToken
	}
	// get user data
	user, err := service.userByEmail(ctx,
		data.Payload["email"].(string))
	if err != nil {
		return err
	}
	// generate new password hash & save it
	if err := user.MakeNewPassword(body.NewPassword); err != nil {
		return err
	}
	return service.sqlCommandRepository.Update(ctx, user)
}

func (service *authService) SignUp(
	ctx context.Context,
	body request.RegisterRequest,
) (map[string]interface{}, error) {
	// Check Email Exist
	user, err := service.userByEmail(ctx, body.Email)
	if err != nil && !errors.Is(err, common.ErrAccountNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, common.ErrEmailExist
	}
	// Create new user account
	newUser := &entity.User{}
	if err := newUser.MakeNewUser(body); err != nil {
		return nil, err
	}
	if err := service.sqlCommandRepository.Insert(ctx, newUser); err != nil {
		return nil, err
	}
	// generate unique exchange token and id for email verification
	exchangeToken, err := newUser.GenerateExchangeToken(
		service.jwtToken, config.GetVerifyAccountSecret())
	if err != nil {
		return nil, common.ErrGenerateExchangeToken
	}
	// Send Event to RabbitMQ to Send Email
	notifyPayload, err := newUser.MakeVerifyEmailNotificationPayload(exchangeToken)
	if err != nil {
		return nil, common.ErrEncodeNotifyPayload
	}
	if err := service.publishEvent(ctx, common.EmailEventTopic, notifyPayload); err != nil {
		return nil, err
	}
	// Generate JWT (AccessToken & RefreshToken)
	tokens, err := newUser.GenerateAccessToken(service.jwtToken,
		config.GetAuthSecret(), false)
	if err != nil {
		return nil, err
	}
	tokens["user"] = newUser
	return tokens, nil
}

func (service *authService) RefreshToken(
	ctx context.Context,
	userRefreshToken string,
) (map[string]interface{}, error) {
	// validate user refresh token
	data, err := security.ExtractAndValidateJWT(
		config.GetAuthSecret(), userRefreshToken)
	if err != nil {
		return nil, err
	}
	// get user data
	user, err := service.userByEmail(ctx,
		data.Payload["email"].(string))
	if err != nil {
		return nil, err
	}
	// generate new token
	tokens, err := user.GenerateAccessToken(service.jwtToken,
		config.GetAuthSecret(), true)
	if err != nil {
		return nil, err
	}
	tokens["user"] = user
	return tokens, nil
}

func (service *authService) VerifyEmail(
	ctx context.Context,
	exchangeToken string,
) error {
	// verify exchange token
	data, err := security.ExtractAndValidateJWT(config.GetVerifyAccountSecret(), exchangeToken)
	if err != nil {
		return common.ErrInvalidExchangeToken
	}
	// get user data
	user, err := service.userByEmail(ctx,
		data.Payload["email"].(string))
	if err != nil {
		return err
	}
	// save user data
	now := time.Now()
	user.EmailVerifiedAt = &now
	user.UpdatedAt = now
	return service.sqlCommandRepository.Update(ctx, user)
}

func (service *authService) publishEvent(
	ctx context.Context,
	topic string, body []byte,
) error {
	ch, err := service.rabbitMQClient.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()
	q, err := ch.QueueDeclare(topic, false,
		false, false, false, nil)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(ctx,
		"", q.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (service *authService) userByEmail(
	ctx context.Context,
	email string,
) (*entity.User, error) {
	users, err := service.sqlQueryRepository.FindWithLimit(
		ctx, 1, 0, sqlRepo.Equal("email", email))
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, common.ErrAccountNotFound
	}
	return users[0], nil
}

func NewAuthService(
	sqlQueryRepository sqlRepo.ISQLQueryRepository[*entity.User, entity.User],
	sqlCommandRepository sqlRepo.ISQLCommandRepository[*entity.User, entity.User],
	rabbitMQClient *amqp.Connection,
	jwtToken *security.JSONWebToken,
) IAuthService {
	return &authService{
		sqlQueryRepository:   sqlQueryRepository,
		sqlCommandRepository: sqlCommandRepository,
		rabbitMQClient:       rabbitMQClient,
		jwtToken:             jwtToken,
	}
}
