# 其他

使用中间件鉴权，handler 中进行颗粒度控制

- query： 前端请求
- vo：响应给前端
- bo: service <-> handler
- dto: service <-> service / repo <-> repo
- po: service <->repo

优先横向调用

````bash
psql -U p-d tiny_forum
```

```sql
DROP DATABASE tiny_forum;
CREATE DATABASE tiny_forum;
````

ubuntu

数据库更改

````bash
\c tiny_forum

sudo -u postgres psql

\c postgres -- 切换到默认的 postgres 数据库（或使用 template1）
DROP DATABASE IF EXISTS tiny_forum;

```bash
sudo -u postgres psql <<EOF
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'tiny_forum' AND pid <> pg_backend_pid();
\c postgres
DROP DATABASE IF EXISTS tiny_forum;
CREATE DATABASE tiny_forum OWNER tinyform;
\q
EOF
````

```bash
sudo nano /etc/postgresql/*/main/pg_hba.conf
```
