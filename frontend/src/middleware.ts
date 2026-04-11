import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import jwt from 'jsonwebtoken'
import { jwtVerify } from 'jose'

// 需要认证的路由
const authRoutes = ['/admin', '/settings', '/profile', '/posts/create']
// 管理员专用路由
const adminRoutes = ['/admin']

console.log('🔥 MIDDLEWARE FILE LOADED')
export async function middleware(request: NextRequest) {
  console.log("proxy is running")
  const token = request.cookies.get('bbs_token')?.value
  const { pathname } = request.nextUrl

  // 检查是否需要认证
  const isAuthRoute = authRoutes.some(route => pathname.startsWith(route))
  const isAdminRoute = adminRoutes.some(route => pathname.startsWith(route))

  if (isAuthRoute && !token) {
    // ✅ 修改为正确的登录路径
    const loginUrl = new URL('/auth/login', request.url)
    loginUrl.searchParams.set('redirect', pathname)
    return NextResponse.redirect(loginUrl)
  }

  if (isAdminRoute && token) {
    try {
      const secret = new TextEncoder().encode(process.env.JWT_SECRET || 'your-secret-key')
      const { payload } = await jwtVerify(token, secret)
      
      // 从 payload 中获取用户信息
      const role = payload.role as string
      const userId = payload.user_id || payload.userId || payload.sub
      
      if (role !== 'admin') {
        return NextResponse.redirect(new URL('/', request.url))
      }
    } catch (error) {
      console.error('JWT verification failed:', error)
      const response = NextResponse.redirect(new URL('/auth/login', request.url))
      response.cookies.delete('bbs_token')
      return response
    }
  }
  
  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|robots.txt).*)',
  ],
}