package errors

type NotFoundErr struct {
	Err error
}

func (e NotFoundErr) Error() string {
	return e.Err.Error()
}
