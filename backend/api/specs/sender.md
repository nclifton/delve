## Get Sender List

Request

```
GET /v1/sender
```

Response

```
200 OK
Content-Type: application/json
{
  senders: [
    {
      "id": "12345",
      "account_id": "23456",
      "address": "61434237636",
      "mms_provider_key": "mgage",
      "channels": ["mms", "sms"],
      "country": "AU",
      "comment": "blah",
      "created_at": "2021-01-12T04:45:46.718261Z",
      "updated_at": "2021-01-12T04:45:46.718261Z"
    },
    {
      "id": "54321",
      "account_id": "65432",
      "address": "TOILET",
      "mms_provider_key": "optus",
      "channels": ["mms", "sms"],
      "country": "PH",
      "comment": "blah",
      "created_at": "2021-01-12T04:45:46.718261Z",
      "updated_at": "2021-01-12T04:45:46.718261Z"
    }
 ]
}
```
