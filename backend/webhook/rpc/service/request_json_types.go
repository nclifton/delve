package service

type PublishStatusData struct {
	MMS_id            string `json:"mms_id"`
	SMS_id            string `json:"sms_id"`
	Message_ref       string `json:"message_ref"`
	Recipient         string `json:"recipient"`
	Sender            string `json:"sender"`
	Status            string `json:"status"`
	Status_updated_at string `json:"status_updated_at"`
}

type PublishMessageData struct {
	Type         string   `json:"type"`
	Id           string   `json:"id"`
	Recipient    string   `json:"recipient"`
	Sender       string   `json:"sender"`
	Subject      string   `json:"subject"`
	Message      string   `json:"message"`
	Content_urls []string `json:"content_urls"`
	Message_ref  string   `json:"message_ref"`
}
type PublishOptOutData struct {
	Source         string             `json:"source"`
	Contact_ref    string             `json:"contact_ref"`
	Timestamp      string             `json:"timestamp"`
	Source_message PublishMessageData `json:"source_message"`
}

type PublishMOData struct {
	MMS_id       string             `json:"mms_id"`
	SMS_id       string             `json:"sms_id"`
	Recipient    string             `json:"recipient"`
	Sender       string             `json:"sender"`
	Subject      string             `json:"subject"`
	Message      string             `json:"message"`
	Content_urls []string           `json:"content_urls"`
	Contact_ref  string             `json:"contact_ref"`
	Timestamp    string             `json:"timestamp"`
	Last_message PublishMessageData `json:"last_message"`
}

type PublishLinkHitData struct {
	URL            string             `json:"url"`
	Hits           int                `json:"hits"`
	Timestamp      string             `json:"timestamp"`
	Source_message PublishMessageData `json:"source_message"`
}
