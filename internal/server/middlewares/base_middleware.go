package middlewares

type (
	BaseMiddleware struct {
		log logger
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}
)

func NewBase(log logger) *BaseMiddleware {
	return &BaseMiddleware{
		log: log,
	}
}
