# MTMO Validator



```
valid.Validate(v)
```

```
	if err = valid.Validate(v); err != nil {
		log.Println(err)
		r.WriteValidatorError(err)
		return errors.New("request was invalid")
	}
```

accepts an interface but must be a struct

```
type WebhookCreatePOSTRequest struct {
	Event     string `json:"event" valid:"contains(link_hit|opt_out|sms_status|mms_status|sms_inbound|mms_inbound)"`
	Name      string `json:"name" valid:"length(2|100)"`
	URL       string `json:"url" valid:"webhook_url"`
	RateLimit int    `json:"rate_limit" valid:"range(0|10000)"`
}
```