// src/middleware.ts
import { jwtVerify } from 'jose';
import createMiddleware from 'next-intl/middleware';
import { NextRequest, NextResponse } from 'next/server';
import { routing } from './i18n/routing';

// 创建 i18n middleware
const intlMiddleware = createMiddleware(routing);

// 需要认证的路由
const authRoutes = ['/admin', '/settings', '/posts/new'];
// 管理
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
  
  // 去除语言前缀的路径
  let pathnameWithoutLocale = pathname;
  for (const locale of routing.locales) {
    if (pathname === `/${locale}` || pathname.startsWith(`/${locale}/`)) {
      pathnameWithoutLocale = pathname.replace(`/${locale}`, '') || '/';
      break;
    }
  }

  const token = request.cookies.get('tiny_forum_token')?.value;
  const isAuthRoute = authRoutes.some(route => pathnameWithoutLocale.startsWith(route));
  const isAdminRoute = adminRoutes.some(route => pathnameWithoutLocale.startsWith(route));

  // 认证路由检查
  if (isAuthRoute && !token) {
    console.log('9. ❌ No token, redirecting to login');
    const loginUrl = new URL(`/${currentLocale}/auth/login`, request.url);
    loginUrl.searchParams.set('redirect', pathnameWithoutLocale);
    return NextResponse.redirect(loginUrl);
  }

  // 管理员路由验证
  if (isAdminRoute && token) {
    const jwt = process.env.JWT_SECRET;
    console.log('10. JWT_SECRET exists:', !!jwt);
    
    try {
      if (!jwt) {
        throw new Error('JWT_SECRET is not set');
      }
      const secret = new TextEncoder().encode(jwt);
      const { payload } = await jwtVerify(token, secret);
      const role = payload.role as string;
      
      console.log('11. ✅ JWT verified successfully');
      console.log('12. User role:', role);
      console.log('13. Role allowed:', allowedRoles.includes(role));
      
      if (!allowedRoles.includes(role)) {
        console.log('14. ❌ Role not allowed, redirecting to home');
        return NextResponse.redirect(new URL(`/${currentLocale}`, request.url));
      }
      
      console.log('15. ✅ Admin access granted');
    } catch (error) {
      console.error('16. ❌ JWT verification failed:', error);
      const response = NextResponse.redirect(new URL(`/${currentLocale}/auth/login`, request.url));
      response.cookies.delete('tiny_forum_token');
      return response;
    }
  }

  // 处理 i18n
  const intlResponse = intlMiddleware(request);
  
  if (intlResponse) {
    console.log('17. i18n middleware returned response');
    return intlResponse;
  }
  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|robots.txt|.*\\.(?:jpg|jpeg|gif|png|svg|ico|webp|css|js)$).*)',
  ],
};