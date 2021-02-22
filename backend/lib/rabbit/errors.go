package rabbit

// TODO rename these confusing things

type ErrWorkerMessageParse struct {
	message string
}

func NewErrWorkerMessageParse(message string) *ErrWorkerMessageParse {
	return &ErrWorkerMessageParse{
		message: message,
	}
}

func (e *ErrWorkerMessageParse) Error() string {
	return e.message
}

type ErrRetryWorkerMessage struct {
	message string
}

func (e *ErrRetryWorkerMessage) Error() string {
	return e.message
}

func NewErrRetryWorkerMessage(message string) *ErrRetryWorkerMessage {
	return &ErrRetryWorkerMessage{
		message: message,
	}
}

type ErrRequeueWorkerMessage struct {
	message string
}

func (e *ErrRequeueWorkerMessage) Error() string {
	return e.message
}

func NewErrRequeueWorkerMessage(message string) *ErrRequeueWorkerMessage {
	return &ErrRequeueWorkerMessage{
		message: message,
	}
}
