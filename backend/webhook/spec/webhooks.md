### Common

Return all hooks in some consistent envelope like so:

```
{
	event: 'sms_inbound|contact_update|<source_event>',
	data: {
		webhook data
	}
}
```

The json request to send a message as defined in API

```
{
	recipient: '61400001701',
	sender: '61404123123',
	subject: 'hello',
	message: 'here is a QR code for your collection',
	content_urls: [ 'https://link.to/file.png' ],
	shorten_urls: false,
	message_ref: 'msg-123',
}
```

### MMS Status Webhook

- need to go through and clean up terminology used (not DLR anymore)

http post request

```
{
	event: 'mms_status',
	data: {
		mms_id: '7zxczaa789c789zxcf',
		message_ref: 'msg-123',
		recipient: '61400001701',
		sender: '61404123123',
		status: 'delivered',
		status_updated_at: '2020-07-16T23:32:38.114Z',
	}
}

```

### MMS Inbound Webhook

http post request

```
{
	event: 'mms_inbound',
	data: {
		mms_id: '6hewczaa71235c7zxscf',
		recipient: '61404123123',
		sender: '61400001701',
		subject: '',
		message: 'thanks! this one looks better',
		content_urls: [ 'https://sendsei.cdn/image.png' ],
		contact_ref: 'customer-123',
		timestamp: '2020-07-16T23:32:38.114Z',
		last_message: {
			type: 'mms',
			id: '7zxczaa789c789zxcf',
			recipient: '61400001701',
			sender: '61404123123',
			subject: 'hello',
			message: 'here is a QR code for your collection',
			content_urls: [ 'https://link.to/file.png' ],
			message_ref: 'msg-123',
		}
	}
}

```

### SMS Inbound Webhook

http post request

```
{
	event: 'sms_inbound',
	data: {
		sms_id: '345njk632k45n2khjlsw',
		recipient: '61404123123',
		sender: '61400001701',
		message: 'thanks! looks great',
		contact_ref: 'customer-123',
		timestamp: '2020-07-16T23:32:38.114Z',
		last_message: {
			type: 'mms',
			id: '7zxczaa789c789zxcf',
			recipient: '61400001701',
			sender: '61404123123',
			subject: 'hello',
			message: 'here is a QR code for your collection',
			content_urls: [ 'https://link.to/file.png' ],
			message_ref: 'msg-123',
		}
	}
}

```

### Unsub Webhook

http POST REQUEST NOT RESPONSE

```
{
	event: 'unsubscribe',
	data: {
		source: 'sms_inbound', // link_hit, etc
		contact_ref: 'customer-123',
		timestamp: '2020-07-16T23:32:38.114Z',
		source_message: {
			type: 'sms',
			id: '34623kl45m23sadfasd',
			recipient: '61400001701',
			sender: '61404123123',
			message: 'sign up for free! Opt-out reply STOP',
			message_ref: 'msg-124',
		}
	}
}
```

### Link Hit Webhook

```
{
	event: 'link_hit',
	data: {
		url: 'https://apple.com',
		hits: 1,
		timestamp: '2020-07-16T23:32:38.114Z',
		source_message: {
			type: 'sms',
			id: '34623kl45m23sadfasd',
			recipient: '61400001701',
			sender: '61404123123',
			message: 'buy stuff here to elevate your social status getinfo.at/dsdhs42i',
			message_ref: 'msg-125',
		}
	}
}
```
