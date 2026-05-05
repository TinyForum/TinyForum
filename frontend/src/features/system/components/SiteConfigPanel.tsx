"use client";

import { Save, Globe, Mail, FileText, Hash, Layout, Moon } from "lucide-react";
import { Switch } from "@headlessui/react";
import type { SiteConfig } from "../hooks/useSiteConfig";

interface SiteConfigPanelProps {
  config: SiteConfig;
  isSaving: boolean;
  update: <K extends keyof SiteConfig>(key: K, value: SiteConfig[K]) => void;
  onSave: () => void;
}

export function SiteConfigPanel({
  config,
  isSaving,
  update,
  onSave,
}: SiteConfigPanelProps) {
  return (
    <div className="space-y-6">
      <SectionCard title="基本信息" subtitle="Basic Information">
        <Field label="网站名称" icon={<Globe className="w-4 h-4" />}>
          <input
            type="text"
            className="input input-bordered w-full"
            value={config.siteName}
            onChange={(e) => update("siteName", e.target.value)}
            placeholder="我的网站"
          />
        </Field>

        <Field label="网站描述" icon={<FileText className="w-4 h-4" />}>
          <textarea
            className="textarea textarea-bordered w-full resize-none"
            rows={3}
            value={config.siteDescription}
            onChange={(e) => update("siteDescription", e.target.value)}
            placeholder="一句话介绍你的网站..."
          />
        </Field>

        <Field label="SEO 关键词" icon={<Hash className="w-4 h-4" />}>
          <input
            type="text"
            className="input input-bordered w-full"
            value={config.siteKeywords}
            onChange={(e) => update("siteKeywords", e.target.value)}
            placeholder="关键词1,关键词2,关键词3"
          />
          <p className="text-xs text-base-content/40 mt-1">
            用英文逗号分隔多个关键词
          </p>
        </Field>

        <Field label="管理员邮箱" icon={<Mail className="w-4 h-4" />}>
          <input
            type="email"
            className="input input-bordered w-full"
            value={config.adminEmail}
            onChange={(e) => update("adminEmail", e.target.value)}
          />
        </Field>
      </SectionCard>

      <SectionCard title="显示偏好" subtitle="Display Preferences">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <Field label="每页显示数量" icon={<Layout className="w-4 h-4" />}>
            <select
              className="select select-bordered w-full"
              value={config.itemsPerPage}
              onChange={(e) => update("itemsPerPage", Number(e.target.value))}
            >
              {[10, 20, 50, 100].map((n) => (
                <option key={n} value={n}>
                  {n} 条/页
                </option>
              ))}
            </select>
          </Field>

          <Field label="默认主题" icon={<Moon className="w-4 h-4" />}>
            <select
              className="select select-bordered w-full"
              value={config.defaultTheme}
              onChange={(e) =>
                update(
                  "defaultTheme",
                  e.target.value as SiteConfig["defaultTheme"],
                )
              }
            >
              <option value="light">浅色模式</option>
              <option value="dark">深色模式</option>
              <option value="auto">跟随系统</option>
            </select>
          </Field>
        </div>
      </SectionCard>

      <SectionCard title="高级选项" subtitle="Advanced">
        <div className="flex items-center justify-between p-4 rounded-xl bg-warning/5 border border-warning/20">
          <div>
            <p className="font-medium text-sm">维护模式</p>
            <p className="text-xs text-base-content/50 mt-0.5">
              开启后仅管理员可访问网站前台，普通用户将看到维护提示
            </p>
          </div>
          <Switch
            checked={config.enableMaintenanceMode}
            onChange={(v) => update("enableMaintenanceMode", v)}
            className={`${
              config.enableMaintenanceMode ? "bg-warning" : "bg-base-300"
            } relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-warning/50 focus:ring-offset-2`}
          >
            <span
              className={`${
                config.enableMaintenanceMode ? "translate-x-6" : "translate-x-1"
              } inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform`}
            />
          </Switch>
        </div>
      </SectionCard>

      <div className="flex justify-end">
        <button
          onClick={onSave}
          disabled={isSaving}
          className="btn btn-primary gap-2"
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

// ── helpers ────────────────────────────────────────────────────────────────────

function SectionCard({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle: string;
  children: React.ReactNode;
}) {
  return (
    <div className="card bg-base-100 border border-base-300 shadow-sm">
      <div className="card-body gap-4">
        <div className="flex items-baseline gap-2 pb-2 border-b border-base-200">
          <h2 className="font-semibold text-base">{title}</h2>
          <span className="text-xs text-base-content/30 tracking-widest uppercase">
            {subtitle}
          </span>
        </div>
        {children}
      </div>
    </div>
  );
}

function Field({
  label,
  icon,
  children,
}: {
  label: string;
  icon?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="form-control gap-1">
      <label className="label py-0">
        <span className="label-text font-medium flex items-center gap-1.5">
          {icon && <span className="text-base-content/40">{icon}</span>}
          {label}
        </span>
      </label>
      {children}
    </div>
  );
}
