// src/middleware.ts
import { jwtVerify } from 'jose';
import createMiddleware from 'next-intl/middleware';
import { NextRequest, NextResponse } from 'next/server';
import { routing } from './i18n/routing';

// 创建 i18n middleware
const intlMiddleware = createMiddleware(routing);

// 需要认证的路由
const authRoutes = ['/admin' ,'/settings'];
const adminRoutes = ['/admin'];
const allowedRoles = ['admin', 'super_admin'];

export async function middleware(request: NextRequest) {
  // 获取路径
  let pathname = request.nextUrl.pathname;
  
  // 提取当前语言（从路径的第一段）
  const pathnameParts = pathname.split('/');
  let currentLocale = pathnameParts[1];
  
  // 验证语言是否有效
  if (!currentLocale || !routing.locales.includes(currentLocale as any)) {
    currentLocale = routing.defaultLocale;
  }
  
  // 先处理 i18n，但不要立即返回，因为我们需要先做认证检查
  // 注意：createMiddleware 返回的是一个函数，我们需要调用它
  // 但由于 i18n middleware 可能会重定向，我们需要小心处理
  
  // 检查是否需要认证（使用原始路径，包含 locale）
  let pathnameWithoutLocale = pathname;
  for (const locale of routing.locales) {
    if (pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) {
      pathnameWithoutLocale = pathname.replace(`/${locale}`, '') || '/';
      break;
    }
  }

  const token = request.cookies.get('tiny_forum_token')?.value;
// console.log('[middleware] pathname:', pathnameWithoutLocale, '| token:', token ?? 'NOT FOUND');
  const isAuthRoute = authRoutes.some(route => pathnameWithoutLocale.startsWith(route));
  const isAdminRoute = adminRoutes.some(route => pathnameWithoutLocale.startsWith(route));

  // 认证路由检查（在 i18n 处理之前）
  if (isAuthRoute && !token) {
    const loginUrl = new URL(`/${currentLocale}/auth/login`, request.url);
    loginUrl.searchParams.set('redirect', pathnameWithoutLocale);
    return NextResponse.redirect(loginUrl);
  }

  // 管理员路由验证
  if (isAdminRoute && token) {
    const jwt = process.env.NEXT_PUBLIC_JWT_SECRET
    try {
      if (!jwt) {
        throw new Error('JWT_SECRET is not set\n Check your environment variables');
      }
      const secret = new TextEncoder().encode(jwt);
      const { payload } = await jwtVerify(token, secret);
      const role = payload.role as string;
      
      if (!allowedRoles.includes(role)) {
  return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
}
    } catch (error) {
      console.error('JWT verification failed:', error);
      const response = NextResponse.redirect(new URL(`/${currentLocale}/auth/login`, request.url));
      response.cookies.delete('tiny_forum_token');
      return response;
    }
  }

  // 最后处理 i18n（只在没有重定向的情况下）
  // 注意：intlMiddleware 本身可能返回重定向（如语言重定向）
  const intlResponse = intlMiddleware(request);
  
  // 如果 i18n middleware 返回了响应，直接返回
  if (intlResponse) {
    return intlResponse;
  }

  // 如果没有特殊处理，继续
  return NextResponse.next();
}

export const config = {
  matcher: [
    // 匹配所有路径，除了 API 和静态文件
    '/((?!api|_next/static|_next/image|favicon.ico|robots.txt|.*\\.(?:jpg|jpeg|gif|png|svg|ico|webp|css|js)$).*)',
  ],
};