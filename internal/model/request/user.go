package request

import (
	"context"

	"github.com/golodash/galidator/v2"
)

type (
	LoginRequest struct {
		ClientSessionID string `json:"-" form:"-"`
		Email           string `json:"email" form:"email"`
		Password        string `json:"password" form:"password"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" form:"email"`
	}

	ResetPasswordRequest struct {
		ExchangeToken   string `json:"exchange_token" form:"exchange_token"`
		NewPassword     string `json:"new_password" form:"new_password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	RegisterRequest struct {
		ClientSessionID string `json:"-" form:"-"`
		DisplayName     string `json:"display_name" form:"display_name"`
		Email           string `json:"email" form:"email"`
		Password        string `json:"password" form:"password"`
	}

	UpdateAccountRequest struct {
		DisplayName string `json:"display_name" form:"display_name"`
	}

	UpdatePasswordRequest struct {
		NewPassword string `json:"new_password" form:"new_password"`
		OldPassword string `json:"old_password" form:"old_password"`
	}
)

const minR, maxR = 5.0, 16.0

func (r *LoginRequest) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"Email":    g.R("email").Required().Email(),
		"Password": g.R("password").Required(),
	}).Validate(context.TODO(), r)
}

func (r *ForgotPasswordRequest) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"Email": g.R("email").Required().Email(),
	}).Validate(context.TODO(), r)
}

func (r *ResetPasswordRequest) Validate() interface{} {
	g := galidator.New()
	if r.NewPassword != r.ConfirmPassword {
		return map[string]string{
			"confirm_password": "new_password and confirm_password does not match",
		}
	}
	return g.ComplexValidator(galidator.Rules{
		"ExchangeToken":   g.R("exchange_token").Required(),
		"NewPassword":     g.R("new_password").Required().Password(),
		"ConfirmPassword": g.R("confirm_password").Required().Password(),
	}).Validate(context.TODO(), r)
}

func (r *RegisterRequest) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"DisplayName": g.R("display_name").Required().Min(minR).Max(maxR),
		"Email":       g.R("email").Required().Email(),
		"Password":    g.R("password").Required().Password(),
	}).Validate(context.TODO(), r)
}

func (r *UpdateAccountRequest) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"DisplayName": g.R("display_name").Required().Min(minR).Max(maxR),
	}).Validate(context.TODO(), r)
}

func (r *UpdatePasswordRequest) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"OldPassword": g.R("old_password").Required(),
		"NewPassword": g.R("new_password").Required().Password(),
	}).Validate(context.TODO(), r)
}
