'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { authApi, userApi } from '@/lib/api';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import { getErrorMessage } from '@/lib/utils';
import { 
  Save, 
  KeyRound, 
  Eye, 
  EyeOff, 
  Shield, 
  CheckCircle, 
  AlertCircle,
  Smartphone,
  Mail,
  History
} from 'lucide-react';

// 弱密码验证（仅用于提示，不强制）
const passwordSchema = z.object({
  oldPassword: z.string().min(1, '请输入旧密码'),
  newPassword: z.string().min(6, '密码长度至少6位'),
  confirmPassword: z.string().min(1, '请确认新密码'),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: '两次输入的密码不一致',
  path: ['confirmPassword'],
});

type PasswordForm = z.infer<typeof passwordSchema>;

// 密码强度计算
function calculatePasswordStrength(password: string) {
  if (!password) return { score: 0, level: '无', color: 'bg-gray-200', message: '' };
  
  let score = 0;
  if (password.length >= 6) score++;
  if (password.length >= 10) score++;
  if (/[a-z]/.test(password)) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^A-Za-z0-9]/.test(password)) score++;
  
  if (score <= 2) return { score, level: '弱', color: 'bg-red-500', message: '密码强度较弱，建议增加复杂度' };
  if (score <= 4) return { score, level: '中', color: 'bg-yellow-500', message: '密码强度中等，可以更强一些' };
  return { score, level: '强', color: 'bg-green-500', message: '密码强度很好！' };
}

interface SecuritySettingsProps {
  user: any;
}

export default function SecuritySettings({ user }: SecuritySettingsProps) {
  const router = useRouter();
  const [passwordLoading, setPasswordLoading] = useState(false);
  const [showOldPassword, setShowOldPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState({ score: 0, level: '无', color: 'bg-gray-200', message: '' });

  const { 
    register, 
    handleSubmit, 
    reset,
    watch,
    formState: { errors, isValid } 
  } = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
    mode: 'onChange',
    defaultValues: {
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    },
  });

  const newPasswordValue = watch('newPassword');

  useEffect(() => {
    setPasswordStrength(calculatePasswordStrength(newPasswordValue || ''));
  }, [newPasswordValue]);

const onSubmit = async (data: PasswordForm) => {
  setPasswordLoading(true);
  try {
    // 修改这里：将驼峰命名转换为下划线命名
    await userApi.changePassword({
      old_password: data.oldPassword, 
      new_password: data.newPassword
    });
  
    toast.success('密码修改成功');
    reset();
    
    setTimeout(() => {
      toast.success('即将退出，请使用新密码重新登录');
      authApi.logout();
      router.push('/auth/login');
    }, 3000);
  } catch (err) {
    toast.error(getErrorMessage(err));
  } finally {
    setPasswordLoading(false);
  }
};
  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold">安全设置</h1>
        <p className="text-sm text-base-content/60 mt-1">保护您的账户安全</p>
      </div>

      {/* 安全概览卡片 */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="card bg-base-100 border border-base-200 shadow-sm">
          <div className="card-body p-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-green-100 dark:bg-green-900/20 rounded-lg">
                <Shield className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p className="text-xs text-base-content/60">安全等级</p>
                <p className="font-semibold text-sm">中等</p>
              </div>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 border border-base-200 shadow-sm">
          <div className="card-body p-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-blue-100 dark:bg-blue-900/20 rounded-lg">
                <History className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p className="text-xs text-base-content/60">最后登录</p>
                <p className="font-semibold text-sm">今天 14:30</p>
              </div>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 border border-base-200 shadow-sm">
          <div className="card-body p-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-purple-100 dark:bg-purple-900/20 rounded-lg">
                <Mail className="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <p className="text-xs text-base-content/60">邮箱状态</p>
                <p className="font-semibold text-sm">已验证</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 修改密码卡片 */}
      <div className="card bg-base-100 border border-base-200 shadow-sm">
        <div className="card-body p-6">
          <div className="flex items-center gap-2 mb-4">
            <KeyRound className="w-5 h-5 text-primary" />
            <h2 className="text-xl font-semibold">修改密码</h2>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {/* 旧密码 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">当前密码</span>
              </label>
              <div className="relative">
                <input
                  {...register('oldPassword')}
                  type={showOldPassword ? "text" : "password"}
                  placeholder="请输入当前密码"
                  className={`input input-bordered w-full pr-10 ${errors.oldPassword ? 'input-error' : ''}`}
                />
                <button
                  type="button"
                  onClick={() => setShowOldPassword(!showOldPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2"
                >
                  {showOldPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
              {errors.oldPassword && (
                <p className="text-error text-xs mt-1 flex items-center gap-1">
                  <AlertCircle className="w-3 h-3" />
                  {errors.oldPassword.message}
                </p>
              )}
            </div>

            {/* 新密码 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">新密码</span>
              </label>
              <div className="relative">
                <input
                  {...register('newPassword')}
                  type={showNewPassword ? "text" : "password"}
                  placeholder="请输入新密码"
                  className={`input input-bordered w-full pr-10 ${errors.newPassword ? 'input-error' : ''}`}
                />
                <button
                  type="button"
                  onClick={() => setShowNewPassword(!showNewPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2"
                >
                  {showNewPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
              
              {/* 密码强度 */}
              {newPasswordValue && (
                <div className="mt-2 space-y-1">
                  <div className="flex justify-between text-xs">
                    <span>密码强度</span>
                    <span className={`font-medium ${
                      passwordStrength.level === '强' ? 'text-green-600' : 
                      passwordStrength.level === '中' ? 'text-yellow-600' : 'text-red-600'
                    }`}>
                      {passwordStrength.level}
                    </span>
                  </div>
                  <div className="w-full h-1.5 bg-gray-200 rounded-full overflow-hidden">
                    <div 
                      className={`h-full transition-all duration-300 ${passwordStrength.color}`}
                      style={{ width: `${(passwordStrength.score / 6) * 100}%` }}
                    />
                  </div>
                  <p className="text-xs text-base-content/60">{passwordStrength.message}</p>
                </div>
              )}
              
              {errors.newPassword && (
                <p className="text-error text-xs mt-1">{errors.newPassword.message}</p>
              )}
              <p className="text-xs text-base-content/40 mt-1">
                💡 建议：使用大小写字母、数字和特殊字符组合
              </p>
            </div>

            {/* 确认密码 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">确认新密码</span>
              </label>
              <div className="relative">
                <input
                  {...register('confirmPassword')}
                  type={showConfirmPassword ? "text" : "password"}
                  placeholder="请再次输入新密码"
                  className={`input input-bordered w-full pr-10 ${errors.confirmPassword ? 'input-error' : ''}`}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2"
                >
                  {showConfirmPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
              {errors.confirmPassword && (
                <p className="text-error text-xs mt-1">{errors.confirmPassword.message}</p>
              )}
              {watch('confirmPassword') && !errors.confirmPassword && watch('newPassword') === watch('confirmPassword') && (
                <p className="text-success text-xs mt-1 flex items-center gap-1">
                  <CheckCircle className="w-3 h-3" />
                  密码匹配
                </p>
              )}
            </div>

            <button
              type="submit"
              className="btn btn-primary w-full gap-2"
              disabled={passwordLoading || !isValid}
            >
              {passwordLoading ? (
                <span className="loading loading-spinner loading-sm" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              {passwordLoading ? '修改中...' : '修改密码'}
            </button>
          </form>
        </div>
      </div>

      {/* 两步验证提示 */}
      <div className="card bg-primary/5 border border-primary/20">
        <div className="card-body p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/20 rounded-lg">
                <Smartphone className="w-5 h-5 text-primary" />
              </div>
              <div>
                <p className="font-medium">两步验证</p>
                <p className="text-xs text-base-content/60">为账户添加额外的安全保护</p>
              </div>
            </div>
            <button className="btn btn-sm btn-outline btn-primary">
              启用
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}