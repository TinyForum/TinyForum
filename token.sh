#!/bin/bash

LOGIN_URL="http://localhost:8080/api/v1/auth/login"
TEST_URL="http://localhost:8080/api/v1/admin/posts?page=1&page_size=20"
COOKIE_JAR="cookies.txt"

# 1. 登录并保存 Cookie
echo "正在登录..."
curl -s -X POST "$LOGIN_URL" \
    -H "Content-Type: application/json" \
    -d '{"email": "admin@test.com", "password": "adminadmin"}' \
    -c "$COOKIE_JAR" -o /dev/null

# 检查 Cookie 是否保存成功
if [ ! -s "$COOKIE_JAR" ]; then
    echo "登录失败：未收到 Cookie"
    exit 1
fi

echo "登录成功，Cookie 已保存"

# 2. 使用保存的 Cookie 访问管理接口
echo -e "\n测试管理接口（使用 Cookie 认证）："
curl -X GET "$TEST_URL" \
    -b "$COOKIE_JAR" \
    -v

# 可选：查看保存的 Cookie 内容
echo -e "\nCookie 内容："
cat "$COOKIE_JAR"