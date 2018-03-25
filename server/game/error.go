package criscross

const (
	UNKNOW_ERROR  = "UNKNOWN_ERROR"
	AUTH_ERROR    = "AUTH_ERROR"
	REG_ERROR     = "REG_ERROR"
	NOT_YOUR_STEP = "NOT_YOUR_STEP"
	VALUE_ERROR   = "VALUE_ERROR"
)

type gameError interface {
	error
	Code() string
}

type gameErrorImpl struct {
	msg, code string
}

func (err *gameErrorImpl) Error() string {
	return err.msg
}

func (err *gameErrorImpl) Code() string {
	return err.code
}

func NewGameError(code, msg string) gameError {
	return &gameErrorImpl{msg, code}
}
