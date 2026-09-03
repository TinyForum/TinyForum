# API Testing

Test the API endpoints using Swagger UI or curl.

## Swagger UI

Start the backend and open:

```
http://localhost:8080/swagger/index.html
```

## curl Examples

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"Test1234"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"Test1234"}'
```

For more examples, see the [Chinese Swagger Guide](/zh-CN/dev/swagger).
