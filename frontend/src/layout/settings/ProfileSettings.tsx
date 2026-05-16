"use client";

import { useRef, useCallback } from "react";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import { Save, AlertCircle, Camera } from "lucide-react";
import Avatar from "@/features/user/components/Avatar";
import { userApi } from "@/shared/api/modules/user";
import { UserDO } from "@/shared/api/types/user.model.do";

import { useUpload } from "@/features/upload/hooks/useUpload";

// 移除 avatar 的 URL 校验，允许任意字符串（包括相对路径、协议相对路径）
const profileSchema = z.object({
  bio: z.string().max(500, "个人简介最多500字").optional().default(""),
  avatar: z.string().optional().default(""),
});

type ProfileForm = z.infer<typeof profileSchema>;

interface ProfileSettingsProps {
  user: UserDO;
}

export default function ProfileSettings({ user }: ProfileSettingsProps) {
  const { updateUser } = useAuthStore();
  const { uploadAvatar, isUploading, uploadError } = useUpload();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const {
    register,
    handleSubmit,
    setValue,
    control,
    formState: { errors, isDirty },
  } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      bio: user?.bio || "",
      avatar: user?.avatar || "",
    },
    mode: "onChange",
  });
  const currentAvatar = useWatch({ control, name: "avatar" });
  const bioValue = useWatch({ control, name: "bio" });
  // const currentAvatar = watch("avatar");
  // const bioValue = watch("bio");

  // 处理文件选择并上传头像
  const handleFileSelect = useCallback(
    async (event: React.ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (!file) return;

      if (!file.type.startsWith("image/")) {
        toast.error("请选择图片文件");
        return;
      }

      if (file.size > 5 * 1024 * 1024) {
        toast.error("图片大小不能超过 5MB");
        return;
      }

      try {
        const uploadedUrl = await uploadAvatar(file);
        if (uploadedUrl) {
          setValue("avatar", uploadedUrl, {
            shouldDirty: true,
            shouldValidate: true,
          });
          toast.success("头像上传成功，请保存资料");
        } else {
          throw new Error("上传失败，未返回URL");
        }
      } catch (err) {
        toast.error(getErrorMessage(err));
      } finally {
        if (fileInputRef.current) fileInputRef.current.value = "";
      }
    },
    [uploadAvatar, setValue],
  );

  const onSubmit = async (data: ProfileForm) => {
    try {
      await userApi.updateProfile(data);
      updateUser({ ...user, ...data });
      toast.success("资料已更新");
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const bioLength = bioValue?.length || 0;

  return (
    <div className="space-y-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">个人资料</h1>
        <p className="text-sm text-base-content/60 mt-1">
          管理您的个人信息和公开资料
        </p>
      </div>

      <div className="card bg-base-100 border border-base-200 shadow-sm">
        <div className="card-body p-6">
          {/* 头像区域 */}
          <div className="flex flex-col items-center gap-4 pb-4 border-b border-base-200">
            <div className="relative group">
              <div className="avatar">
                <div className="w-28 h-28 rounded-2xl ring-2 ring-primary/20 ring-offset-2 transition-all group-hover:ring-primary/40">
                  <Avatar
                    avatarUrl={currentAvatar || user.avatar}
                    username={user.username}
                    shape="square"
                  />
                </div>
              </div>
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className="absolute -bottom-2 -right-2 p-1.5 bg-primary rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer hover:bg-primary-focus"
                disabled={isUploading}
              >
                {isUploading ? (
                  <span className="loading loading-spinner loading-xs" />
                ) : (
                  <Camera className="w-3 h-3" />
                )}
              </button>
              {currentAvatar !== user?.avatar && currentAvatar && (
                <div className="absolute -top-2 -right-2 badge badge-primary badge-sm animate-pulse">
                  未保存
                </div>
              )}
            </div>
            <div className="text-center">
              <p className="font-semibold text-lg">@{user.username}</p>
              <p className="text-sm text-base-content/40">{user.email}</p>
            </div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileSelect}
              className="hidden"
            />
            {uploadError && (
              <div className="text-error text-xs flex items-center gap-1">
                <AlertCircle className="w-3 h-3" />
                {uploadError}
              </div>
            )}
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5 mt-4">
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
                  className={`text-xs ${
                    bioLength > 450 ? "text-warning" : "text-base-content/40"
                  }`}
                >
                  {bioLength}/500
                </span>
              </div>
            </div>

            {/* 提交按钮 */}
            <button
              type="submit"
              className="btn btn-primary w-full gap-2 transition-all transform active:scale-95"
              disabled={isUploading || !isDirty}
            >
              {isUploading ? (
                <span className="loading loading-spinner loading-sm" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              {isUploading ? "上传中..." : "保存资料"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
