// app/[locale]/auth/forgot-password/page.tsx
"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import toast from "react-hot-toast";
import { Mail, ArrowLeft, Send } from "lucide-react";
import { useTranslations, } from "next-intl";

const forgotPasswordSchema = z.object({
  email: z.string().email("请输入有效的邮箱"),
});

type ForgotPasswordForm = z.infer<typeof forgotPasswordSchema>;

export default function ForgotPasswordPage() {
  const t = useTranslations("Auth");

  const [isLoading, setIsLoading] = useState(false);
  const [isEmailSent, setIsEmailSent] = useState(false);

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

  const onSubmit = async () => {
    setIsLoading(true);

    try {
      // TODO: 调用忘记密码 API
      // const response = await fetch('/api/auth/forgot-password', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ email: data.email }),
      // });

      // 模拟 API 调用
      await new Promise((resolve) => setTimeout(resolve, 1500));

      // 假设总是成功（实际应该根据 API 响应判断）
      setIsEmailSent(true);
      toast.success("重置密码链接已发送到您的邮箱");
    } catch {
      toast.error("发送失败，请稍后重试");
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendEmail = () => {
    const email = getValues("email");
    if (email) {
      onSubmit();
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
              <h1 className="text-2xl font-bold">{t("forgot_password")}</h1>
              <p className="text-base-content/50 text-sm mt-1">
                {isEmailSent ? "请查收邮件重置密码" : "输入邮箱以重置密码"}
              </p>
            </div>

            {!isEmailSent ? (
              // 表单页面
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                {/* Email */}
                <div className="form-control">
                  <label className="label pb-1">
                    <span className="label-text font-medium">{t("email")}</span>
                  </label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                    <input
                      {...register("email")}
                      type="email"
                      placeholder="your@email.com"
                      className={`input input-bordered w-full pl-10 focus:outline-none focus:border-primary ${errors.email ? "input-error" : ""}`}
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
                  {isLoading ? "发送中..." : t("send_reset_link")}
                </button>

                {/* Back to Login */}
                <div className="text-center">
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
              // 成功页面
              <div className="space-y-6">
                <div className="alert alert-success">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="stroke-current shrink-0 h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>重置密码链接已发送到您的邮箱</span>
                </div>

                <div className="bg-base-200 rounded-lg p-4 text-sm space-y-2">
                  <p className="text-base-content/70">📧 请检查您的邮箱：</p>
                  <p className="font-mono font-medium break-all">
                    {getValues("email")}
                  </p>
                  <p className="text-base-content/60 text-xs mt-2">
                    如果没有收到邮件，请检查垃圾邮件箱，或尝试重新发送。
                  </p>
                </div>

                <div className="flex flex-col gap-3">
                  <button
                    onClick={handleResendEmail}
                    className="btn btn-outline btn-sm w-full gap-2"
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <span className="loading loading-spinner loading-xs" />
                    ) : (
                      <Send className="w-3 h-3" />
                    )}
                    重新发送邮件
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
