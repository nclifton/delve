package rabbit

// TODO rename these confusing things

type ErrWorkerMessageParse struct {
	message string
}

func NewErrWorkerMessageParse(message string) *ErrRetryWorkerMessage {
	return &ErrRetryWorkerMessage{
		message: message,
	}
}

func (e *ErrWorkerMessageParse) Error() string {
	return e.message
}

type ErrRetryWorkerMessage struct {
	message string
}

func NewErrRetryWorkerMessage(message string) *ErrRetryWorkerMessage {
	return &ErrRetryWorkerMessage{
		message: message,
	}
}

func (e *ErrRetryWorkerMessage) Error() string {
	return e.message
}
