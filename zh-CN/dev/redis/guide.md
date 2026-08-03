# 使用说明

后端使用 redis 进行限流，以避免用户灌水。


## 查找所有限流 key（仅测试环境）

```bash
KEYS rl:*
```

## 查看当前窗口内的请求数量

```bash
ZCARD rl:1:create_post
```

## 查看所有请求记录（score 和 member）

```bash
ZRANGE rl:1:create_post 0 -1 WITHSCORES
```

## 重置限流

```bash
DEL rl:1:create_post
```