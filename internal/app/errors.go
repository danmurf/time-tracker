package app

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEventNotFound      = Error("event not found")
	ErrTaskAlreadyStarted = Error("task already started")
)
