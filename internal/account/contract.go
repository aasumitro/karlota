package account

import (
	"context"

	"karlota.aasumitro.id/internal/model/request"
)

type IAccountService interface {
	Profile(ctx context.Context, id uint) (interface{}, error)
	UpdateDisplayName(ctx context.Context, id uint, displayName string) (interface{}, error)
	UpdatePassword(ctx context.Context, id uint, oldPwd, newPwd string) (interface{}, error)
}

type IAuthService interface {
	SignIn(
		ctx context.Context,
		body request.LoginRequest,
	) (map[string]interface{}, error)
	RequestForgotLink(
		ctx context.Context,
		body request.ForgotPasswordRequest,
	) (interface{}, error)
	SetNewPassword(
		ctx context.Context,
		body request.ResetPasswordRequest,
	) error
	SignUp(
		ctx context.Context,
		body request.RegisterRequest,
	) (map[string]interface{}, error)
	RefreshToken(
		ctx context.Context,
		userRefreshToken string,
	) (map[string]interface{}, error)
	VerifyEmail(
		ctx context.Context,
		exchangeToken string,
	) error
}
