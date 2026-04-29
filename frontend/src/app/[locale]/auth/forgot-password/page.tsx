// app/[locale]/auth/forgot-password/page.tsx
"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import toast from "react-hot-toast";
import { Mail, ArrowLeft, Send, CheckCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import { authApi } from "@/shared/api";

const forgotPasswordSchema = z.object({
  email: z.string().email("请输入有效的邮箱"),
});

type ForgotPasswordForm = z.infer<typeof forgotPasswordSchema>;

export default function ForgotPasswordPage() {
  const t = useTranslations("Auth");

  const [isLoading, setIsLoading] = useState(false);
  const [isEmailSent, setIsEmailSent] = useState(false);
  const [submittedEmail, setSubmittedEmail] = useState("");
  console.log(submittedEmail);

  const {
    register,
    handleSubmit,
    formState: { errors },
    getValues,
  } = useForm<ForgotPasswordForm>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: ForgotPasswordForm) => {
    setIsLoading(true);

    try {
      // 调用后端 API（始终返回成功，不暴露邮箱是否存在）
      const response = await authApi.forgotPassword({ email: data.email });

      // 后端始终返回 200，无论邮箱是否存在
      if (response.data.code === 0 || response.status === 200) {
        setSubmittedEmail(data.email);
        setIsEmailSent(true);
        // 不显示具体的成功消息，使用统一的提示
        toast.success(t("reset_link_sent") || "重置密码链接已发送");
      } else {
        // 理论上不会进入这里，但做个兜底
        toast.error(t("send_failed") || "发送失败，请稍后重试");
      }
    } catch (error) {
      // 网络错误或其他异常才显示错误
      console.error("Forgot password error:", error);
      toast.error(t("network_error") || "网络错误，请稍后重试");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendEmail = () => {
    const email = getValues("email");
    if (email && !isLoading) {
      onSubmit({ email });
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-100 shadow-xl border border-base-300">
          <div className="card-body p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-2xl font-black mx-auto mb-4">
                B
              </div>
              <h1 className="text-2xl font-bold">
                {!isEmailSent ? t("forgot_password") : t("check_email")}
              </h1>
              <p className="text-base-content/50 text-sm mt-1">
                {!isEmailSent
                  ? t("forgot_password_desc")
                  : t("reset_link_sent_desc")}
              </p>
            </div>

            {!isEmailSent ? (
              // 表单页面
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                {/* Email Input */}
                <div className="form-control">
                  <label className="label pb-1">
                    <span className="label-text font-medium">
                      {t("email_address")}
                    </span>
                  </label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                    <input
                      {...register("email")}
                      type="email"
                      placeholder="your@email.com"
                      className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${
                        errors.email ? "input-error" : ""
                      }`}
                      autoComplete="email"
                      disabled={isLoading}
                    />
                  </div>
                  {errors.email && (
                    <label className="label pt-1">
                      <span className="label-text-alt text-error">
                        {errors.email.message}
                      </span>
                    </label>
                  )}
                </div>

                {/* Security Note */}
                <div className="text-xs text-base-content/50 text-center">
                  <p>{t("security_note")}</p>
                </div>

                {/* Submit Button */}
                <button
                  type="submit"
                  className="btn btn-primary w-full gap-2"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <span className="loading loading-spinner loading-sm" />
                  ) : (
                    <Send className="w-4 h-4" />
                  )}
                  {isLoading ? t("sending") : t("send_reset_link")}
                </button>

                {/* Back to Login */}
                <div className="text-center pt-2">
                  <Link
                    href="/auth/login"
                    className="inline-flex items-center gap-1 text-sm text-primary hover:underline"
                  >
                    <ArrowLeft className="w-3 h-3" />
                    {t("back_to_login")}
                  </Link>
                </div>
              </form>
            ) : (
              // 成功页面（统一消息，不暴露邮箱是否存在）
              <div className="space-y-6">
                {/* Success Alert */}
                <div className="alert alert-success bg-success/10 border-success/20">
                  <CheckCircle className="h-5 w-5 text-success" />
                  <div className="flex-1">
                    <span className="font-medium">{t("email_sent")}</span>
                    <p className="text-sm opacity-90 mt-0.5">
                      {t("check_your_email")}
                    </p>
                  </div>
                </div>

                {/* Instructions */}
                <div className="bg-base-200 rounded-lg p-5 text-sm space-y-3">
                  <p className="text-base-content/80 font-medium">
                    📧 {t("what_to_do_next")}:
                  </p>
                  <ol className="space-y-2 text-base-content/70 list-decimal list-inside">
                    <li>{t("step1_check_email")}</li>
                    <li>{t("step2_click_link")}</li>
                    <li>{t("step3_reset_password")}</li>
                  </ol>
                  <div className="mt-3 pt-2 border-t border-base-300">
                    <p className="text-base-content/60 text-xs">
                      💡 {t("email_tip")}
                    </p>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex flex-col gap-3">
                  <button
                    onClick={handleResendEmail}
                    className="btn btn-outline w-full gap-2"
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <span className="loading loading-spinner loading-sm" />
                    ) : (
                      <Send className="w-4 h-4" />
                    )}
                    {t("resend_email")}
                  </button>

                  <Link
                    href="/auth/login"
                    className="btn btn-primary w-full gap-2"
                  >
                    {t("back_to_login")}
                  </Link>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
