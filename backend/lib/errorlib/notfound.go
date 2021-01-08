package errorlib

type NotFoundErr struct {
	Message string
}

func (e NotFoundErr) Error() string {
	return e.Message
}
