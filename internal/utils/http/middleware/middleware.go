package middleware

var (
	tokenSecret string
)

func SetRequirement(secret string) {
	tokenSecret = secret
}
