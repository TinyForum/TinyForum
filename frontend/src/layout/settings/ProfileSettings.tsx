"use client";

import { useRef, useCallback, useMemo } from "react";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import { Save, AlertCircle, Camera } from "lucide-react";
import Avatar from "@/shared/ui/common/Avatar";
import { UserDO } from "@/shared/api/types/user.model.do";
import { useTranslations } from "next-intl";

import { useUpload } from "@/features/upload/hooks/useUpload";
import { useProfileSettings } from "@/features/user/hooks/useProfileSettings";

// 移除 avatar 的 URL 校验，允许任意字符串（包括相对路径、协议相对路径）

type ProfileForm = { bio?: string; avatar_url?: string };

interface ProfileSettingsProps {
  user: UserDO;
}

export default function ProfileSettings({ user }: ProfileSettingsProps) {
  const t = useTranslations("Settings");
  const { updateUser } = useAuthStore();

  const profileSchema = useMemo(
    () =>
      z.object({
        bio: z.string().max(500, t("bio_max_length")).optional().default(""),
        avatar_url: z.string().optional().default(""),
      }),
    [t],
  );
  const { uploadAvatar, isUploading, uploadError } = useUpload();
  const { updateProfile } = useProfileSettings();
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
      avatar_url: user?.avatar_url || "",
    },
    mode: "onChange",
  });
  const currentAvatar = useWatch({ control, name: "avatar_url" });
  const bioValue = useWatch({ control, name: "bio" });
  // const currentAvatar = watch("avatar");
  // const bioValue = watch("bio");

  // 处理文件选择并上传头像
  const handleFileSelect = useCallback(
    async (event: React.ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (!file) return;

      if (!file.type.startsWith("image/")) {
        toast.error(t("select_image_file"));
        return;
      }

      if (file.size > 5 * 1024 * 1024) {
        toast.error(t("image_size_limit"));
        return;
      }

      try {
        const uploadedUrl = await uploadAvatar(file);
        if (uploadedUrl) {
          setValue("avatar_url", uploadedUrl, {
            shouldDirty: true,
            shouldValidate: true,
          });
          toast.success(t("avatar_upload_success"));
        } else {
          throw new Error(t("upload_failed_no_url"));
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
      await updateProfile(data);
      updateUser({ ...user, ...data });
      toast.success(t("profile_updated"));
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const bioLength = bioValue?.length || 0;

  return (
    <div className="space-y-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">{t("profile")}</h1>
        <p className="text-sm text-base-content/60 mt-1">
          {t("profile_page_desc")}
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
                    avatarUrl={currentAvatar || user.avatar_url}
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
              {currentAvatar !== user?.avatar_url && currentAvatar && (
                <div className="absolute -top-2 -right-2 badge badge-primary badge-sm animate-pulse">
                  {t("not_saved")}
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
                <span className="label-text font-medium">{t("about_me")}</span>
                <span className="label-text-alt text-base-content/40">
                  {t("optional")}
                </span>
              </label>
              <textarea
                {...register("bio")}
                rows={4}
                placeholder={t("about_me_placeholder")}
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
                  {t("markdown_supported")}
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
              {isUploading ? t("uploading") : t("save_profile")}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
