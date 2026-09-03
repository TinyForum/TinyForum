import { useState } from "react";
import { useTranslations } from "next-intl";

// ---------- 配置弹窗（编辑节点的 params）----------
interface ConfigModalProps {
  title: string;
  params: Record<string, unknown>;
  onSave: (newParams: Record<string, unknown>) => void;
  onClose: () => void;
}

export function ConfigModal({
  title,
  params,
  onSave,
  onClose,
}: ConfigModalProps) {
  const t = useTranslations("Bot");
  const [config, setConfig] = useState<Record<string, unknown>>(params);

  const handleSave = () => {
    onSave(config);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-4 rounded shadow-lg w-96">
        <h3 className="font-bold text-lg mb-2">{t("config_node", { title })}</h3>
        <textarea
          className="w-full border p-2 mt-2 font-mono text-sm"
          rows={6}
          value={JSON.stringify(config, null, 2)}
          onChange={(e) => {
            try {
              const parsed = JSON.parse(e.target.value) as Record<
                string,
                unknown
              >;
              setConfig(parsed);
            } catch {
              // 保留原文本，不更新
            }
          }}
        />
        <div className="flex justify-end mt-4 space-x-2">
          <button
            className="px-3 py-1 border rounded hover:bg-gray-100"
            onClick={onClose}
          >
            {t("cancel")}
          </button>
          <button
            className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
            onClick={handleSave}
          >
            {t("save")}
          </button>
        </div>
      </div>
    </div>
  );
}
