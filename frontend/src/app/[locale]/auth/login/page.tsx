'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useAuthStore } from '@/store/auth';
import { useLoginStore } from '@/store/login';
import toast from 'react-hot-toast';
import { Mail, Lock, Eye, EyeOff, LogIn } from 'lucide-react';
import { useTranslations, useLocale } from 'next-intl';
import { authApi } from '@/lib/api';

const loginSchema = z.object({
  email: z.string().email('请输入有效的邮箱'),
  password: z.string().min(1, '请输入密码'),
});

type LoginForm = z.infer<typeof loginSchema>;

interface DeletionStatus {
  is_deleted: boolean;
  deleted_at?: string;
  can_restore: boolean;
  remaining_days?: number;
}

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const locale = useLocale();
  const t = useTranslations('Auth');
  
  const { user, isAuthenticated, isHydrated, logout: storeLogout } = useAuthStore();
  const { 
    email, 
    password, 
    rememberMe,
    isLoading, 
    error,
    setEmail, 
    setPassword, 
    setRememberMe,
    setError,
    login 
  } = useLoginStore();
  
  const [showPassword, setShowPassword] = useState(false);
  const [showRestoreDialog, setShowRestoreDialog] = useState(false);
  
  const redirectTo = searchParams.get('redirect') || '/';

  useEffect(() => {
    if (isHydrated && isAuthenticated && user && !showRestoreDialog) {
      const destination = redirectTo === '/' ? `/${locale}` : redirectTo;
      router.replace(destination);
    }
  }, [isHydrated, isAuthenticated, user, router, redirectTo, locale, showRestoreDialog]);

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
  } = useForm<LoginForm>({ 
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: email || '',
      password: password || '',
    }
  });

  useEffect(() => {
    setValue('email', email);
    setValue('password', password);
  }, [email, password, setValue]);

  // 检查账户删除状态
  const checkDeletionStatus = async () => {
    try {
      const response = await authApi.getDeletionStatus();
      const status = response.data.data;
      
      if (status.is_deleted && status.can_restore) {
        setShowRestoreDialog(true);
        return true;
      } else if (status.is_deleted && !status.can_restore) {
        toast.error('您的账户已被永久删除，请联系管理员');
        await handleForceLogout();
        return false;
      }
      return false;
    } catch (error) {
      console.error('获取删除状态失败:', error);
      return false;
    }
  };

  // 强制退出（用于永久删除的账户）
  const handleForceLogout = async () => {
    await authApi.logout();
    storeLogout();
    router.push('/');
  };

  const onSubmit = async (data: LoginForm) => {
    setEmail(data.email);
    setPassword(data.password);
    
    const result = await login();
    
    if (result.success) {
      // 检查账户删除状态
      await checkDeletionStatus();
    } else {
      toast.error(error || t("login_failed"));
    }
  };

  if (!isHydrated) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  return (
    <>
      <div className="min-h-[80vh] flex items-center justify-center px-4 py-8">
        <div className="w-full max-w-md">
          <div className="card bg-base-100 shadow-xl border border-base-300">
            <div className="card-body p-8">
              <div className="text-center mb-8">
                <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-2xl font-black mx-auto mb-4">
                  B
                </div>
                <h1 className="text-2xl font-bold">{t("welcome_back")}</h1>
                <p className="text-base-content/50 text-sm mt-1">{t("login_to_your_account")}</p>
              </div>

              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
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
                      className={`input input-bordered w-full pl-10 ${errors.email ? 'input-error' : ''}`}
                      autoComplete="email"
                      onChange={(e) => {
                        register('email').onChange(e);
                        setEmail(e.target.value);
                        setError(null);
                      }}
                    />
                  </div>
                  {errors.email && (
                    <label className="label pt-1">
                      <span className="label-text-alt text-error">{errors.email.message}</span>
                    </label>
                  )}
                </div>

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
                      className={`input input-bordered w-full pl-10 pr-10 ${errors.password ? 'input-error' : ''}`}
                      autoComplete="current-password"
                      onChange={(e) => {
                        register('password').onChange(e);
                        setPassword(e.target.value);
                        setError(null);
                      }}
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

                <div className="flex items-center justify-between">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={rememberMe}
                      onChange={(e) => setRememberMe(e.target.checked)}
                      className="checkbox checkbox-sm checkbox-primary"
                    />
                    <span className="text-sm">{t("remember_me")}</span>
                  </label>
                  <Link 
                    href="/auth/forgot-password" 
                    className="text-sm text-primary hover:underline"
                  >
                    {t("forgot_password")}
                  </Link>
                </div>

                {error && (
                  <div className="alert alert-error py-2">
                    <span className="text-sm">{error}</span>
                  </div>
                )}

                <button
                  type="submit"
                  className="btn btn-primary w-full gap-2 mt-2"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <span className="loading loading-spinner loading-sm" />
                  ) : (
                    <LogIn className="w-4 h-4" />
                  )}
                  {isLoading ? t("logging_in") : t("login")}
                </button>
              </form>

              <div className="divider text-base-content/30 text-xs">{t("dont_have_account")}</div>
              <Link href="/auth/register" className="btn btn-ghost btn-sm w-full">
                {t("sign_up_for_free")}
              </Link>
            </div>
          </div>
        </div>
      </div>

      
    </>
  );
}