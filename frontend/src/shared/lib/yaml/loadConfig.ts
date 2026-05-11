// lib/config-loader.ts
import { readFileSync } from "fs";
import yaml from "js-yaml";
import type { RemotePattern } from "next/dist/shared/lib/image-config";
import { ConvertedConfig, OutputType, RawConfigShape } from "./config.type";

// ---------- 工具函数 ----------
export function snakeToCamel(str: string): string {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
}

export function convertKeys<T>(obj: T): T {
  if (Array.isArray(obj)) {
    return obj.map(convertKeys) as T;
  }
  if (obj !== null && typeof obj === "object") {
    const newObj: { [key: string]: unknown } = {};
    for (const [key, value] of Object.entries(obj)) {
      const camelKey = snakeToCamel(key);
      newObj[camelKey] = convertKeys(value);
    }
    return newObj as T;
  }
  return obj;
}

export function validateOutput(value: unknown): OutputType | undefined {
  if (value === "standalone" || value === "export") return value;
  if (value !== undefined) {
    console.warn(`⚠️ output 值 "${value}" 无效，回退至 "standalone"`);
  }
  return undefined;
}

// 加载原始 YAML 配置
export function loadRawConfig(configPath?: string): RawConfigShape {
  const path = configPath || process.env.CONFIG_PATH || "./config.yml";
  try {
    const fileContents = readFileSync(path, "utf8");
    const config = yaml.load(fileContents) as RawConfigShape | undefined;
    if (!config) {
      console.warn(`⚠️ 配置文件 ${path} 内容为空，将使用默认配置。`);
      return {};
    }
    return config;
  } catch (err) {
    const error = err as { code?: string; message?: string };
    if (error.code === "ENOENT") {
      console.warn(`⚠️ 配置文件 ${path} 不存在，将使用默认配置。`);
    } else {
      console.error(`❌ 读取配置文件 ${path} 失败: ${error.message}`);
    }
    return {};
  }
}

// 加载并转换配置（返回完全驼峰化的配置对象）
export function loadConvertedConfig(configPath?: string): ConvertedConfig {
  const raw = loadRawConfig(configPath);
  const output = validateOutput(raw.output);
  const converted = convertKeys(raw) as ConvertedConfig;
  if (output !== undefined) {
    converted.output = output;
  }
  return converted;
}

// 获取经过验证的 remotePatterns（符合 Next.js RemotePattern 类型）
export function getRemotePatterns(converted: ConvertedConfig): RemotePattern[] {
  return (
    converted.images?.remotePatterns?.map((p) => ({
      protocol: p.protocol as "http" | "https",
      hostname: p.hostname,
      port: p.port,
      pathname: p.pathname,
    })) ?? [
      { protocol: "https", hostname: "api.dicebear.com" },
      { protocol: "https", hostname: "images.unsplash.com" },
      { protocol: "http", hostname: "localhost" },
    ]
  );
}

// 获取 allowedDevOrigins
export function getAllowedDevOrigins(converted: ConvertedConfig): string[] {
  return (
    converted.allowedDevOrigins ?? [
      "192.168.5.180",
      "localhost",
      "*.local-origin.dev",
    ]
  );
}

// 检查后端是否可访问（异步）
export async function checkBackendReachable(url: string): Promise<void> {
  if (!url) return;
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 2000);
    await fetch(url, { method: "HEAD", signal: controller.signal });
    clearTimeout(timeoutId);
    console.log(`✅ 后端服务检测正常: ${url}`);
  } catch (err) {
    console.warn(
      `⚠️ 无法连接到后端服务 ${url}，请确保后端已启动。\nERROR: ${err}`,
    );
  }
}
