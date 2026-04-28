"use client";

import { useState, useEffect } from "react";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { userApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/lib/utils";
import { Save, AlertCircle, Camera } from "lucide-react";
import Avatar from "@/components/user/Avatar";
import type { User } from "@/lib/api/types";

// 注意：User 接口中没有 displayName, location, website 字段
// 这些可能是扩展字段，或者需要使用其他字段名
const profileSchema = z.object({
  bio: z.string().max(500, "个人简介最多500字").optional().default(""),
  avatar: z.string().url("请输入有效的图片URL").optional().or(z.literal("")),
  // displayName: z.string().optional().default(""), // User 接口中没有此字段
  // location: z.string().optional().default(""),   // User 接口中没有此字段
  // website: z.string().url("请输入有效的网址").optional().or(z.literal("")), // User 接口中没有此字段
});

type ProfileForm = z.infer<typeof profileSchema>;

interface ProfileSettingsProps {
  user: User;
}

export default function ProfileSettings({ user }: ProfileSettingsProps) {
  const { updateUser } = useAuthStore();
  const [loading, setLoading] = useState<boolean>(false);
  const [previewAvatarUrl, setPreviewAvatarUrl] = useState<string>("");

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isDirty },
  } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      bio: user?.bio || "",
      avatar: user?.avatar || "",
    },
  });

  // 使用 useWatch 替代 watch
  const avatarValue = useWatch({
    control,
    name: "avatar",
    defaultValue: "",
  });

  const bioValue = useWatch({
    control,
    name: "bio",
    defaultValue: "",
  });

  useEffect(() => {
    if (avatarValue && avatarValue !== user?.avatar) {
      setPreviewAvatarUrl(avatarValue);
    } else {
      setPreviewAvatarUrl("");
    }
  }, [avatarValue, user?.avatar]);

  const onSubmit = async (data: ProfileForm): Promise<void> => {
    setLoading(true);
    try {
      await userApi.updateProfile(data);
      // 更新本地用户信息
      updateUser({ ...user, ...data });
      toast.success("资料已更新");
      setPreviewAvatarUrl("");
    } catch (err: unknown) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const bioLength = bioValue?.length || 0;

  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold">个人资料</h1>
        <p className="text-sm text-base-content/60 mt-1">
          管理您的个人信息和公开资料
        </p>
      </div>

      {/* 头像卡片 */}
      <div className="card bg-base-100 border border-base-200 shadow-sm">
        <div className="card-body p-6">
          <div className="flex flex-col items-center gap-4 pb-4 border-b border-base-200">
            <div className="relative group">
              <div className="avatar">
                <div className="w-28 h-28 rounded-2xl ring-2 ring-primary/20 ring-offset-2 transition-all group-hover:ring-primary/40">
                  <Avatar
                    avatarUrl={previewAvatarUrl || user.avatar}
                    username={user.username}
                    shape="square"
                  />
                </div>
              </div>
              <div className="absolute -bottom-2 -right-2 p-1.5 bg-primary rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                <Camera className="w-3 h-3" />
              </div>
              {previewAvatarUrl && (
                <div className="absolute -top-2 -right-2 badge badge-primary badge-sm animate-pulse">
                  预览
                </div>
              )}
            </div>
            <div className="text-center">
              <p className="font-semibold text-lg">@{user.username}</p>
              <p className="text-sm text-base-content/40">{user.email}</p>
            </div>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5 mt-4">
            {/* 头像URL */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">头像URL</span>
                <span className="label-text-alt text-base-content/40">
                  可选
                </span>
              </label>
              <input
                {...register("avatar")}
                type="text"
                placeholder="https://example.com/avatar.jpg"
                className={`input input-bordered focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all ${errors.avatar ? "input-error" : ""}`}
              />
              {errors.avatar && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error flex items-center gap-1">
                    <AlertCircle className="w-3 h-3" />
                    {errors.avatar.message}
                  </span>
                </label>
              )}
              <label className="label pt-1">
                <span className="label-text-alt text-base-content/40 flex items-center gap-1">
                  💡 提示：输入URL后会自动预览效果
                </span>
              </label>
            </div>

            {/* 个人简介 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">个人简介</span>
                <span className="label-text-alt text-base-content/40">
                  选填
                </span>
              </label>
              <textarea
                {...register("bio")}
                rows={4}
                placeholder="介绍一下自己吧..."
                className={`textarea textarea-bordered focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all resize-none ${errors.bio ? "textarea-error" : ""}`}
              />
              {errors.bio && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error flex items-center gap-1">
                    <AlertCircle className="w-3 h-3" />
                    {errors.bio.message}
                  </span>
                </label>
              )}
              <div className="flex justify-between items-center mt-1">
                <span className="text-xs text-base-content/40">
                  支持Markdown格式
                </span>
                <span
                  className={`text-xs ${bioLength > 450 ? "text-warning" : "text-base-content/40"}`}
                >
                  {bioLength}/500
                </span>
              </div>
            </div>

            {/* 提交按钮 */}
            <button
              type="submit"
              className="btn btn-primary w-full gap-2 transition-all transform active:scale-95"
              disabled={loading || !isDirty}
            >
              {loading ? (
                <span className="loading loading-spinner loading-sm" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              {loading ? "保存中..." : "保存资料"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}