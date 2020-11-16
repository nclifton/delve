### SMS Send

Request

```
POST /v1/sms
Content-Type: application/json
{
	"message": "string",
	"message_ref": "string",
	"recipient": "string",
	"sender": "string",
	"country": "string"
}
```

Response

```
200 OK
Content-Type: application/json
{
  sms: {
    "message_id": "string",
    "account_id": "string",
    "message_ref": "string",
    "updated": "time.Time",
    "created": "time.Time",
    "country": "string",
    "message": "string",
    "sms_count": "int",
    "is_gsm": "bool",
    "recipient": "string",
    "sender": "string",
    "status": "string"
  },
  parts: "int"
}
```
