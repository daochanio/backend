package common

type ContextKey string

const contextKeyPrefix = "etheralley context key "

func (c ContextKey) String() string {
	return contextKeyPrefix + string(c)
}

var (
	ContextKeyTraceID          = ContextKey("request id")
	ContextKeyRequestStartTime = ContextKey("request start time")
	ContextKeyAddress          = ContextKey("address")
)
