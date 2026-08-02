// src/middleware.ts
import { jwtVerify } from "jose";
import createMiddleware from "next-intl/middleware";
import { NextRequest, NextResponse } from "next/server";
import { routing } from "./i18n/routing";

// 创建 i18n 中间件实例
const intlMiddleware = createMiddleware(routing);

// ---------- 路由配置 ----------
const publicAuthRoutes = ["/auth/login", "/auth/register"]; // 公开但需特殊处理（已登录时重定向）
const authRoutes = ["/dashboard/admin", "/settings", "/posts/new"]; // 需要登录的页面
const adminRoutes = ["/dashboard/admin"]; // 需要管理员角色的页面
const allowedRoles = ["admin", "super_admin"];

// ---------- 调试日志（开发环境可用，生产可移除） ----------
function logRequestDetails(request: NextRequest, stage: string) {
  console.log(`\n🔍 [${stage}] Request Details:`);
  console.log(`  URL: ${request.url}`);
  console.log(`  Method: ${request.method}`);
  console.log(`  Headers:`);
  console.log(`    - Cookie: ${request.headers.get("cookie") || "none"}`);
  console.log(
    `    - Authorization: ${request.headers.get("authorization") || "none"}`,
  );
  console.log(
    `    - User-Agent: ${request.headers.get("user-agent") || "unknown"}`,
  );
}

function logAllCookies(request: NextRequest) {
  const cookieHeader = request.headers.get("cookie");
  console.log(`\n🍪 All Cookies: ${cookieHeader || "none"}`);
  if (cookieHeader) {
    const cookies = cookieHeader.split(";").map((c) => c.trim());
    cookies.forEach((cookie) => {
      const [name, value] = cookie.split("=");
      console.log(
        `  - ${name}: ${value ? `${value.substring(0, 20)}...` : "empty"}`,
      );
    });
  }
}

// ---------- 核心中间件 ----------
export async function middleware(request: NextRequest): Promise<NextResponse> {
  console.log("\n" + "=".repeat(60));
  console.log("🚀 MIDDLEWARE START");
  console.log("=".repeat(60));
  logRequestDetails(request, "START");

  // ----- 1. 解析路径和语言 -----
  const pathname = request.nextUrl.pathname;
  console.log(`\n📂 Original Pathname: ${pathname}`);

  // 提取当前语言（路径第一段）
  const pathnameParts = pathname.split("/");
  let currentLocale = pathnameParts[1] || routing.defaultLocale;

  // 验证语言是否有效
  if (!routing.locales.includes(currentLocale as any)) {
    currentLocale = routing.defaultLocale;
    console.log(`🌍 Using default locale: ${currentLocale}`);
  } else {
    console.log(`🌍 Detected locale: ${currentLocale}`);
  }

  // 去除语言前缀，得到不带 locale 的路径
  let pathnameWithoutLocale = pathname;
  for (const locale of routing.locales) {
    if (pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) {
      pathnameWithoutLocale = pathname.replace(`/${locale}`, "") || "/";
      console.log(
        `  ✅ Matched locale ${locale}, path without locale: ${pathnameWithoutLocale}`,
      );
      break;
    }
  }

  // ----- 2. 获取 Token（Cookie 优先，Authorization Header 作为备选） -----
  const token = request.cookies.get("tiny_forum_token")?.value;
  const authHeader = request.headers.get("authorization");
  const headerToken = authHeader?.startsWith("Bearer ")
    ? authHeader.substring(7)
    : null;
  const finalToken = token || headerToken;

  console.log(
    `\n🔑 Token from cookie: ${token ? `${token.substring(0, 20)}... (length: ${token.length})` : "NOT FOUND"}`,
  );
  if (headerToken && !token) {
    console.log(
      `📝 Token found in Authorization header: ${headerToken.substring(0, 20)}...`,
    );
  }
  logAllCookies(request);

  console.log(`\n🛣️ Route Analysis:`);
  console.log(`  Path without locale: ${pathnameWithoutLocale}`);
  console.log(`  Has token: ${!!finalToken}`);

  // ----- 3. 公开认证路由（登录/注册）处理 -----
  // 这类页面：已登录用户重定向到首页，未登录或无效 token 则正常显示
  if (publicAuthRoutes.includes(pathnameWithoutLocale)) {
    // 如果有 token，尝试验证
    if (finalToken) {
      try {
        const secret = new TextEncoder().encode(process.env.JWT_SECRET);
        await jwtVerify(finalToken, secret);
        // token 有效 → 重定向到首页
        console.log(`✅ 已登录用户访问公开认证页，重定向到首页`);
        return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
      } catch {
        // token 无效 → 清除 cookie，并返回 i18n 响应（显示登录页）
        console.warn(`⚠️ 无效 token，已清除，允许访问登录页`);
        const response = intlMiddleware(request);
        if (response) {
          response.cookies.delete("tiny_forum_token");
          return response;
        }
        // 保底：生成新响应并清除 cookie
        const fallback = NextResponse.next();
        fallback.cookies.delete("tiny_forum_token");
        return fallback;
      }
    }
    // 无 token → 直接渲染登录页（不再执行任何其他检查）
    console.log(`🔓 未登录，允许访问登录页`);
    const response = intlMiddleware(request);
    return response || NextResponse.next();
  }

  // ----- 4. API 请求处理（仅做简单的身份验证，不修改请求头） -----
  if (pathname.includes("/api/")) {
    const needsAuth = [
      "/auth/logout",
      "/users/me",
      "/timeline/following",
      "/notifications",
    ].some((api) => pathname.includes(api));
    if (needsAuth && !finalToken) {
      console.log(`❌ API ${pathname} 需要认证但无 token`);
      return NextResponse.json(
        { code: 40101, message: "未认证，请先登录" },
        { status: 401 },
      );
    }
    // 其他 API 请求放行（客户端自行携带 token）
    console.log(`🔌 API 请求放行`);
    return intlMiddleware(request) || NextResponse.next();
  }

  // ----- 5. 普通受保护路由（需要登录） -----
  const isAuthRoute = authRoutes.some((route) =>
    pathnameWithoutLocale.startsWith(route),
  );
  if (isAuthRoute && !finalToken) {
    console.log(
      `❌ 受保护路由 ${pathnameWithoutLocale} 需要登录，重定向到登录页`,
    );
    const loginUrl = new URL(`/${currentLocale}/auth/login`, request.url);
    loginUrl.searchParams.set("redirect", pathnameWithoutLocale);
    return NextResponse.redirect(loginUrl);
  }

  // ----- 6. 管理员路由（需要特定角色） -----
  const isAdminRoute = adminRoutes.some((route) =>
    pathnameWithoutLocale.startsWith(route),
  );
  if (isAdminRoute && finalToken) {
    try {
      const secret = new TextEncoder().encode(process.env.JWT_SECRET);
      const { payload } = await jwtVerify(finalToken, secret);
      const role = payload.role as string;
      if (!allowedRoles.includes(role)) {
        console.log(`❌ 角色 ${role} 无权访问管理员页面，重定向到首页`);
        return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
      }
      console.log(`✅ 管理员访问允许（角色：${role}）`);
    } catch (error) {
      // token 无效 → 清除并重定向到登录页
      console.warn(`⚠️ 管理员路由验证失败，清除 token 并跳转登录`);
      const response = NextResponse.redirect(
        new URL(`/${currentLocale}/auth/login`, request.url),
      );
      response.cookies.delete("tiny_forum_token");
      return response;
    }
  }

  // ----- 7. 其他所有请求（默认走 i18n 中间件） -----
  console.log(`🌐 执行 i18n 中间件（默认）`);
  const intlResponse = intlMiddleware(request);
  if (intlResponse) {
    console.log(`✅ i18n 中间件返回响应`);
    console.log("=".repeat(60) + "\n");
    return intlResponse;
  }

  console.log(`✅ 无特殊处理，继续执行`);
  console.log("=".repeat(60) + "\n");
  return NextResponse.next();
}

// ----- 中间件匹配器（排除静态资源等） -----
export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|robots.txt|.*\\.(?:jpg|jpeg|gif|png|svg|ico|webp|css|js)$).*)",
  ],
};
