import { useState } from "react";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { authApi, userAPI } from "@/shared/api";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import {
  Save,
  KeyRound,
  Eye,
  EyeOff,
  AlertCircle,
  CheckCircle,
} from "lucide-react";
import { usePasswordStrength } from "@/features/settings/hooks/usePasswordStrength";
// 修改密码表单的 Schema 和类型
const passwordSchema = z
  .object({
    oldPassword: z.string().min(1, "请输入旧密码"),
    newPassword: z.string().min(6, "密码长度至少6位"),
    confirmPassword: z.string().min(1, "请确认新密码"),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "两次输入的密码不一致",
    path: ["confirmPassword"],
  });

type PasswordForm = z.infer<typeof passwordSchema>;
export function ChangePasswordForm() {
  const router = useRouter();
  const [passwordLoading, setPasswordLoading] = useState<boolean>(false);
  const [showOldPassword, setShowOldPassword] = useState<boolean>(false);
  const [showNewPassword, setShowNewPassword] = useState<boolean>(false);
  const [showConfirmPassword, setShowConfirmPassword] =
    useState<boolean>(false);

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

  const onSubmit = async (data: PasswordForm): Promise<void> => {
    setPasswordLoading(true);
    try {
      await userAPI.changePassword({
        old_password: data.oldPassword,
        new_password: data.newPassword,
      });
      toast.success("密码修改成功");
      reset();
      setTimeout(() => {
        toast.success("即将退出，请使用新密码重新登录");
        authApi.logout();
        router.push("/auth/login");
      }, 3000);
    } catch (err: unknown) {
      toast.error(getErrorMessage(err));
    } finally {
      setPasswordLoading(false);
    }
  };

  return (
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
                {...register("oldPassword")}
                type={showOldPassword ? "text" : "password"}
                placeholder="请输入当前密码"
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
              <span className="label-text font-medium">新密码</span>
            </label>
            <div className="relative">
              <input
                {...register("newPassword")}
                type={showNewPassword ? "text" : "password"}
                placeholder="请输入新密码"
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
                  <span>密码强度</span>
                  <span
                    className={`font-medium ${
                      passwordStrength.level === "强"
                        ? "text-green-600"
                        : passwordStrength.level === "中"
                          ? "text-yellow-600"
                          : "text-red-600"
                    }`}
                  >
                    {passwordStrength.level}
                  </span>
                </div>
                <div className="w-full h-1.5 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className={`h-full transition-all duration-300 ${passwordStrength.color}`}
                    style={{ width: `${(passwordStrength.score / 6) * 100}%` }}
                  />
                </div>
                <p className="text-xs text-base-content/60">
                  {passwordStrength.message}
                </p>
              </div>
            )}

            {errors.newPassword && (
              <p className="text-error text-xs mt-1">
                {errors.newPassword.message}
              </p>
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
                {...register("confirmPassword")}
                type={showConfirmPassword ? "text" : "password"}
                placeholder="请再次输入新密码"
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
            {passwordLoading ? "修改中..." : "修改密码"}
          </button>
        </form>
      </div>
    </div>
  );
}
