import { useTranslations } from "next-intl";

// 空状态组件
export function EmptyState() {
  const t = useTranslations("Leaderboard");
  return (
    <div className="bg-base-100 rounded-2xl shadow-sm p-12 text-center border border-base-200">
      <div className="text-6xl mb-4 opacity-50">📊</div>
      <h3 className="text-lg font-semibold text-base-content mb-2">
        {t("no_data")}
      </h3>
      <p className="text-base-content/60 text-sm">{t("no_data_description")}</p>
    </div>
  );
}
