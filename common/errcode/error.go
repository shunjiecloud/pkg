package ec

type Error struct {
	Result   string
	Msg      string
	HttpCode int
	IsWarn   bool
}

func (e Error) Error() string {
	return e.Result
}

func (e Error) Message() string {
	return e.Msg
}
