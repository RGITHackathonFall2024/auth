package user

type ErrNoSuchUser struct {
}

func (e *ErrNoSuchUser) Error() string {
	return "no such user"
}
