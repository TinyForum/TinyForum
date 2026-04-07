'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { authApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import toast from 'react-hot-toast';
import { Mail, Lock, User, Eye, EyeOff, UserPlus } from 'lucide-react';
import { getErrorMessage } from '@/lib/utils';

const registerSchema = z.object({
  username: z.string().min(2, '用户名至少2个字符').max(50, '用户名最多50个字符'),
  email: z.string().email('请输入有效的邮箱'),
  password: z.string().min(6, '密码至少6个字符'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: '两次密码输入不一致',
  path: ['confirmPassword'],
});

type RegisterForm = z.infer<typeof registerSchema>;

export default function RegisterPage() {
  const router = useRouter();
  const { setAuth } = useAuthStore();
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({ resolver: zodResolver(registerSchema) });

  const onSubmit = async (data: RegisterForm) => {
    setLoading(true);
    try {
      const res = await authApi.register({
        username: data.username,
        email: data.email,
        password: data.password,
      });
      const { token, user } = res.data.data;
      setAuth(user, token);
      toast.success('注册成功，欢迎加入！');
      router.push('/');
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

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
              <h1 className="text-2xl font-bold">创建账号</h1>
              <p className="text-base-content/50 text-sm mt-1">加入 BBS Forum 社区</p>
            </div>

            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              {/* Username */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">用户名</span>
                </label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    {...register('username')}
                    type="text"
                    placeholder="你的用户名"
                    className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.username ? 'input-error' : ''}`}
                    autoComplete="username"
                  />
                </div>
                {errors.username && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.username.message}</span>
                  </label>
                )}
              </div>

              {/* Email */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">邮箱</span>
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
                  <span className="label-text font-medium">密码</span>
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    {...register('password')}
                    type={showPassword ? 'text' : 'password'}
                    placeholder="至少6个字符"
                    className={`input input-bordered w-full pl-10 pr-10 focus:outline-none focus:border-primary ${errors.password ? 'input-error' : ''}`}
                    autoComplete="new-password"
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

              {/* Confirm Password */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">确认密码</span>
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    {...register('confirmPassword')}
                    type={showPassword ? 'text' : 'password'}
                    placeholder="再次输入密码"
                    className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.confirmPassword ? 'input-error' : ''}`}
                    autoComplete="new-password"
                  />
                </div>
                {errors.confirmPassword && (
                  <label className="label pt-1">
                    <span className="label-text-alt text-error">{errors.confirmPassword.message}</span>
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
                  <UserPlus className="w-4 h-4" />
                )}
                注册
              </button>
            </form>

            <div className="divider text-base-content/30 text-xs">已有账号？</div>
            <Link href="/auth/login" className="btn btn-ghost btn-sm w-full">
              立即登录
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
