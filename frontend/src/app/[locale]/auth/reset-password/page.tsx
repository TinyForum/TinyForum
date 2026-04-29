// app/[locale]/auth/reset-password/page.tsx

"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useLocale } from "next-intl";

import {
  Lock,
  Eye,
  EyeOff,
  KeyRound,
  ArrowLeft,
  CheckCircle,
} from "lucide-react";
import Link from "next/link";
import { useValidateTokenStore, useResetPasswordStore } from "@/store/token";

export default function ResetPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const locale = useLocale();

  const tokenFromUrl = searchParams.get("token"); // 从 URL 获取 token

  // ✅ 从 store 获取状态和方法
  const {
    isValid,
    isLoading: isValidating,
    error: validationError,
    validateToken,
    resetValidation,
  } = useValidateTokenStore();

  const {
    token, // ✅ 从 store 获取 token
    password,
    confirmPassword,
    isLoading: isResetting,
    isSuccess,
    error: resetError,
    setToken, // ✅ 获取 setToken 方法
    setPassword,
    setConfirmPassword,
    resetPassword,
    resetForm: resetResetForm,
  } = useResetPasswordStore();

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  // ✅ 将 URL 中的 token 设置到 store
  useEffect(() => {
    if (tokenFromUrl) {
      console.log("Setting token to store:", tokenFromUrl);
      setToken(tokenFromUrl); // 关键：将 token 设置到 store
    }
  }, [tokenFromUrl, setToken]);

  // 验证 token
  useEffect(() => {
    console.log("=== Effect triggered ===");
    console.log("Token from URL:", tokenFromUrl);

    if (tokenFromUrl) {
      console.log("Calling validateToken...");
      validateToken(tokenFromUrl).then((result) => {
        console.log("Validation promise resolved:", result);
      });
    } else {
      console.log("No token, resetting validation");
      resetValidation();
    }

    return () => {
      console.log("Cleanup");
      resetValidation();
      resetResetForm();
    };
  }, [tokenFromUrl]); // 依赖 tokenFromUrl

  // 监听状态变化
  useEffect(() => {
    console.log("Store token value:", token);
    console.log("IsValid:", isValid);
  }, [token, isValid]);

  // 验证中
  if (isValidating) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  // Token 无效
  if (isValid === false) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center px-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-100 shadow-xl border border-base-300">
            <div className="card-body p-8 text-center">
              <div className="text-6xl mb-4">🔒</div>
              <h2 className="text-2xl font-bold mb-2">Invalid Reset Link</h2>
              <p className="text-base-content/70 mb-6">
                {validationError ||
                  "This password reset link is invalid or has expired."}
              </p>
              <Link
                href={`/${locale}/auth/forgot-password`}
                className="btn btn-primary"
              >
                Request New Reset Link
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // 等待验证
  if (isValid === null) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center">
        <div className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  // 重置成功
  if (isSuccess) {
    return (
      <div className="min-h-[80vh] flex items-center justify-center px-4">
        <div className="w-full max-w-md">
          <div className="card bg-base-100 shadow-xl border border-base-300">
            <div className="card-body p-8 text-center">
              <div className="text-6xl mb-4">
                <CheckCircle className="w-16 h-16 text-success mx-auto" />
              </div>
              <h2 className="text-2xl font-bold mb-2">
                Password Reset Success!
              </h2>
              <p className="text-base-content/70 mb-6">
                Your password has been reset successfully. Please login with
                your new password.
              </p>
              <Link
                href={`/${locale}/auth/login`}
                className="btn btn-primary w-full"
              >
                Go to Login
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // 显示重置表单
  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <div className="card bg-base-100 shadow-xl border border-base-300">
          <div className="card-body p-8">
            <div className="text-center mb-8">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-2xl font-black mx-auto mb-4">
                B
              </div>
              <h1 className="text-2xl font-bold">Set New Password</h1>
              <p className="text-base-content/50 text-sm mt-1">
                Please enter your new password
              </p>
            </div>

            <form
              onSubmit={async (e) => {
                e.preventDefault();
                console.log("Submitting reset password...");
                console.log("Current token in store:", token); // 调试日志
                const result = await resetPassword();
                console.log("Reset result:", result);
                if (result.success) {
                  setTimeout(() => {
                    router.push(`/${locale}/auth/login`);
                  }, 2000);
                }
              }}
              className="space-y-4"
            >
              {/* New Password */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">New Password</span>
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    type={showPassword ? "text" : "password"}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    className="input input-bordered w-full pl-10 pr-10 focus:outline-none focus:border-primary"
                    disabled={isResetting}
                    required
                    minLength={6}
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-base-content/40 hover:text-base-content"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                </div>
              </div>

              {/* Confirm Password */}
              <div className="form-control">
                <label className="label pb-1">
                  <span className="label-text font-medium">
                    Confirm Password
                  </span>
                </label>
                <div className="relative">
                  <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                  <input
                    type={showConfirmPassword ? "text" : "password"}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="••••••••"
                    className="input input-bordered w-full pl-10 pr-10 focus:outline-none focus:border-primary"
                    disabled={isResetting}
                    required
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-base-content/40 hover:text-base-content"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  >
                    {showConfirmPassword ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                </div>
              </div>

              {resetError && (
                <div className="alert alert-error text-sm p-3">
                  {resetError}
                </div>
              )}

              <button
                type="submit"
                className="btn btn-primary w-full gap-2 mt-4"
                disabled={isResetting}
              >
                {isResetting ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  <KeyRound className="w-4 h-4" />
                )}
                {isResetting ? "Resetting..." : "Reset Password"}
              </button>

              <div className="text-center">
                <Link
                  href={`/${locale}/auth/login`}
                  className="inline-flex items-center gap-1 text-sm text-primary hover:underline"
                >
                  <ArrowLeft className="w-3 h-3" />
                  Back to Login
                </Link>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}
