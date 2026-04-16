// app/[locale]/auth/register/page.tsx
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuthStore } from '@/store/auth';
import { useRegisterStore } from '@/store/register';
import toast from 'react-hot-toast';
import { Mail, Lock, User, Eye, EyeOff, UserPlus } from 'lucide-react';
import { useTranslations } from 'next-intl';
import Image from "next/image";

export default function RegisterPage() {
  const router = useRouter();
  const t = useTranslations('auth');
  const { user, isAuthenticated, isHydrated } = useAuthStore();
  
  const {
    username,
    email,
    password,
    confirmPassword,
    isLoading,
    errors,
    serverError,
    setUsername,
    setEmail,
    setPassword,
    setConfirmPassword,
    register: registerAction,
  } = useRegisterStore();
  
  const [showPassword, setShowPassword] = useState(false);

  // 已登录用户重定向到首页
  useEffect(() => {
    if (isHydrated && isAuthenticated && user) {
      router.replace('/');
    }
  }, [isHydrated, isAuthenticated, user, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const result = await registerAction();
    
    if (result.success) {
      const currentUser = useAuthStore.getState().user;
      toast.success(`${t("registration_successful")}，${currentUser?.username || ''}！`);
      router.push('/');
    } else {
      toast.error(serverError || result.message || t("registration_failed"));
    }
  };

  // 等待 hydration 完成
  if (!isHydrated) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  // 已登录用户不显示注册表单
  if (isAuthenticated && user) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-md">
        <div className="card bg-base-100 shadow-xl border border-base-300">
          <div className="card-body p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <Image 
                  src="/assets/brand/logo.svg" 
                  width={64} 
                  height={64} 
                  alt="logo" 
                  className="w-full h-full object-contain"
                />
              </div>
              <h1 className="text-2xl font-bold">{t("create_account")}</h1>
              <p className="text-base-content/50 text-sm mt-1">{t("join_the_forum")}</p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              {/* Username */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">{t("username")}</span>
                </label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder={t("username_placeholder")}
                    className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.username ? 'input-error' : ''}`}
                    autoComplete="username"
                    disabled={isLoading}
                  />
                </div>
                {errors.username && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.username}</span>
                  </label>
                )}
              </div>

              {/* Email */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">{t("email")}</span>
                </label>
                <div className="relative">
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="your@email.com"
                    className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.email ? 'input-error' : ''}`}
                    autoComplete="email"
                    disabled={isLoading}
                  />
                </div>
                {errors.email && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.email}</span>
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
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t("password_placeholder")}
                    className={`input input-bordered w-full pl-10 pr-10 focus:outline-none focus:border-primary ${errors.password ? 'input-error' : ''}`}
                    autoComplete="new-password"
                    disabled={isLoading}
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-base-content/40 hover:text-base-content"
                    onClick={() => setShowPassword(!showPassword)}
                    tabIndex={-1}
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                {errors.password && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.password}</span>
                  </label>
                )}
              </div>

              {/* Confirm Password */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">{t("confirm_password")}</span>
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder={t("confirm_password_placeholder")}
                    className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.confirmPassword ? 'input-error' : ''}`}
                    autoComplete="new-password"
                    disabled={isLoading}
                  />
                </div>
                {errors.confirmPassword && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.confirmPassword}</span>
                  </label>
                )}
              </div>

              {/* Server Error */}
              {serverError && !Object.keys(errors).length && (
                <div className="alert alert-error py-2">
                  <span className="text-sm">{serverError}</span>
                </div>
              )}

              {/* Submit Button */}
              <button
                type="submit"
                className="btn btn-primary w-full gap-2 mt-2"
                disabled={isLoading}
              >
                {isLoading ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  <UserPlus className="w-4 h-4" />
                )}
                {isLoading ? t("registering") : t("register")}
              </button>
            </form>

            <div className="divider text-base-content/30 text-xs">{t("already_have_an_account")}</div>
            <Link href="/auth/login" className="btn btn-ghost btn-sm w-full">
              {t("to_login")}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}