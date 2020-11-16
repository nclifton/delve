### Inbound DLR

POST /v1/dlr/alaris

```
Content-Type: application/x-www-form-urlencoded

msgid=<MessageID>&state=<Status>&reasoncode=<ReasonCode>&to=<Recipient>&time=2020-02-05Z00%3A09%3A00&mcc=<Mcc>&mnc=<Mnc>

```

Possible status values:

```
ENROUTE Message is in routing stage
SENT Message is delivered to the SMSC (Short Message Service Centre)
DELIVRD Message is delivered to the subscriber
EXPIRED Message storage period expired
DELETED Message was deleted
UNDELIV Message cannot be delivered
ACCEPTD Message was accepted by SMSC
REJECTD Message was rejected by SMSC
UNKNOWN Unknown message status
```

### Inbound MO

POST /v1/mo/alaris

```
Content-Type: application/x-www-form-urlencoded

message=<Message>&to=<Recipient>&from=<Sender>&msgid=<MessageID>&sarId=<sarId>&sarPartNumber=<SarPartNumber>&SarParts=<SarParts>
```
