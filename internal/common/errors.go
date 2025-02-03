package common

import "errors"

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidExchangeToken  = errors.New("invalid exchange token")
	ErrEmailExist            = errors.New("email already exists")
	ErrOldPasswordNotValid   = errors.New("old password is not valid")
	ErrEncodeNotifyPayload   = errors.New("encode payload")
	ErrGenerateExchangeToken = errors.New("exchange token generation failed")
	ErrDeleteGroup           = errors.New("you are not allowed to delete this conversation")
)
