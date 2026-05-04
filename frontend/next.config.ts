import { readFileSync } from "fs";
import yaml from "js-yaml";
import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";
import { RemotePattern } from "next/dist/shared/lib/image-config";

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

type OutputType = "standalone" | "export";

interface ProxyConfig {
  enabled_dev_only?: boolean;
  backend_url?: string;
  source?: string;
  destination_pattern?: string;
}

// YAML 原始配置结构（全部 snake_case）
interface RawConfigShape {
  output?: OutputType;
  images?: {
    dangerously_allow_svg?: boolean;
    remote_patterns?: Array<{
      protocol: string;
      hostname: string;
      port?: string;
      pathname?: string;
    }>;
  };
  allowed_dev_origins?: string[];
  experimental?: {
    proxy_timeout?: number;
  };
  proxy?: ProxyConfig;
}

// 辅助：将对象的所有键从 snake_case 转为 camelCase（递归）
function snakeToCamel(str: string): string {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
}

function convertKeys<T>(obj: T): T {
  if (Array.isArray(obj)) {
    return obj.map(convertKeys) as T;
  }
  if (obj !== null && typeof obj === "object") {
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
  return obj;
}

// 加载原始配置（snake_case）
function loadRawConfig(): RawConfigShape {
  const configPath = process.env.CONFIG_PATH || "./config.yaml";
  try {
    const fileContents = readFileSync(configPath, "utf8");
    const config = yaml.load(fileContents) as RawConfigShape | undefined;
    if (!config) {
      console.warn(`⚠️ 配置文件 ${configPath} 内容为空，将使用默认配置。`);
      return {};
    }
    return config;
  } catch (err) {
    const error = err as { code?: string; message?: string };
    if (error.code === "ENOENT") {
      console.warn(`⚠️ 配置文件 ${configPath} 不存在，将使用默认配置。`);
    } else {
      console.error(`❌ 读取配置文件 ${configPath} 失败: ${error.message}`);
    }
    return {};
  }
}

// 校验 output
function validateOutput(value: unknown): OutputType | undefined {
  if (value === "standalone" || value === "export") return value;
  if (value !== undefined) {
    console.warn(`⚠️ output 值 "${value}" 无效，回退至 "standalone"`);
  }
  return undefined;
}

// 加载并转换为驼峰
const rawSnake = loadRawConfig();
const validatedOutput = validateOutput(rawSnake.output);

// 转换整个配置（蛇形→驼峰）
const converted = convertKeys(rawSnake) as {
  output?: OutputType;
  images?: {
    dangerouslyAllowSVG?: boolean;
    remotePatterns?: Array<{
      protocol: string;
      hostname: string;
      port?: string;
      pathname?: string;
    }>;
  };
  allowedDevOrigins?: string[];
  experimental?: {
    proxyTimeout?: number;
  };
  proxy?: {
    enabledDevOnly?: boolean;
    backendUrl?: string;
    source?: string;
    destinationPattern?: string;
  };
};

// 确保 remotePatterns 符合 Next.js 的 RemotePattern 类型（添加可选 port 和 pathname）
const remotePatterns: RemotePattern[] | undefined =
  converted.images?.remotePatterns?.map((p) => ({
    protocol: p.protocol as "http" | "https",
    hostname: p.hostname,
    port: p.port,
    pathname: p.pathname,
  }));

// 基础 NextConfig（所有字段已是驼峰，符合 Next.js 要求）
const baseConfig: Partial<NextConfig> = {
  output: validatedOutput ?? "standalone",
  images: {
    dangerouslyAllowSVG: converted.images?.dangerouslyAllowSVG ?? true,
    remotePatterns: remotePatterns ?? [
      { protocol: "https", hostname: "api.dicebear.com" },
      { protocol: "https", hostname: "images.unsplash.com" },
      { protocol: "http", hostname: "localhost" },
    ],
  },
  distDir: "build",
  experimental: converted.experimental ?? {
    proxyTimeout: 10000,
  },
};

// 自定义配置：allowedDevOrigins 不放在 NextConfig 顶层
const customConfig = {
  allowedDevOrigins: converted.allowedDevOrigins ?? [
    "192.168.5.180",
    "localhost",
    "*.local-origin.dev",
  ],
};

// 辅助：检查后端
async function checkBackendReachable(url: string): Promise<void> {
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

// 最终 NextConfig
const nextConfig: NextConfig = {
  ...baseConfig,
  serverRuntimeConfig: customConfig,
  async rewrites() {
    const proxy = converted.proxy ?? {};
    const enabledDevOnly = proxy.enabledDevOnly !== false;

    if (enabledDevOnly && process.env.NODE_ENV === "production") {
      return [];
    }

    const backendUrl = process.env.BACKEND_URL || proxy.backendUrl;
    if (!backendUrl) {
      console.error(
        "❌ 代理目标未配置，请在环境变量 BACKEND_URL 或 config.yaml 的 proxy.backend_url 中设置。",
      );
      return [];
    }

    const cleanUrl = backendUrl.replace(/\/$/, "");
    const source = proxy.source ?? "/api/v1/:path*";
    const destination = `${cleanUrl}${proxy.destinationPattern ?? "/api/v1/:path*"}`;

    checkBackendReachable(cleanUrl).catch(() => {});
    console.log(`🔁 启用代理: ${source} → ${destination}`);
    return [{ source, destination }];
  },
};

export default withNextIntl(nextConfig);
