package constant

type ContextKey string

const (
	ContextKeyClientIP  ContextKey = "client_ip"
	ContextKeyUserAgent ContextKey = "user_agent"
	ContextKeyReferer   ContextKey = "referer"
)
