import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Save,
  KeyRound,
  Eye,
  EyeOff,
  AlertCircle,
  CheckCircle,
} from "lucide-react";
import { useTranslations } from "next-intl";
import { usePasswordStrength } from "@/features/settings/hooks/usePasswordStrength";
import { useState, useMemo } from "react";
import { useChangePassword } from "@/store/password";

type PasswordForm = { oldPassword: string; newPassword: string; confirmPassword: string };

export function ChangePasswordForm() {
  const t = useTranslations("Settings");

  const passwordSchema = useMemo(
    () =>
      z
        .object({
          oldPassword: z.string().min(1, t("enter_old_password")),
          newPassword: z.string().min(8, t("password_length_min_eight")),
          confirmPassword: z.string().min(1, t("please_confirm_new_password")),
        })
        .refine((data) => data.newPassword === data.confirmPassword, {
          message: t("passwords_not_match"),
          path: ["confirmPassword"],
        }),
    [t],
  );
  const [showOldPassword, setShowOldPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const { changePassword, isLoading } = useChangePassword();

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors, isValid },
  } = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
    mode: "onChange",
    defaultValues: {
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
  });

  const newPasswordValue = useWatch({
    control,
    name: "newPassword",
    defaultValue: "",
  });
  const confirmPasswordValue = useWatch({
    control,
    name: "confirmPassword",
    defaultValue: "",
  });
  const passwordStrength = usePasswordStrength(newPasswordValue || "");

  const onSubmit = async (data: PasswordForm) => {
    const success = await changePassword(data.oldPassword, data.newPassword);
    if (success) {
      reset();
    }
  };

  return (
    <div className="card bg-base-100 border border-base-200 shadow-sm">
      <div className="card-body p-6">
        <div className="flex items-center gap-2 mb-4">
          <KeyRound className="w-5 h-5 text-primary" />
          <h2 className="text-xl font-semibold">{t("change_password")}</h2>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* 旧密码 */}
          <div className="form-control">
            <label className="label pb-1">
              <span className="label-text font-medium">{t("current_password")}</span>
            </label>
            <div className="relative">
              <input
                {...register("oldPassword")}
                type={showOldPassword ? "text" : "password"}
                placeholder={t("enter_current_password_placeholder")}
                className={`input input-bordered w-full pr-10 ${errors.oldPassword ? "input-error" : ""}`}
              />
              <button
                type="button"
                onClick={() => setShowOldPassword(!showOldPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2"
              >
                {showOldPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
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
              <span className="label-text font-medium">{t("new_password")}</span>
            </label>
            <div className="relative">
              <input
                {...register("newPassword")}
                type={showNewPassword ? "text" : "password"}
                placeholder={t("enter_new_password_placeholder")}
                className={`input input-bordered w-full pr-10 ${errors.newPassword ? "input-error" : ""}`}
              />
              <button
                type="button"
                onClick={() => setShowNewPassword(!showNewPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2"
              >
                {showNewPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>

            {/* 密码强度显示 */}
            {newPasswordValue && (
              <div className="mt-2 space-y-1">
                <div className="flex justify-between text-xs">
                  <span>{t("password_strength")}</span>
                  <span
                    className={`font-medium ${
                      passwordStrength.level === "password_strength_strong"
                        ? "text-green-600"
                        : passwordStrength.level === "password_strength_medium"
                          ? "text-yellow-600"
                          : "text-red-600"
                    }`}
                  >
                    {t(passwordStrength.level)}
                  </span>
                </div>
                <div className="w-full h-1.5 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className={`h-full transition-all duration-300 ${passwordStrength.color}`}
                    style={{ width: `${(passwordStrength.score / 6) * 100}%` }}
                  />
                </div>
                <p className="text-xs text-base-content/60">
                  {passwordStrength.message && t(passwordStrength.message)}
                </p>
              </div>
            )}

            {errors.newPassword && (
              <p className="text-error text-xs mt-1">
                {errors.newPassword.message}
              </p>
            )}
            <p className="text-xs text-base-content/40 mt-1">
              {t("password_suggestion")}
            </p>
          </div>

          {/* 确认密码 */}
          <div className="form-control">
            <label className="label pb-1">
              <span className="label-text font-medium">{t("confirm_new_password")}</span>
            </label>
            <div className="relative">
              <input
                {...register("confirmPassword")}
                type={showConfirmPassword ? "text" : "password"}
                placeholder={t("reenter_new_password_placeholder")}
                className={`input input-bordered w-full pr-10 ${errors.confirmPassword ? "input-error" : ""}`}
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2"
              >
                {showConfirmPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>
            {errors.confirmPassword && (
              <p className="text-error text-xs mt-1">
                {errors.confirmPassword.message}
              </p>
            )}
            {confirmPasswordValue &&
              !errors.confirmPassword &&
              newPasswordValue === confirmPasswordValue && (
                <p className="text-success text-xs mt-1 flex items-center gap-1">
                  <CheckCircle className="w-3 h-3" />
                  {t("password_match")}
                </p>
              )}
          </div>

          <button
            type="submit"
            className="btn btn-primary w-full gap-2"
            disabled={isLoading || !isValid}
          >
            {isLoading ? (
              <span className="loading loading-spinner loading-sm" />
            ) : (
              <Save className="w-4 h-4" />
            )}
            {isLoading ? t("changing") : t("change_password")}
          </button>
        </form>
      </div>
    </div>
  );
}
