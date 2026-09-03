// components/question/RewardScoreInput.tsx
"use client";

import { AskFormData } from "@/features/qustion/hooks/useQuestionForm";
import {
  CurrencyDollarIcon,
  SparklesIcon,
  InformationCircleIcon,
} from "@heroicons/react/24/outline";
import { useTranslations } from "next-intl";
import { UseFormRegister } from "react-hook-form";

interface RewardScoreInputProps {
  register: UseFormRegister<AskFormData>;
  rewardScore: number;
  userScore?: number;
  error?: string;
}

export function RewardScoreInput({
  register,
  rewardScore,
  userScore = 0,
  error,
}: RewardScoreInputProps) {
  const t = useTranslations("Questions");
  // 计算是否超出可用积分
  const isExceeded = rewardScore > userScore;
  // 预设悬赏金额
  const presetScores = [0, 10, 20, 50, 100];

  return (
    <div className="card bg-gradient-to-r from-amber-50 to-orange-50 dark:from-amber-900/10 dark:to-orange-900/10 border border-amber-200 dark:border-amber-800/30 shadow-sm">
      <div className="card-body p-5">
        {/* 标题 */}
        <div className="flex items-center gap-2 mb-3">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-amber-500 to-orange-500 flex items-center justify-center shadow-sm">
            <CurrencyDollarIcon className="w-4 h-4 text-white" />
          </div>
          <div>
            <h3 className="font-semibold text-base-content">{t("reward_score_title")}</h3>
            <p className="text-xs text-base-content/60">
              {t("reward_score_desc")}
            </p>
          </div>
        </div>

        {/* 输入区域 */}
        <div className="space-y-3">
          {/* 预设金额快捷按钮 */}
          <div className="flex flex-wrap gap-2">
            {presetScores.map((score) => (
              <button
                key={score}
                type="button"
                onClick={() => {
                  // 需要通过 react-hook-form 设置值
                  const input = document.querySelector(
                    'input[name="reward_score"]',
                  ) as HTMLInputElement;
                  if (input) {
                    input.value = String(score);
                    input.dispatchEvent(new Event("change", { bubbles: true }));
                  }
                }}
                className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-all ${
                  rewardScore === score
                    ? "bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-md"
                    : "bg-base-100 border border-base-200 text-base-content/70 hover:border-amber-300 hover:text-amber-600 dark:hover:border-amber-700"
                }`}
              >
                {score === 0 ? t("no_reward") : t("score_points", { score })}
              </button>
            ))}
          </div>

          {/* 自定义输入 */}
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3">
            <div className="relative flex-1">
              <CurrencyDollarIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
              <input
                {...register("reward_score", {
                  min: { value: 0, message: t("reward_no_negative") },
                  max: {
                    value: Math.min(100, userScore),
                    message: t("reward_max", { max: Math.min(100, userScore) }),
                  },
                  valueAsNumber: true,
                })}
                type="number"
                min="0"
                max={Math.min(100, userScore)}
                step="5"
                className={`w-full pl-9 pr-3 py-2 rounded-lg border bg-base-100 text-base-content placeholder-base-content/40 focus:outline-none focus:ring-2 focus:ring-primary/50 transition-all ${
                  error || isExceeded
                    ? "border-red-500 focus:ring-red-500"
                    : "border-base-200 focus:border-primary"
                }`}
                placeholder={t("reward_custom")}
              />
            </div>

            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-base-100 border border-base-200">
              <SparklesIcon className="w-3.5 h-3.5 text-amber-500" />
              <span className="text-xs text-base-content/70">
                {t("available_score")}
                <span
                  className={`font-semibold ml-1 ${userScore >= rewardScore ? "text-primary" : "text-red-500"}`}
                >
                  {userScore}
                </span>
              </span>
            </div>
          </div>

          {/* 错误提示 */}
          {(error || isExceeded) && (
            <div className="alert alert-error alert-sm p-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800/30">
              <InformationCircleIcon className="w-4 h-4 text-red-500" />
              <span className="text-xs text-red-600 dark:text-red-400">
                {error ||
                  t("reward_insufficient", { need: rewardScore, have: userScore })}
              </span>
            </div>
          )}

          {/* 提示信息 */}
          {rewardScore > 0 && !isExceeded && !error && (
            <div className="alert alert-info alert-sm p-2 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800/30">
              <SparklesIcon className="w-4 h-4 text-blue-500" />
              <div className="text-xs text-blue-600 dark:text-blue-400">
                <p>{t("reward_exposure", { score: rewardScore })}</p>
                <p className="text-xs opacity-75 mt-0.5">
                  {t("reward_deduct", { score: rewardScore })}
                </p>
              </div>
            </div>
          )}

          {rewardScore === 0 && (
            <div className="alert alert-ghost alert-sm p-2 bg-base-200/50 border border-base-200">
              <InformationCircleIcon className="w-4 h-4 text-base-content/40" />
              <span className="text-xs text-base-content/60">
                {t("reward_hint")}
              </span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
