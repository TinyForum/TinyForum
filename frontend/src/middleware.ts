// src/middleware.ts
import { jwtVerify, JWTVerifyResult } from "jose";
import createMiddleware from "next-intl/middleware";
import { NextRequest, NextResponse } from "next/server";
import { routing } from "./i18n/routing";

// 创建 i18n middleware
const intlMiddleware = createMiddleware(routing);

// 需要认证的路由
const authRoutes: string[] = ["/dashboard/admin", "/settings", "/posts/new"];
// 管理路由
const adminRoutes: string[] = ["/dashboard/admin"];
const allowedRoles: string[] = ["admin", "super_admin"];

// 调试函数：打印请求详情
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

// 调试函数：打印所有 Cookie
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

export async function middleware(request: NextRequest): Promise<NextResponse> {
  console.log("\n" + "=".repeat(60));
  console.log("🚀 MIDDLEWARE START");
  console.log("=".repeat(60));

  // 打印请求详情
  logRequestDetails(request, "START");

  // 获取路径
  let pathname: string = request.nextUrl.pathname;
  console.log(`\n📂 Original Pathname: ${pathname}`);

  // 提取当前语言（从路径的第一段）
  const pathnameParts: string[] = pathname.split("/");
  let currentLocale: string = pathnameParts[1];

  // 验证语言是否有效
  if (
    !currentLocale ||
    !routing.locales.includes((currentLocale as "en-US") || "zh-CN")
  ) {
    currentLocale = routing.defaultLocale;
    console.log(`🌍 Using default locale: ${currentLocale}`);
  } else {
    console.log(`🌍 Detected locale: ${currentLocale}`);
  }

  // 去除语言前缀的路径
  let pathnameWithoutLocale: string = pathname;
  for (const locale of routing.locales) {
    console.log(`  Checking locale: ${locale}`);
    if (pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) {
      pathnameWithoutLocale = pathname.replace(`/${locale}`, "") || "/";
      console.log(
        `  ✅ Matched locale ${locale}, path without locale: ${pathnameWithoutLocale}`,
      );
      break;
    }
  }

  // 获取 Token（从 Cookie）
  const token: string | undefined =
    request.cookies.get("tiny_forum_token")?.value;
  console.log(
    `\n🔑 Token from cookie: ${token ? `${token.substring(0, 20)}... (length: ${token.length})` : "NOT FOUND"}`,
  );

  // 也检查 Authorization header（备份方案）
  const authHeader = request.headers.get("authorization");
  const headerToken = authHeader?.startsWith("Bearer ")
    ? authHeader.substring(7)
    : null;
  if (headerToken && !token) {
    console.log(
      `📝 Token found in Authorization header: ${headerToken.substring(0, 20)}...`,
    );
  }

  // 打印所有 Cookie（调试用）
  logAllCookies(request);

  // 检查是否是 API 请求
  const isApiRequest = pathname.includes("/api/");
  if (isApiRequest) {
    console.log(`🔌 API Request detected: ${pathname}`);
    // 对于 API 请求，检查是否需要认证
    const needsAuth =
      pathname.includes("/auth/logout") ||
      pathname.includes("/users/me") ||
      pathname.includes("/timeline/following") ||
      pathname.includes("/notifications");

    if (needsAuth && !token && !headerToken) {
      console.log(
        `❌ API ${pathname} requires authentication but no token provided`,
      );
      // 返回 401 响应
      const response = NextResponse.json(
        { code: 40101, message: "未认证，请先登录" },
        { status: 401 },
      );
      console.log("=".repeat(60) + "\n");
      return response;
    }
  }

  // 检查是否是认证路由
  const isAuthRoute: boolean = authRoutes.some((route: string) =>
    pathnameWithoutLocale.startsWith(route),
  );
  const isAdminRoute: boolean = adminRoutes.some((route: string) =>
    pathnameWithoutLocale.startsWith(route),
  );

  console.log(`\n🛣️ Route Analysis:`);
  console.log(`  Path without locale: ${pathnameWithoutLocale}`);
  console.log(`  Is auth route: ${isAuthRoute}`);
  console.log(`  Is admin route: ${isAdminRoute}`);
  console.log(`  Has token: ${!!token || !!headerToken}`);

  // 认证路由检查
  if (isAuthRoute && !token && !headerToken) {
    console.log(
      `❌ Auth route ${pathnameWithoutLocale} requires authentication but no token found`,
    );
    const loginUrl: URL = new URL(`/${currentLocale}/auth/login`, request.url);
    loginUrl.searchParams.set("redirect", pathnameWithoutLocale);
    console.log(`🔄 Redirecting to: ${loginUrl.toString()}`);
    console.log("=".repeat(60) + "\n");
    return NextResponse.redirect(loginUrl);
  }

  // 管理员路由验证
  if (isAdminRoute && (token || headerToken)) {
    const finalToken = token || headerToken;
    const jwt: string | undefined = process.env.JWT_SECRET;
    console.log(`\n🔐 Admin route verification:`);
    console.log(`  JWT_SECRET exists: ${!!jwt}`);
    console.log(`  Token length: ${finalToken?.length}`);
    console.log(`🔐 JWT_SECRET from env: "${process.env.JWT_SECRET}"`);
    console.log(`🔐 JWT_SECRET length: ${process.env.JWT_SECRET?.length}`);

    try {
      if (!jwt) {
        throw new Error("JWT_SECRET is not set in environment variables");
      }

      const secret: Uint8Array = new TextEncoder().encode(jwt);
      console.log(`  Verifying token...`);
      const { payload }: JWTVerifyResult = await jwtVerify(finalToken!, secret);
      const role: string = payload.role as string;
      const userId: number = payload.id as number;
      const username: string = payload.username as string;

      console.log(`  ✅ Token verified successfully:`);
      console.log(`    - User ID: ${userId}`);
      console.log(`    - Username: ${username}`);
      console.log(`    - Role: ${role}`);
      console.log(
        `    - Expiration: ${new Date((payload.exp || 0) * 1000).toLocaleString()}`,
      );
      console.log(`  Role allowed: ${allowedRoles.includes(role)}`);

      if (!allowedRoles.includes(role)) {
        console.log(
          `❌ Role ${role} not allowed for admin route, redirecting to home`,
        );
        console.log("=".repeat(60) + "\n");
        return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
      }

      console.log(`✅ Admin access granted for ${username}`);
    } catch (error: unknown) {
      console.error(`❌ JWT verification failed:`, error);
      if (error instanceof Error) {
        console.error(`  Error message: ${error.message}`);
        console.error(`  Error stack: ${error.stack}`);
      }

      const response: NextResponse = NextResponse.redirect(
        new URL(`/${currentLocale}/auth/login`, request.url),
      );
      response.cookies.delete("tiny_forum_token");
      console.log(`  Deleted invalid token cookie`);
      console.log(`🔄 Redirecting to login`);
      console.log("=".repeat(60) + "\n");
      return response;
    }
  }

  // 对于 API 请求，如果有 token，添加到 headers
  if (isApiRequest && (token || headerToken)) {
    console.log(`📤 Adding token to API request headers`);
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set("Authorization", `Bearer ${token || headerToken}`);

    // 修改请求以包含 Authorization header
    const response = intlMiddleware(request);
    if (response) {
      console.log("  Token added to request headers");
      console.log("=".repeat(60) + "\n");
      return response;
    }
  }

  // 处理 i18n
  console.log(`\n🌐 Processing i18n middleware...`);
  const intlResponse: NextResponse | undefined = intlMiddleware(request);

  if (intlResponse) {
    console.log(`✅ i18n middleware returned response`);
    console.log("=".repeat(60) + "\n");
    return intlResponse;
  }

  console.log(`✅ No middleware action needed, continuing`);
  console.log("=".repeat(60) + "\n");
  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|robots.txt|.*\\.(?:jpg|jpeg|gif|png|svg|ico|webp|css|js)$).*)",
  ],
};
