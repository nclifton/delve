## GET Webhooks

Request

```
GET /v1/webhooks
```

Response

```
200 OK
Content-Type: application/json
{
  webhooks: [
    {
      id: string,
      account_id: string,
      event: string,
      name: string,
      url: string,
      ratelimit: int,
      created_at: string,
      updated_at: string,
    }
  ]
}
```

## Create Webhook

Request

```
POST /v1/webhooks
Content-Type: application/json
{
  event: string,
  name: string,
  url: string,
  ratelimit: int,
}
```

Response

```
201 Created
Content-Type: application/json
{
  id: string,
  account_id: string,
  event: string,
  name: string,
  url: string,
  ratelimit: int,
  created_at: string,
  updated_at: string,
}

```

## Update Webhook

Request

```
PUT /v1/webhooks
Content-Type: application/json
{
  event: string,
  name: string,
  url: string,
  ratelimit: int,
}
```

Response

```
200 OK
Content-Type: application/json
{
  id: string,
  account_id: string,
  event: string,
  name: string,
  url: string,
  ratelimit: int,
  created_at: string,
  updated_at: string,
}

```

## GET Single Webhook

Request

```
GET /v1/webhooks/<id>
```

Response

```
200 OK
Content-Type: application/json
{
  id: string,
  account_id: string,
  event: string,
  name: string,
  url: string,
  ratelimit: int,
  created_at: string,
  updated_at: string,
}

```

## DELETE Single Webhook

Request

```
DELETE /v1/webhooks/<id>
```

Response

```
200 OK
Content-Type: application/json
{
  id: string,
}
```
