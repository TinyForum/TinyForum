# 用户验证信息

## 用户注册

### 请求

地址：

请求体：

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6InVzZXIiLCJyb2xlIjoidXNlciIsImlzcyI6InRpbnktZm9ydW0iLCJleHAiOjE3Nzg1NzU2NDEsImlhdCI6MTc3ODQ4OTI0MSwianRpIjoiODRjYmRhOGQtNDMzYS00ZDJhLWFmYTctOGVlMjAwYWM5ODAyIn0.bhTPqKTrtK7JqPY76CDQEOIYXi33SWbqn_B7GBDOErg",
    "user": {
      "id": 3,
      "created_at": "2026-05-11T16:47:21.984598+08:00",
      "updated_at": "2026-05-11T16:47:21.984598+08:00",
      "deleted_at": null,
      "username": "user",
      "email": "user@test.com",
      "password": "$2a$10$Y8hgWtV/d6fCzM/C4yJUBuPm5Qmv3nXUd9UXzwD8Td1JJwJJW8fVa",
      "avatar": "https://api.dicebear.com/8.x/lorelei/svg?seed=user",
      "bio": "",
      "role": "user",
      "score": 0,
      "is_active": true,
      "is_blocked": false,
      "last_login": null,
      "invited_by_id": null
    }
  },
  "timestamp": 1778489241
}
```





## 用户登陆

### 请求

地址：`/api/v1/auth/login`

请求体：

```json
{
  "email": "user_name@test.com",
  "password": "password"
}
```



### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z",
      "username": "admin",
      "avatar": "",
      "bio": "",
      "role": "super_admin",
      "score": 0,
      "is_active": false,
      "is_blocked": false
    }
  },
  "timestamp": 1778488814
}
```

