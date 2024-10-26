package auth

type ErrInvalidHash struct {
}

func (e *ErrInvalidHash) Error() string {
	return "invalid hash"
}

type ErrMissingToken struct {
}

func (e *ErrMissingToken) Error() string {
	return "missing token"
}

type ErrInvalidToken struct {
}

func (e *ErrInvalidToken) Error() string {
	return "invalid token"
}
