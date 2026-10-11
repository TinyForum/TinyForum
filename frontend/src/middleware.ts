// src/middleware.ts
import { jwtVerify, JWTVerifyResult } from "jose";
import createMiddleware from "next-intl/middleware";
import { NextRequest, NextResponse } from "next/server";
import { routing } from "./i18n/routing";

// 创建 i18n middleware
const intlMiddleware = createMiddleware(routing);

// 需要认证的路由
const authRoutes: string[] = ["/dashboard/admin", "/settings", "/posts/new"];
// 管理路由（需要校验 JWT 与角色）
const adminRoutes: string[] = ["/dashboard/admin"];
const allowedRoles: string[] = ["admin", "super_admin"];

export async function middleware(request: NextRequest): Promise<NextResponse> {
  const pathname: string = request.nextUrl.pathname;

  // 提取当前语言（从路径的第一段）
  const pathnameParts: string[] = pathname.split("/");
  const rawLocale: string = pathnameParts[1];
  const currentLocale =
    rawLocale && routing.locales.includes(rawLocale as (typeof routing.locales)[number])
      ? rawLocale
      : routing.defaultLocale;

  // 去除语言前缀的路径
  let pathnameWithoutLocale: string = pathname;
  for (const locale of routing.locales) {
    if (pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) {
      pathnameWithoutLocale = pathname.replace(`/${locale}`, "") || "/";
      break;
    }
  }

  // 获取 Token（Cookie 优先，Authorization Header 兜底）
  const token: string | undefined =
    request.cookies.get("tiny_forum_token")?.value ||
    request.headers.get("authorization")?.replace(/^Bearer\s+/, "") ||
    undefined;

  // 认证路由：无 token 则跳转登录
  const isAuthRoute: boolean = authRoutes.some((route: string) =>
    pathnameWithoutLocale.startsWith(route),
  );
  if (isAuthRoute && !token) {
    const loginUrl: URL = new URL(`/${currentLocale}/auth/login`, request.url);
    loginUrl.searchParams.set("redirect", pathnameWithoutLocale);
    return NextResponse.redirect(loginUrl);
  }

  // 管理路由：校验 JWT 签名与角色
  const isAdminRoute: boolean = adminRoutes.some((route: string) =>
    pathnameWithoutLocale.startsWith(route),
  );
  if (isAdminRoute && token) {
    const jwt: string | undefined = process.env.JWT_SECRET;
    if (!jwt) {
      return NextResponse.redirect(
        new URL(`/${currentLocale}/auth/login`, request.url),
      );
    }

    try {
      const secret: Uint8Array = new TextEncoder().encode(jwt);
      const { payload }: JWTVerifyResult = await jwtVerify(token, secret);
      const role: string = payload.role as string;
      if (!allowedRoles.includes(role)) {
        return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
      }
    } catch {
      const response: NextResponse = NextResponse.redirect(
        new URL(`/${currentLocale}/auth/login`, request.url),
      );
      response.cookies.delete("tiny_forum_token");
      return response;
    }
  }

  const response: NextResponse | undefined = intlMiddleware(request);
  return response ?? NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|robots.txt|uploads|store|.*\\.(?:jpg|jpeg|gif|png|svg|ico|webp|mp4|webm|ogg|mov|avi|mkv|css|js)$).*)",
  ],
};