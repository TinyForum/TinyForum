"use client";

import { useMemo } from "react";
import {
  Save,
  RefreshCw,
  Server,
  Globe,
  HardDrive,
  Puzzle,
  Layers,
  FileText,
  EyeOff,
} from "lucide-react";
import { Switch } from "@headlessui/react";
import type { ConfigField } from "../hooks/useSiteConfig";
import type { ConfigReloadResult } from "@/shared/api/modules/config";

interface SiteConfigPanelProps {
  groups: Record<string, ConfigField[]>;
  groupOrder: string[];
  isSaving: boolean;
  isReloading: boolean;
  isLoading: boolean;
  configFile: string;
  onChangeFile: (file: string) => void;
  onUpdateField: (key: string, value: unknown) => void;
  onSave: () => void;
  onReload: () => Promise<ConfigReloadResult | null>;
}

const GROUP_ICONS: Record<string, React.ReactNode> = {
  server: <Server className="w-4 h-4" />,
  api: <Globe className="w-4 h-4" />,
  frontend: <Globe className="w-4 h-4" />,
  log: <FileText className="w-4 h-4" />,
  attachment: <HardDrive className="w-4 h-4" />,
  plugins: <Puzzle className="w-4 h-4" />,
  general: <Layers className="w-4 h-4" />,
};

const GROUP_LABELS: Record<string, string> = {
  server: "服务器",
  api: "API 配置",
  frontend: "前端配置",
  log: "日志配置",
  attachment: "附件配置",
  plugins: "插件配置",
  allow_origins: "跨域来源",
  version: "版本",
  general: "通用配置",
};

export function SiteConfigPanel({
  groups,
  groupOrder,
  isSaving,
  isReloading,
  isLoading,
  configFile,
  onChangeFile,
  onUpdateField,
  onSave,
  onReload,
}: SiteConfigPanelProps) {
  const hasSensitive = useMemo(() => {
    for (const g of groupOrder) {
      for (const f of groups[g]) {
        if (f.sensitive) return true;
      }
    }
    return false;
  }, [groups, groupOrder]);

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-20 gap-4">
        <span className="loading loading-spinner loading-lg text-primary" />
        <p className="text-sm text-base-content/50">加载配置中...</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* 配置文件选择 */}
      <div className="flex items-center gap-3">
        <label className="text-sm font-medium">配置文件:</label>
        <select
          className="select select-bordered select-sm"
          value={configFile}
          onChange={(e) => onChangeFile(e.target.value)}
        >
          <option value="basic">basic.yml</option>
          <option value="app">app.yml</option>
        </select>
      </div>

      {groupOrder.length === 0 && (
        <div className="text-center py-12 text-base-content/50">
          暂无配置项
        </div>
      )}

      {groupOrder.map((groupKey) => {
        const groupFields = groups[groupKey];
        if (!groupFields || groupFields.length === 0) return null;
        return (
          <ConfigSection
            key={groupKey}
            title={GROUP_LABELS[groupKey] ?? groupKey}
            icon={GROUP_ICONS[groupKey]}
          >
            {groupFields.map((field) => (
              <ConfigFieldRow
                key={field.key}
                field={field}
                onChange={(v) => onUpdateField(field.key, v)}
              />
            ))}
          </ConfigSection>
        );
      })}

      {hasSensitive && (
        <div className="flex items-center gap-2 text-xs text-warning/80 p-3 bg-warning/5 rounded-lg border border-warning/20">
          <EyeOff className="w-3.5 h-3.5" />
          部分敏感字段已被隐藏，不参与保存。
        </div>
      )}

      {/* 操作按钮 */}
      <div className="flex justify-end gap-3 pt-2">
        <button
          onClick={onReload}
          disabled={isReloading}
          className="btn btn-outline gap-2 btn-sm"
        >
          <RefreshCw className={`w-4 h-4 ${isReloading ? "animate-spin" : ""}`} />
          {isReloading ? "重载中..." : "重载配置"}
        </button>
        <button
          onClick={onSave}
          disabled={isSaving}
          className="btn btn-primary gap-2 btn-sm"
        >
          {isSaving ? (
            <span className="loading loading-spinner loading-sm" />
          ) : (
            <Save className="w-4 h-4" />
          )}
          {isSaving ? "保存中..." : "保存配置"}
        </button>
      </div>
    </div>
  );
}

// ── Section ────────────────────────────────────────────────────────────

function ConfigSection({
  title,
  icon,
  children,
}: {
  title: string;
  icon?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <details className="collapse collapse-arrow bg-base-100 border border-base-300 shadow-sm" open>
      <summary className="collapse-title font-semibold flex items-center gap-2">
        {icon}
        {title}
      </summary>
      <div className="collapse-content space-y-3 pt-2">{children}</div>
    </details>
  );
}

// ── Field Row ──────────────────────────────────────────────────────────

function ConfigFieldRow({
  field,
  onChange,
}: {
  field: ConfigField;
  onChange: (value: unknown) => void;
}) {
  if (field.sensitive) return null;

  return (
    <div className="flex items-center justify-between gap-4 py-2 border-b border-base-200 last:border-b-0">
      <div className="flex-1 min-w-0">
        <label className="text-sm font-medium">{field.label}</label>
        <p className="text-xs text-base-content/40 truncate">{field.key}</p>
      </div>
      <div className="shrink-0">
        <FieldInput field={field} onChange={onChange} />
      </div>
    </div>
  );
}

// ── Type-specific inputs ───────────────────────────────────────────────

function FieldInput({
  field,
  onChange,
}: {
  field: ConfigField;
  onChange: (v: unknown) => void;
}) {
  const { type } = field;

  if (type === "boolean") {
    return (
      <Switch
        checked={Boolean(field.value)}
        onChange={onChange}
        className={`${
          field.value ? "bg-primary" : "bg-base-300"
        } relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary/50 focus:ring-offset-2`}
      >
        <span
          className={`${
            field.value ? "translate-x-6" : "translate-x-1"
          } inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform`}
        />
      </Switch>
    );
  }

  if (type === "number") {
    return (
      <input
        type="number"
        className="input input-bordered input-sm w-36 text-right"
        value={Number(field.value)}
        onChange={(e) => onChange(Number(e.target.value))}
      />
    );
  }

  if (type === "select" && field.options) {
    return (
      <select
        className="select select-bordered select-sm w-36"
        value={String(field.value)}
        onChange={(e) => onChange(e.target.value)}
      >
        {field.options.map((opt) => (
          <option key={opt} value={opt}>
            {opt}
          </option>
        ))}
      </select>
    );
  }

  if (type === "duration") {
    return (
      <input
        type="text"
        className="input input-bordered input-sm w-36 text-right font-mono"
        value={String(field.value)}
        onChange={(e) => onChange(e.target.value)}
        placeholder="30s / 5m / 24h"
      />
    );
  }

  if (type === "list") {
    const arr = Array.isArray(field.value) ? field.value : [];
    return (
      <input
        type="text"
        className="input input-bordered input-sm w-64 font-mono text-xs"
        value={arr.join(", ")}
        onChange={(e) => {
          const items = e.target.value
            .split(",")
            .map((s) => s.trim())
            .filter(Boolean);
          onChange(items);
        }}
        placeholder="逗号分隔"
      />
    );
  }

  return (
    <input
      type="text"
      className="input input-bordered input-sm"
      style={{ width: Math.max(180, Math.min(400, String(field.value).length * 10)) }}
      value={String(field.value ?? "")}
      onChange={(e) => onChange(e.target.value)}
    />
  );
}
