# 常见问题

## 为什么无法打开后台管理页面？

前后端访问使用 jwt 鉴权，确保前后端 jwt 是相同的。请查看 `.env` 以及 `config/private.yaml` 文件。

## 如何设置默认管理员？

在 `config/private.yaml` 文件中，设置 `admin` 字段即可。

