# ── 构建 ────────────────────────────────────────────────────
# go.mod 声明 go 1.26.2；Docker Hub 当前稳定最新镜像为 golang:1.24
# GOTOOLCHAIN=local+auto 允许工具链在 go.mod 要求更高版本时自动下载
# 同时建议将 go.mod 中 go 1.26.2 改为 go 1.24.x 以彻底消除警告
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 允许用比 go.mod 声明版本低的工具链构建（避免 GOTOOLCHAIN=local 报错）
ENV GOTOOLCHAIN=auto

RUN apk add --no-cache git

# 单独拷贝 go.mod/go.sum 并下载依赖，充分利用层缓存
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源码并编译
# GOGC=off 关闭 GC 以减少编译期内存碎片压力
# -p 1 限制并行包数，避免 arm64/低内存环境 OOM
COPY . .
RUN GOGC=off go build -p 1 -ldflags="-w -s" -o bin/server ./cmd/server/main.go

# ── 运行 ────────────────────────────────────────────────────
FROM alpine:3.20 AS runner

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
 && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo "Asia/Shanghai" > /etc/timezone \
 && apk del tzdata

# 非 root 用户运行
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/bin/server ./server

# 创建必要目录并授权
RUN mkdir -p logs uploads && chown -R appuser:appgroup /app

USER appuser
EXPOSE 8080

CMD ["/app/bin/server"]