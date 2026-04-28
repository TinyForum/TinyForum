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

// JWT Payload 类型定义
// interface JWTPayload {
//   id: number;
//   username: string;
//   email: string;
//   role: string;
//   exp: number;
//   [key: string]: unknown;
// }

export async function middleware(request: NextRequest): Promise<NextResponse> {
  // 获取路径
  let pathname: string = request.nextUrl.pathname;

  // 提取当前语言（从路径的第一段）
  const pathnameParts: string[] = pathname.split("/");
  let currentLocale: string = pathnameParts[1];

  // 验证语言是否有效
  if (!currentLocale || !routing.locales.includes(currentLocale as "en-US" || "zh-CN")) {
    currentLocale = routing.defaultLocale;
  }

  // 去除语言前缀的路径
  let pathnameWithoutLocale: string = pathname;
  for (const locale of routing.locales) {
    console.log("8. Checking locale:", locale);
    if (pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) {
      pathnameWithoutLocale = pathname.replace(`/${locale}`, "") || "/";
      break;
    }
  }

  const token: string | undefined = request.cookies.get("tiny_forum_token")?.value;
  const isAuthRoute: boolean = authRoutes.some((route: string) =>
    pathnameWithoutLocale.startsWith(route),
  );
  const isAdminRoute: boolean = adminRoutes.some((route: string) =>
    pathnameWithoutLocale.startsWith(route),
  );

  // 认证路由检查
  if (isAuthRoute && !token) {
    console.log("❌ No token, redirecting to login");
    const loginUrl: URL = new URL(`/${currentLocale}/auth/login`, request.url);
    loginUrl.searchParams.set("redirect", pathnameWithoutLocale);
    return NextResponse.redirect(loginUrl);
  }

  // 管理员路由验证
  if (isAdminRoute && token) {
    const jwt: string | undefined = process.env.JWT_SECRET;
    console.log("10. JWT_SECRET exists:", !!jwt);

    try {
      if (!jwt) {
        throw new Error("JWT_SECRET is not set");
      }
      const secret: Uint8Array = new TextEncoder().encode(jwt);
      const { payload }: JWTVerifyResult = await jwtVerify(token, secret);
      const role: string = payload.role as string;


      console.log("User role:", role);
      console.log("Role allowed:", allowedRoles.includes(role));

      if (!allowedRoles.includes(role)) {
        console.log("❌ Role not allowed, redirecting to home");
        return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
      }


    } catch (error: unknown) {
      console.error("❌ JWT verification failed:", error);
      const response: NextResponse = NextResponse.redirect(
        new URL(`/${currentLocale}/auth/login`, request.url),
      );
      response.cookies.delete("tiny_forum_token");
      return response;
    }
  }

  // 处理 i18n
  const intlResponse: NextResponse | undefined = intlMiddleware(request);

  if (intlResponse) {
    console.log("17. i18n middleware returned response");
    return intlResponse;
  }
  
  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|robots.txt|.*\\.(?:jpg|jpeg|gif|png|svg|ico|webp|css|js)$).*)",
  ],
};