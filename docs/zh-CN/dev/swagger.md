# Swagger 接口测试

## 先注册一个用户

```bash
http://localhost:8080/swagger/index.html#/%E9%AA%8C%E8%AF%81%E7%AE%A1%E7%90%86/post_auth_register
```
或者 使用命令行

```bash
curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d '{
"email": "admin@email.com",
"username": "admin",
"password": "adminadmin"
}'
```

## 登陆测试

```bash
curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{
"email": "test@email.com",
"password": "testtest"
}'
```

注册成功的响应：

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6InVzZXIiLCJpc3MiOiJ0aW55LWZvcnVtIiwiZXhwIjoxNzc2MjI2NTU3LCJpYXQiOjE3NzU5NjczNTd9.Btu72AoKu7dXQSo2uTYpFtKIt9-k8U44SKWEdYdAcro",
        "user": {
            "id": 3,
            "created_at": "2026-04-12T12:15:57.838887+08:00",
            "updated_at": "2026-04-12T12:15:57.838887+08:00",
            "deleted_at": null,
            "username": "admin",
            "email": "admin@email.com",
            "avatar": "https://api.dicebear.com/8.x/lorelei/svg?seed=admin",
            "bio": "",
            "role": "user",
            "score": 0,
            "is_active": true,
            "is_blocked": false,
            "last_login": null
        }
    }
}                                        
```
> 现在的角色还是 user（用户），为了测试方便，建议修改为 admin（管理员），需要在数据库中操作

```bash
psql -U postgres -d tiny_forum -h localhost -p 5432
# 先看一下所有用户
SELECT id, username, email, role, created_at FROM users;
# 修改角色
UPDATE users SET role = 'admin' WHERE username = 'admin';
# 再确认修改成功
SELECT id, username, email, role, created_at FROM users;
\q # 退出
```

然后使用测试意思登陆，获取到 token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{
"email": "admin@email.com",
"password": "adminadmin"
}'
```

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiaXNzIjoidGlueS1mb3J1bSIsImV4cCI6MTc3NjIyNjkzMywiaWF0IjoxNzc1OTY3NzMzfQ.e_J2oKS9gAC3ITrPQdJTkVlc_ujeiMvwVZL6DQJN6zw",
        "user": {
            "id": 3,
            "created_at": "2026-04-12T12:15:57.838887+08:00",
            "updated_at": "2026-04-12T12:22:13.308231+08:00",
            "deleted_at": null,
            "username": "admin",
            "email": "admin@email.com",
            "avatar": "https://api.dicebear.com/8.x/lorelei/svg?seed=admin",
            "bio": "",
            "role": "admin",
            "score": 0,
            "is_active": true,
            "is_blocked": false,
            "last_login": "2026-04-12T12:22:13.307013+08:00"
        }
    }
}
```

现在已经是 admin 了。


## 复制 token

输入到 swagger 的 Authorize 中，例如


```bash
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiaXNzIjoidGlueS1mb3J1bSIsImV4cCI6MTc3NjIyNjkzMywiaWF0IjoxNzc1OTY3NzMzfQ.e_J2oKS9gAC3ITrPQdJTkVlc_ujeiMvwVZL6DQJN6zw
```

可以在终端测试一下，访问管理员专用接口：

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiaXNzIjoidGlueS1mb3J1bSIsImV4cCI6MTc3NjIyNjkzMywiaWF0IjoxNzc1OTY3NzMzfQ.e_J2oKS9gAC3ITrPQdJTkVlc_ujeiMvwVZL6DQJN6zw"
```

```bash
$ curl -X GET "http://localhost:8080/api/v1/admin/posts?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN" \
  -v
```

输出：

```out
Note: Unnecessary use of -X or --request, GET is already inferred.
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /api/v1/admin/posts?page=1&page_size=20 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiaXNzIjoidGlueS1mb3J1bSIsImV4cCI6MTc3NjIyNjkzMywiaWF0IjoxNzc1OTY3NzMzfQ.e_J2oKS9gAC3ITrPQdJTkVlc_ujeiMvwVZL6DQJN6zw
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 12 Apr 2026 04:26:22 GMT
< Transfer-Encoding: chunked
< 
* Connection #0 to host localhost left intact
{"code":0,"message":"success","data":{"list":[{"id":4,"created_at":"2026-04-12T05:26:08.433528+08:00","updated_at":"2026-04-12T05:26:08.433528+08:00","deleted_at":null,"title":"封禁测试","content":"\u003cp\u003e账号一杯封禁\u003c/p\u003e","summary":"测试概述","cover":"","type":"post","status":"published","author_id":2,"view_count":4,"like_count":0,"pin_top":false,"author":{"id":2,"created_at":"2026-04-11T23:30:19.801771+08:00","updated_at":"2026-04-12T12:12:00.765913+08:00","deleted_at":null,"username":"test","email":"test@email.com","avatar":"https://api.dicebear.com/8.x/lorelei/svg?seed=test","bio":"","role":"user","score":30,"is_active":true,"is_blocked":false,"last_login":"2026-04-12T12:12:00.759127+08:00"},"board_id":0,"pin_in_board":false,"is_question":false,"board":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"name":"","slug":"","description":"","icon":"","cover":"","parent_id":null,"sort_order":0,"view_role":"","post_role":"","reply_role":"","post_count":0,"thread_count":0,"today_count":0}},{"id":3,"created_at":"2026-04-11T23:30:51.517235+08:00","updated_at":"2026-04-11T23:30:51.517235+08:00","deleted_at":null,"title":"fdsaf","content":"\u003cp\u003efdsaf\u003c/p\u003e","summary":"fdsa","cover":"","type":"post","status":"published","author_id":2,"view_count":1,"like_count":0,"pin_top":false,"author":{"id":2,"created_at":"2026-04-11T23:30:19.801771+08:00","updated_at":"2026-04-12T12:12:00.765913+08:00","deleted_at":null,"username":"test","email":"test@email.com","avatar":"https://api.dicebear.com/8.x/lorelei/svg?seed=test","bio":"","role":"user","score":30,"is_active":true,"is_blocked":false,"last_login":"2026-04-12T12:12:00.759127+08:00"},"board_id":0,"pin_in_board":false,"is_question":false,"board":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"name":"","slug":"","description":"","icon":"","cover":"","parent_id":null,"sort_order":0,"view_role":"","post_role":"","reply_role":"","post_count":0,"thread_count":0,"today_count":0}},{"id":2,"created_at":"2026-04-11T20:22:26.028212+08:00","updated_at":"2026-04-11T20:22:26.028212+08:00","deleted_at":null,"title":"标题测试","content":"\u003cp\u003e文章内容\u003c/p\u003e","summary":"文章概述","cover":"","type":"article","status":"published","author_id":1,"view_count":30,"like_count":0,"pin_top":false,"author":{"id":1,"created_at":"2026-04-11T09:18:19.160415+08:00","updated_at":"2026-04-11T23:22:06.63796+08:00","deleted_at":null,"username":"caoyang","email":"reggiesimons2cy@gmail.com","avatar":"https://api.dicebear.com/8.x/lorelei/svg?seed=caoyang","bio":"","role":"admin","score":31,"is_active":true,"is_blocked":false,"last_login":"2026-04-11T23:22:06.636098+08:00"},"board_id":0,"pin_in_board":false,"is_question":false,"board":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"name":"","slug":"","description":"","icon":"","cover":"","parent_id":null,"sort_order":0,"view_role":"","post_role":"","reply_role":"","post_count":0,"thread_count":0,"today_count":0}}],"total":3,"page":1,"page_size":20}}%                                                                   
```

没问题。

打开 http://localhost:8080/swagger/index.html#/%E7%AE%A1%E7%90%86%E6%8E%A5%E5%8F%A3/get_admin_posts 输入后进行测试即可。


