package common

type ContextKey string

const contextKeyPrefix = "daochan context key "

func (c ContextKey) String() string {
	return contextKeyPrefix + string(c)
}

var (
	ContextKeyTraceID          = ContextKey("request id")
	ContextKeyRequestStartTime = ContextKey("request start time")
	ContextKeyRemoteAddress    = ContextKey("request remote address")
	ContextKeyUser             = ContextKey("user")
)
