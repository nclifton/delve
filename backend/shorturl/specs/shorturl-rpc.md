### RPC shorturl.GenerateShortUrls

request

```
type struct GenerateShortUrlsParam {
  AccountID   string
  Message     string
}
```

reply

```
type struct GenerateShortUrlsReply {
  Message string
}
```

### RPC shorturl.LinkHit

request

```
type struct LinkHitParams {
  ShortID string
  UA string
}
```

response

```
type struct LinkHitReply {
}
```
