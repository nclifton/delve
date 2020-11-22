package msg

type MessageSpec struct {
	Queue        string
	Exchange     string
	ExchangeType string
	RouteKey     string
}
