### RPC sms.send

Request

```
type SMSSendParams struct {
	AccountID       OID
	Message         string
	Recipient       string
	Sender          string
	Country         string
  MessageRef      string
}
```

Reply

```
type SMS struct {
	ID        OID       `json:"_id" bson:"_id"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	SMSData   `bson:",inline"`
	SMSStatus `bson:",inline"`
	DLR       DLR `json:"dlr" bson:"dlr"`
}

type SMSData struct {
  MessageID   string
  MessageRef  string
	Country     string          `json:"country" bson:"country"`
	Message     string          `json:"message" bson:"message"`
	SMSCount    int             `json:"sms_count" bson:"sms_count"`
	GSM         bool            `json:"gsm" bson:"gsm"`
	Recipient   string          `json:"recipient" bson:"recipient"`
	Sender      string          `json:"sender" bson:"sender"`
	AccountID   OID             `json:"account_id" bson:"account_id"`
}

type SMSStatus struct {
	Status string `json:"status" bson:"status" valid:"required"`
}

type SMSSendReply struct {
	SMS   *SMS
	Parts int
}
```

### RPC sms.MOProcess

Request

```
type MOProcessParams struct {
  MessageID     string
	Message       string `json:"message"`
	Source        string `json:"source"`
	Dest          string `json:"dest"`
  SarId         int
  SarPartNumber int
  SarParts      int
}
```

Reply

```
type MOProcessReply struct {
}
```

### RPC sms.DLRProcess

Request

```
type DLRProcessParams struct {
	MessageID   string
  State       string
  ReasonCode  string
  To          string
  Time        time.Date
  Mcc         string
  Mnc         string
}
```

Reply

```
type DLRProcessReply struct {
}
```
