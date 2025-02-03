package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJSONWebToken interface {
	Claim(payload interface{}) (string, error)
}

type JSONWebTokenClaim struct {
	jwt.RegisteredClaims
	Payload map[string]interface{} `json:"payload"`
}

type JSONWebToken struct {
	Issuer    string
	SecretKey []byte
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func (j *JSONWebToken) Claim(payload map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JSONWebTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			IssuedAt:  &jwt.NumericDate{Time: j.IssuedAt},
			ExpiresAt: &jwt.NumericDate{Time: j.ExpiredAt},
		},
		Payload: payload,
	})

	return token.SignedString(j.SecretKey)
}

func ExtractAndValidateJWT(
	secret string, token string,
) (claim *JSONWebTokenClaim, err error) {
	var parseToken *jwt.Token
	var ok bool

	parseToken, err = jwt.ParseWithClaims(token, &JSONWebTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if parseToken == nil {
		return nil, errors.New("library error")
	}
	if err != nil && !parseToken.Valid {
		return nil, err
	}

	if claim, ok = parseToken.Claims.(*JSONWebTokenClaim); !ok {
		return nil, errors.New("invalid claim token")
	}

	return claim, nil
}
