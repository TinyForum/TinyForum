"use client";

import { useState, useCallback, useEffect, useMemo } from "react";
import toast from "react-hot-toast";
import { configApi } from "@/shared/api/modules/config";
import type { ConfigReloadResult } from "@/shared/api/modules/config";

export type ConfigFieldType = "string" | "number" | "boolean" | "duration" | "select" | "list";

export interface ConfigField {
  key: string;
  value: unknown;
  type: ConfigFieldType;
  label: string;
  group: string;
  description?: string;
  options?: string[];
  sensitive: boolean;
}

const SENSITIVE_KEY_PATTERNS = [/password/i, /secret/i, /token/i, /key/i, /dsn/i];

function isSensitive(key: string): boolean {
  return SENSITIVE_KEY_PATTERNS.some((p) => p.test(key));
}

function inferType(key: string, value: unknown): ConfigFieldType {
  if (typeof value === "boolean") return "boolean";
  if (typeof value === "number") return "number";
  if (Array.isArray(value)) return "list";
  const s = String(value);
  if (/^\d+(s|ms|m|h|d)$/.test(s) || key.endsWith("_timeout") || key.endsWith("_every"))
    return "duration";
  if (["debug", "release"].includes(s) || key === "server.mode" || key === "log.level")
    return "select";
  return "string";
}

function inferOptions(key: string): string[] | undefined {
  if (key === "server.mode") return ["debug", "release", "test"];
  if (key === "log.level") return ["debug", "info", "warn", "error"];
  return undefined;
}

function inferGroup(key: string): string {
  const dot = key.indexOf(".");
  return dot > 0 ? key.substring(0, dot) : "general";
}

const KEY_LABELS_ZH: Record<string, string> = {
  "server.host": "监听地址",
  "server.port": "监听端口",
  "server.mode": "运行模式",
  "server.read_timeout": "读超时",
  "server.write_timeout": "写超时",
  "server.max_header_bytes": "最大请求头",
  "log.level": "日志级别",
  "log.filename": "日志文件",
  "log.max_size": "日志最大大小 (MB)",
  "log.max_age": "日志保留天数",
  "log.max_backups": "最大备份数",
  "log.compress": "压缩日志",
  "log.console": "输出到控制台",
  "log.json_format": "JSON 格式",
  "api.host": "API 主机",
  "api.port": "API 端口",
  "api.prefix": "API 前缀",
  "api.version": "API 版本",
  "api.protocol": "API 协议",
  "attachment.max_size": "附件最大大小",
  "attachment.allowed_ext": "允许的扩展名",
  "attachment.upload_dir": "上传目录",
  "attachment.url_prefix": "URL 前缀",
  "frontend.host": "前端主机",
  "frontend.port": "前端端口",
  "frontend.protocol": "前端协议",
  "plugins.storage_dir": "插件存储目录",
  "version": "应用版本",
};

function fieldLabel(key: string): string {
  return KEY_LABELS_ZH[key] ?? key.replace(/\./g, " ").replace(/_/g, " ");
}

export function useSiteConfig() {
  const [fields, setFields] = useState<ConfigField[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isReloading, setIsReloading] = useState(false);
  const [configFile, setConfigFile] = useState("basic");

  const loadConfig = useCallback(async (file: string) => {
    setIsLoading(true);
    try {
      const res = await configApi.getKV(file);
      const kv = res.data.data?.kv;
      if (kv) {
        const parsed: ConfigField[] = Object.entries(kv)
          .filter(([key]) => key !== "")
          .map(([key, value]) => ({
            key,
            value,
            type: inferType(key, value),
            label: fieldLabel(key),
            group: inferGroup(key),
            options: inferOptions(key),
            sensitive: isSensitive(key),
          }));
        setFields(parsed);
      }
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })
        ?.response?.data?.message || "加载配置失败";
      toast.error(msg);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    loadConfig(configFile);
  }, [configFile, loadConfig]);

  const grouped = useMemo(() => {
    const groups: Record<string, ConfigField[]> = {};
    for (const f of fields) {
      (groups[f.group] ??= []).push(f);
    }
    return groups;
  }, [fields]);

  const groupOrder = useMemo(() => {
    const order = ["server", "api", "frontend", "log", "attachment", "plugins", "allow_origins", "version", "general"];
    return order.filter((g) => grouped[g]?.length);
  }, [grouped]);

  const updateField = useCallback((key: string, value: unknown) => {
    setFields((prev) => prev.map((f) => (f.key === key ? { ...f, value } : f)));
  }, []);

  const save = useCallback(async () => {
    setIsSaving(true);
    try {
      const kvPayload: Record<string, string> = {};
      for (const f of fields) {
        if (f.sensitive) continue;
        kvPayload[f.key] = String(f.value);
      }
      await configApi.updateKV(configFile, kvPayload);
      toast.success("配置保存成功，已自动重载");
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })
        ?.response?.data?.message || "保存失败";
      toast.error(msg);
    } finally {
      setIsSaving(false);
    }
  }, [fields, configFile]);

  const reloadConfig = useCallback(async (): Promise<ConfigReloadResult | null> => {
    setIsReloading(true);
    try {
      const res = await configApi.reload();
      const result = res.data.data;
      if (result) {
        toast.success(`配置已重载 (${result.reloadedFiles.length} 个文件)`);
      }
      return result || null;
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })
        ?.response?.data?.message || "重载失败";
      toast.error(msg);
      return null;
    } finally {
      setIsReloading(false);
    }
  }, []);

  const switchConfig = useCallback((file: string) => {
    setConfigFile(file);
  }, []);

  return {
    fields,
    grouped,
    groupOrder,
    configFile,
    isLoading,
    isSaving,
    isReloading,
    updateField,
    save,
    reloadConfig,
    switchConfig,
    loadConfig,
  };
}
