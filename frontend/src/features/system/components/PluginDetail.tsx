import { useTranslations } from "next-intl";
import { PluginMeta } from "@/shared/api/types/plugin.model";
import { DetailRow } from "./DetailRow";

export function PluginDetail({ plugin }: { plugin: PluginMeta }) {
  const t = useTranslations("System");

  return (
    <div className="px-4 pb-4 pt-0 bg-base-200/40 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm border-t border-base-200">
      <div className="space-y-2 pt-3">
        <DetailRow label={t("description")} value={plugin.description || "—"} />
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            {t("slots")}
          </span>
          <div className="flex flex-wrap gap-1 mt-1">
            {(plugin.slots ?? []).length === 0 ? (
              <span className="text-xs text-base-content/30 italic">
                {t("not_declared")}
              </span>
            ) : (
              (plugin.slots ?? []).map((s) => (
                <span key={s} className="badge badge-outline badge-sm">
                  {s}
                </span>
              ))
            )}
          </div>
        </div>
      </div>
      <div className="space-y-2 pt-3">
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            Script URL
          </span>
          <p className="mt-0.5 font-mono text-xs break-all text-base-content/50 bg-base-300 p-2 rounded">
            {plugin.scriptUrl}
          </p>
        </div>
        {plugin.updatedAt && (
          <DetailRow
            label={t("last_updated")}
            value={new Date(plugin.updatedAt).toLocaleString("zh-CN")}
          />
        )}
      </div>
    </div>
  );
}
