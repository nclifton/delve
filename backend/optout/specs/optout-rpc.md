### RPC unsub.GetUnsubLink

Request

```
type struct GetUnubLinkParams {
  AccountID   string
  LongURL     string
  ShortURL    string
  Sender      string
}
```

Reply

```
type struct GetUnsubLinkReply {

}
```

### RPC unsub.OptOut

Request

```
type struct OptOutParams {
  AccountID string
  Sender string
}
```

Reply

```
type struct OptOutReply {}
```
