### SMS Service -> SMS Send Queue

Queued Message

```
type Job struct {
	ID				  OID
  AccountID   OID
	Message 		string
	Recipient		string
  GSM         bool
	Sender			string
  AlarisUser  string
  AlarisPass  string
}
```

### HTTP Request to Alaris

Request

```
GET https://api.mtmo.io/api?username=<Job.AlarisUser>&password=<Job.AlarisPass>&ani=<Job.Sender>&dnis=<Job.Recipient>&message=<Job.Message>&command=submit&longMessageMode=split
```

Reply

```
HTTP/1.1 200 OK Content-Type: text/html; charset=UTF-8
{"message_id":"alss-a1b2c3d4-e5f67890"}
```

Bad Request Reply

```
HTTP/1.1 400 Bad Request Content-Type: text/html; charset=UTF-8
NO ROUTES
```
