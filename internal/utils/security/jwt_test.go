package security_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"karlota.aasumitro.id/internal/utils/security"
)

func TestJSONWebToken_ClaimJWTToken(t *testing.T) {
	type fields struct {
		Issuer    string
		SecretKey []byte
		Payload   map[string]interface{}
		IssuedAt  time.Time
		ExpiredAt time.Time
	}

	issuedAt, _ := time.Parse("",
		"2022-12-05 17:57:44.321843 +0800 WITA m=+25.737606459")
	expiredAt, _ := time.Parse("",
		"2022-12-06 17:57:44.321851 +0800 WITA m=+86425.737614876")

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "NEW JWT TEST SHOULD SUCCESS",
			fields: fields{
				Issuer: "BECOOP_TEST",
				Payload: map[string]interface{}{
					"data": "hello world",
				},
				SecretKey: []byte("123"),
				IssuedAt:  issuedAt,
				ExpiredAt: expiredAt,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJCRUNPT1BfVEVTVCIsImV4cCI6LTYyMTM1NTk2ODAwLCJpYXQiOi02MjEzNTU5NjgwMCwicGF5bG9hZCI6eyJkYXRhIjoiaGVsbG8gd29ybGQifX0.uImxm8cezx_7DBlXV0E3YJT4FwNUioZBRoWfYPs5H1g",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt := &security.JSONWebToken{
				Issuer:    tt.fields.Issuer,
				SecretKey: tt.fields.SecretKey,
				IssuedAt:  tt.fields.IssuedAt,
				ExpiredAt: tt.fields.ExpiredAt,
			}
			got, err := jwt.Claim(tt.fields.Payload)
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.AddCookie(&http.Cookie{Name: "access_token", Value: got})
			if data, err := req.Cookie("access_token"); err == nil {
				_, _ = security.ExtractAndValidateJWT(string(tt.fields.SecretKey), data.String())
			}
			if !tt.wantErr(t, err, "ClaimJWTToken()") {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
