'use client';

import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { userApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import { getErrorMessage } from '@/lib/utils';
import { Settings, Save } from 'lucide-react';
import Image from 'next/image';
import { useState } from 'react';

const schema = z.object({
  bio: z.string().max(500, '个人简介最多500字'),
  avatar: z.string().url('请输入有效的图片URL').optional().or(z.literal('')),
});

type SettingsForm = z.infer<typeof schema>;

export default function SettingsPage() {
  const { user, isAuthenticated, updateUser } = useAuthStore();
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!isAuthenticated) router.push('/auth/login');
  }, [isAuthenticated, router]);

  const { register, handleSubmit, watch, formState: { errors } } = useForm<SettingsForm>({
    resolver: zodResolver(schema),
    defaultValues: {
      bio: user?.bio || '',
      avatar: user?.avatar || '',
    },
  });

  const avatarValue = watch('avatar');

  const onSubmit = async (data: SettingsForm) => {
    setLoading(true);
    try {
      await userApi.updateProfile({ bio: data.bio, avatar: data.avatar || undefined });
      updateUser({ bio: data.bio, avatar: data.avatar || user?.avatar });
      toast.success('资料已更新');
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!user) return null;

  return (
    <div className="max-w-lg mx-auto">
      <h1 className="text-2xl font-bold flex items-center gap-2 mb-6">
        <Settings className="w-6 h-6 text-primary" /> 账号设置
      </h1>

      <div className="card bg-base-100 border border-base-300 shadow-sm">
        <div className="card-body p-6">
          {/* Avatar preview */}
          <div className="flex flex-col items-center gap-3 mb-6">
            <div className="avatar">
              <div className="w-24 h-24 rounded-2xl ring ring-primary ring-offset-2">
                <Image
                  src={avatarValue || user.avatar || `https://api.dicebear.com/8.x/initials/svg?seed=${user.username}`}
                  alt={user.username}
                  width={96}
                  height={96}
                  className="rounded-2xl object-cover"
                  onError={(e) => {
                    (e.target as HTMLImageElement).src = `https://api.dicebear.com/8.x/initials/svg?seed=${user.username}`;
                  }}
                />
              </div>
            </div>
            <div className="text-center">
              <p className="font-semibold">{user.username}</p>
              <p className="text-sm text-base-content/40">{user.email}</p>
            </div>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">头像图片 URL</span>
              </label>
              <input
                {...register('avatar')}
                type="text"
                placeholder="https://example.com/avatar.jpg"
                className={`input input-bordered focus:outline-none focus:border-primary ${errors.avatar ? 'input-error' : ''}`}
              />
              {errors.avatar && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error">{errors.avatar.message}</span>
                </label>
              )}
              <label className="label pt-1">
                <span className="label-text-alt text-base-content/40">留空则使用默认头像</span>
              </label>
            </div>

            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">个人简介</span>
              </label>
              <textarea
                {...register('bio')}
                rows={4}
                placeholder="介绍一下自己..."
                className={`textarea textarea-bordered focus:outline-none focus:border-primary resize-none ${errors.bio ? 'textarea-error' : ''}`}
              />
              {errors.bio && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error">{errors.bio.message}</span>
                </label>
              )}
            </div>

            <button type="submit" className="btn btn-primary w-full gap-2" disabled={loading}>
              {loading ? (
                <span className="loading loading-spinner loading-sm" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              保存设置
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
