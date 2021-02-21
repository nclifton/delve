# MTMO Validator



```
valid.Validate(v, customValidators ...CustomValidator)
```

```
	if err = valid.Validate(v); err != nil {
		log.Println(err)
		r.WriteValidatorError(err)
		return errors.New("request was invalid")
	}
```

Accepts an interface but it must be a struct

Specify validation rules using struct tags

```
type WebhookCreatePOSTRequest struct {
	Event     string `json:"event" valid:"contains(link_hit|opt_out|sms_status|mms_status|sms_inbound|mms_inbound)"`
	Name      string `json:"name" valid:"length(2|100)"`
	URL       string `json:"url" valid:"webhook_url"`
	RateLimit int    `json:"rate_limit" valid:"range(0|10000)"`
}
```


### Built-In Available Rules

 - `required`:  
    - if the value is a string it cannot be blank (""), an empty string
	- if the value is an array it cannot be empty, length 0
	- if an array of strings, the strings in the array cannot be blank, an empty string
	- if not a string, the value or values if an array, cannot be the type's zero value
 - `url`:,
 - `email`: 
 - `integer`:
 - `alpha`:
 - `length`:
 - `rune_length`:
 - `range`:
 - `contains`:
    - Contains parameter is an array of strings. 
	- Validation is true if the value equals one of the strings in the params string array.
    - This is effectively a "one-of" validation and not a string contains string validation
 - `webhook_url`:


### Injecting a Custom Validator

```
	if err = valid.Validate(v,myCustomValidator("nope","this one")); err != nil {
		log.Println(err)
		r.WriteValidatorError(err)
		return errors.New("request was invalid")
	}
```

This adds a validator rule named `custom`
```
func myCustomValidator(message string, match string) CustomValidator {
	return CustomValidator{
		Name: "custom",
		Fn: func(i interface{}, parent interface{}, params []string) error {
			if i != match {
				return errors.New(message)
			}
			return nil
		},
	}
}
```
In the struct include the name of the validation rule on the field to be validated using your custom validation
```
type ValidateThis {
	Data string `valid:"required,custom"`
}
```