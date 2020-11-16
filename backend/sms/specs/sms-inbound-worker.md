### SMS MO Inbound API -> SMS Inbound Queue

Queued Message

```
type Job struct {
  MessageID     string
	Message       string `json:"message"`
	Source        string `json:"source"`
	Dest          string `json:"dest"`
  SarId         int
  SarPartNumber int
  SarParts      int
}
```

### SMS DLR Inbound API -> SMS DLR Queue

Queued Message

```
type Job struct {
	MessageID   string
  State       string
  ReasonCode  string
  To          string
  Time        time.Date
  Mcc         string
  Mnc         string
}
```
