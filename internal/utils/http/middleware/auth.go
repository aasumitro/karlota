package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/internal/utils/security"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		tokens := strings.Split(header, " ")
		if header == "" || len(tokens) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := tokens[1]
		claim, err := security.ExtractAndValidateJWT(tokenSecret, accessToken)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("id", claim.Payload["id"])
		ctx.Set("email", claim.Payload["email"])
		ctx.Set("display_name", claim.Payload["display_name"])
		ctx.Next()
	}
}

func SSEAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set header
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		// validate token
		token, err := ctx.Cookie("KATS")
		if err != nil || token == "" {
			ctx.SSEvent("error", http.StatusUnauthorized)
			ctx.Writer.Flush()
			return
		}
		token = strings.Trim(token, `"`)
		claim, err := security.ExtractAndValidateJWT(tokenSecret, token)
		if err != nil {
			ctx.SSEvent("error", http.StatusUnauthorized)
			ctx.Writer.Flush()
			return
		}
		// Set context items
		ctx.Set("id", claim.Payload["id"])
		ctx.Set("email", claim.Payload["email"])
		ctx.Set("display_name", claim.Payload["display_name"])
		ctx.Next()
	}
}

func WSAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// validate token
		token, err := ctx.Cookie("KATS")
		if err != nil || token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token = strings.Trim(token, `"`)
		claim, err := security.ExtractAndValidateJWT(tokenSecret, token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Set context items
		ctx.Set("id", claim.Payload["id"])
		ctx.Set("email", claim.Payload["email"])
		ctx.Set("display_name", claim.Payload["display_name"])
		ctx.Next()
	}
}
