## Get Webhook List

Request

```
GET /v1/webhook
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

## Create a Webhook

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

## Update a Webhook

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

## Get a Single Webhook

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

## Delete a Single Webhook

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
