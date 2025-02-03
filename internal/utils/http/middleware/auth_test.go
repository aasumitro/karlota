package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/utils/http/middleware"
	"karlota.aasumitro.id/internal/utils/security"
)

func TestAuthMiddleware(t *testing.T) {
	viper.Reset()
	viper.SetConfigFile("../../../../.example.env")
	viper.SetConfigType("dotenv")
	config.Load().With(context.TODO())
	middleware.SetRequirement(config.GetAuthSecret())
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJQT1NCRV9URVNUIiwiZXhwIjotNjIxMzU1OTY4MDAsImlhdCI6LTYyMTM1NTk2ODAwLCJwYXlsb2FkIjp7ImRhdGEiOiJoZWxsbyB3b3JsZCJ9fQ.-_tfeKKhqSRP2H_pVg4f_spkX_Z1Lo1nuiu09OFFvO0"

	router := gin.Default()
	router.Use(middleware.Auth())

	t.Run("ERROR AT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ERROR PARSE TOKEN WITH CLAIM", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("SUCCESS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		jwt := security.JSONWebToken{
			Issuer:    "MIDDLEWARE_TEST",
			SecretKey: []byte(config.GetAuthSecret()),
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(1 * time.Minute),
		}
		at, err := jwt.Claim(nil)
		assert.Nil(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", at))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
