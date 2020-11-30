package types

type QueueMOParams struct {
	MessageID     string
	Message       string
	To            string
	From          string
	SARID         string
	SARPartNumber string
	SARParts      string
}

type ProcessMOParams struct {
	MessageID     string
	Message       string
	To            string
	From          string
	SARID         string
	SARPartNumber string
	SARParts      string
}
