'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { authApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import toast from 'react-hot-toast';
import { Mail, Lock, Eye, EyeOff, LogIn } from 'lucide-react';
import { getErrorMessage } from '@/lib/utils';
import { useTranslations } from 'next-intl';

const loginSchema = z.object({
  email: z.string().email('请输入有效的邮箱'),
  password: z.string().min(1, '请输入密码'),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setAuth, user, isAuthenticated, isHydrated } = useAuthStore();
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const t = useTranslations('auth');

  // ✅ 获取重定向地址，默认为首页
  const redirectTo = searchParams.get('redirect') || '/';
  
  // ✅ 等待 hydration 完成后再检查登录状态
  useEffect(() => {
    if (isHydrated && isAuthenticated && user) {
      router.replace(redirectTo);
    }
  }, [isHydrated, isAuthenticated, user, router, redirectTo]);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({ resolver: zodResolver(loginSchema) });

  const onSubmit = async (data: LoginForm) => {
    setLoading(true);
    try {
      const res = await authApi.login(data);
      const { user } = res.data.data;
      setAuth(user);
      toast.success(t("welcome_back") + `，${user.username}！`);
      router.push(redirectTo);
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  // ✅ 等待 hydration 完成，避免闪烁
  if (!isHydrated) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  // ✅ 已登录用户不显示登录表单
  if (isAuthenticated && user) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  return (
    <div className="min-h-[80vh] flex items-center justify-center">
      <div className="w-full max-w-md">
        <div className="card bg-base-100 shadow-xl border border-base-300">
          <div className="card-body p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-2xl font-black mx-auto mb-4">
                B
              </div>
              <h1 className="text-2xl font-bold">{t("welcome_back")}</h1>
              <p className="text-base-content/50 text-sm mt-1">{t("login_to_your_account")}</p>
            </div>

            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              {/* Email */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">{t("email")}</span>
                </label>
                <div className="relative">
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    {...register('email')}
                    type="email"
                    placeholder="your@email.com"
                    className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.email ? 'input-error' : ''}`}
                    autoComplete="email"
                  />
                </div>
                {errors.email && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.email.message}</span>
                  </label>
                )}
              </div>

              {/* Password */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">{t("password")}</span>
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    {...register('password')}
                    type={showPassword ? 'text' : 'password'}
                    placeholder="••••••••"
                    className={`input input-bordered w-full pl-10 pr-10 focus:outline-none focus:border-primary ${errors.password ? 'input-error' : ''}`}
                    autoComplete="current-password"
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-base-content/40 hover:text-base-content"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                {errors.password && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.password.message}</span>
                  </label>
                )}
              </div>

              <button
                type="submit"
                className="btn btn-primary w-full gap-2 mt-2"
                disabled={loading}
              >
                {loading ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  <LogIn className="w-4 h-4" />
                )}
                {t("login")}
              </button>
            </form>

            <div className="divider text-base-content/30 text-xs">{t("Dont_have_an_account")}</div>
            <Link href="/auth/register" className="btn btn-ghost btn-sm w-full">
              {t("sign_up_for_free")}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}