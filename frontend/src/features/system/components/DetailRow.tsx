import { useTranslations } from "next-intl";
import { Puzzle, Plus } from "lucide-react";

export function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <span className="text-xs text-base-content/40 uppercase tracking-wide">
        {label}
      </span>
      <p className="mt-0.5 text-base-content/70 text-xs">{value}</p>
    </div>
  );
}

export function EmptyState({
  keyword,
  onInstall,
}: {
  keyword: string;
  onInstall: () => void;
}) {
  const t = useTranslations("System");

  return (
    <div className="flex flex-col items-center justify-center h-40 gap-3 text-base-content/40">
      <Puzzle className="w-8 h-8" />
      <p className="text-sm">
        {keyword ? t("no_plugins_found", { keyword }) : t("no_plugins_installed")}
      </p>
      {!keyword && (
        <button onClick={onInstall} className="btn btn-ghost btn-xs gap-1">
          <Plus className="w-3 h-3" /> {t("install_first_plugin")}
        </button>
      )}
    </div>
  );
}
