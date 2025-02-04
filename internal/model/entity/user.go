package entity

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/request"
	"karlota.aasumitro.id/internal/utils/security"
)

type User struct {
	ID               uint       `gorm:"primaryKey;column:id" json:"id"`
	DisplayName      string     `gorm:"column:display_name" json:"display_name"`
	Email            string     `gorm:"column:email" json:"email"`
	Password         string     `gorm:"column:password" json:"-"`
	EmailVerifiedAt  *time.Time `gorm:"column:email_verified_at" json:"-"`
	PasswordUpdateAt *time.Time `gorm:"column:password_updated_at" json:"-"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"-"`
	Online           bool       `gorm:"-" json:"online"`
	LastOnlineAt     *time.Time `gorm:"-" json:"last_online_at"`
}

func (u *User) ToResponse() *User {
	return &User{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Password:    u.Password,
		EmailVerifiedAt: func() *time.Time {
			if u.EmailVerifiedAt != nil {
				return u.EmailVerifiedAt
			}
			return nil
		}(),
		PasswordUpdateAt: func() *time.Time {
			if u.PasswordUpdateAt != nil {
				return u.PasswordUpdateAt
			}
			return nil
		}(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) FromResponse(user *User) interface{} {
	return &User{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Password:    user.Password,
		EmailVerifiedAt: func() *time.Time {
			if user.EmailVerifiedAt != nil {
				return user.EmailVerifiedAt
			}
			return nil
		}(),
		PasswordUpdateAt: func() *time.Time {
			if user.PasswordUpdateAt != nil {
				return user.PasswordUpdateAt
			}
			return nil
		}(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *User) MakeNewUser(body request.RegisterRequest) error {
	u.DisplayName = body.DisplayName
	u.Email = body.Email
	h := security.PasswordHashArgon2{Stored: "", Supplied: body.Password}
	pwd, err := h.MakePassword()
	if err != nil {
		return err
	}
	u.Password = pwd
	now := time.Now()
	u.PasswordUpdateAt = &now
	return nil
}

func (u *User) MakeNewPassword(newPassword string) error {
	h := security.PasswordHashArgon2{Stored: "", Supplied: newPassword}
	pwd, err := h.MakePassword()
	if err != nil {
		return err
	}
	u.Password = pwd
	now := time.Now()
	u.PasswordUpdateAt = &now
	u.UpdatedAt = now
	return nil
}

func (u *User) ValidatePassword(pwd string) error {
	h := security.PasswordHashArgon2{Stored: u.Password, Supplied: pwd}
	isValid, err := h.ComparePassword()
	if err != nil {
		return err
	}
	if !isValid {
		return common.ErrInvalidPassword
	}
	return nil
}

func (u *User) GenerateAccessToken(
	j *security.JSONWebToken,
	secret string,
	refresh bool,
) (map[string]interface{}, error) {
	var accessToken, refreshToken map[string]interface{}
	issuedAt := time.Now()
	// Function to generate a token
	generateToken := func(
		t *security.JSONWebToken,
		d time.Duration,
		secret string,
	) (map[string]interface{}, error) {
		localToken := *t
		localToken.SecretKey = []byte(secret)
		localToken.IssuedAt = issuedAt
		expiry := localToken.IssuedAt.Add(d)
		localToken.ExpiredAt = expiry
		jwt, err := localToken.Claim(map[string]interface{}{
			"id":           u.ID,
			"email":        u.Email,
			"display_name": u.DisplayName,
		})
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"str": jwt,
			"dtm": expiry,
		}, nil
	}
	eg := new(errgroup.Group)
	// Goroutine for generating the access token
	eg.Go(func() error {
		var err error
		accessToken, err = generateToken(j, time.Minute*30, secret)
		return err
	})
	// Goroutine for generating the refresh token
	eg.Go(func() error {
		if refresh {
			return nil
		}
		var err error
		refreshToken, err = generateToken(j, time.Hour*24, secret)
		return err
	})
	// Wait for both goroutines to finish and check for any error
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	// Return the generated tokens
	return map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (u *User) GenerateExchangeToken(
	j *security.JSONWebToken,
	secret string,
) (string, error) {
	j.SecretKey = []byte(secret)
	j.IssuedAt = time.Now()
	j.ExpiredAt = time.Now().Add(time.Hour * 2)
	return j.Claim(map[string]interface{}{
		"email":        u.Email,
		"display_name": u.DisplayName,
	})
}

func (u *User) MakePasswordLinkNotificationPayload(
	exchangeToken string,
) ([]byte, error) {
	// {url}/reset?exchange_token={exchangeToken}
	resetPasswordLink := fmt.Sprintf(
		"%s/reset-password?exchange_token=%s",
		config.GetBaseURL(), exchangeToken)
	notifyBody := fmt.Sprintf(common.ForgotPasswordTemplate,
		u.DisplayName, resetPasswordLink)
	return (&Notify{
		Target: u.Email, Body: notifyBody,
		Subject: "Reset Password Link"}).Encode()
}

func (u *User) MakeVerifyEmailNotificationPayload(
	exchangeToken string,
) ([]byte, error) {
	// {url}/verify-email?exchange_token={exchangeToken}
	verifyEmailLink := fmt.Sprintf(
		"%s/api/v1/auth/verify-email?exchange_token=%s",
		config.GetBaseURL(), exchangeToken)
	notifyBody := fmt.Sprintf(common.VerifyEmailTemplate,
		u.DisplayName, verifyEmailLink)
	return (&Notify{
		Target: u.Email, Body: notifyBody,
		Subject: "Verify Email Link"}).Encode()
}
