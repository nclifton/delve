package client

type AlarisClientError struct {
	RetryAble bool
	Message   string
}

func NewError(message string, retry bool) *AlarisClientError {
	return &AlarisClientError{
		Message:   message,
		RetryAble: retry,
	}
}

func (e *AlarisClientError) Error() string {
	return e.Message
}
