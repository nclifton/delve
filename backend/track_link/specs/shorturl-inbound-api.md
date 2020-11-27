### Short URL Link Hit

```
GET /<uniqueid>
```

Response

```
HTTP/1.1 301 Moved Permanently
Location: <shorturl.long_url>
```

Response if id not found

```
HTTP/1.1 400 Not Found
Content-Type: text/html

<insert link not found HTML here>
```
