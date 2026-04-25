# ── 依赖安装 ────────────────────────────────────────────────
FROM node:24-alpine AS deps

WORKDIR /app
COPY package.json ./
# package-lock.json 与 package.json 版本不同步（锁文件是 next@14 时代产物）
# 用 npm install 重新解析依赖树，跳过 peer dep 冲突
# 构建成功后同步锁文件回宿主机：
#   podman cp $(podman ps -lq):/app/package-lock.json ./frontend/
RUN npm install --legacy-peer-deps

# ── 构建 ────────────────────────────────────────────────────
FROM node:20-alpine AS builder

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/package-lock.json ./package-lock.json
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# ── 运行 ────────────────────────────────────────────────────
FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs \
 && adduser  --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT=3000

CMD ["node", "server.js"]