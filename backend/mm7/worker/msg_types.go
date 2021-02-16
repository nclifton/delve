package worker

const (
	FakeProviderKey  = "fake"
	OptusProviderKey = "optus"
	MgageProviderKey = "mgage"

	QueueNameSubmitFake  = "fake"
	QueueNameSubmitOptus = "optus"
	QueueNameSubmitMgage = "mgage.submit"

	QueueNameDLDRFake = "dldr.fake"
)

type SubmitMessage struct {
	ID          string
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
}
