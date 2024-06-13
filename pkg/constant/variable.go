package constant

type ContextKey string

var (
	JWTClaimsContextKey ContextKey
)

const (
	RedisRefreshTokenKeyPrefix = "refresh_token:"
)
