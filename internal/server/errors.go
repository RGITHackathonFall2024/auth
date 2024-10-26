package server

type ErrNoServerInContext struct {
}

func (e *ErrNoServerInContext) Error() string {
	return "no server in context"
}
