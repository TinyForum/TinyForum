import { readFileSync } from "fs";
import yaml from "js-yaml";
import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

type OutputType = "standalone" | "export";

interface ProxyConfig {
  enabledDevOnly?: boolean;
  backendUrl?: string;
  source?: string;
  destinationPattern?: string;
}

interface ConfigShape {
  output?: OutputType;
  images?: NextConfig["images"];
  allowedDevOrigins?: string[];
  experimental?: NextConfig["experimental"];
  proxy?: ProxyConfig;
}

// 1. 加载 YAML 配置，带详细错误提示
function loadConfig(): ConfigShape {
  const configPath = process.env.CONFIG_PATH || "./config.yaml";
  try {
    const fileContents = readFileSync(configPath, "utf8");
    const config = yaml.load(fileContents) as ConfigShape | undefined;
    if (!config) {
      console.warn(`⚠️ 配置文件 ${configPath} 内容为空，将使用默认配置。`);
      return {} as ConfigShape;
    }
    return config;
  } catch (err) {
    const error = err as { code?: string; message?: string };
    if (error.code === "ENOENT") {
      console.warn(`⚠️ 配置文件 ${configPath} 不存在，将使用默认配置。`);
    } else {
      console.error(`❌ 读取配置文件 ${configPath} 失败: ${error.message}`);
    }
    return {} as ConfigShape;
  }
}

// 2. 校验 output 字段，不合法时给出警告
function validateOutput(value: unknown): OutputType | undefined {
  if (value === "standalone" || value === "export") {
    return value;
  }
  if (value !== undefined) {
    console.warn(
      `⚠️ 配置中的 output 值 "${value}" 无效，已回退到默认值 "standalone"。`,
    );
  }
  return undefined;
}

const rawConfig = loadConfig();
const validatedOutput = validateOutput(rawConfig.output);

// 3. 基础配置（带默认值）
const baseConfig: Partial<NextConfig> = {
  output: validatedOutput ?? "standalone",
  images: rawConfig.images ?? {
    dangerouslyAllowSVG: true,
    remotePatterns: [
      { protocol: "https", hostname: "api.dicebear.com" },
      { protocol: "https", hostname: "images.unsplash.com" },
      { protocol: "http", hostname: "localhost" },
    ],
  },
  allowedDevOrigins: rawConfig.allowedDevOrigins ?? [
    "192.168.5.180",
    "localhost",
    "*.local-origin.dev",
  ],
  experimental: rawConfig.experimental ?? {
    proxyTimeout: 10000,
  },
};

// 辅助函数：检查后端服务是否可访问（仅用于启动时提示，不阻塞启动）
async function checkBackendReachable(url: string): Promise<void> {
  if (!url) return;
  try {
    // 使用 fetch 检测根路径，超时设为 2 秒
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 2000);
    await fetch(url, { method: "HEAD", signal: controller.signal });
    clearTimeout(timeoutId);
    console.log(`✅ 后端服务检测正常: ${url}`);
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : String(err);
    console.warn(
      `⚠️ 无法连接到后端服务 ${url}，请确保后端已启动。\n` +
        `   错误信息: ${errorMessage}`,
    );
  }
}

const nextConfig: NextConfig = {
  ...baseConfig,

  async rewrites() {
    const proxyConfig = rawConfig.proxy ?? {};
    const enabledDevOnly = proxyConfig.enabledDevOnly !== false;

    // 生产环境下若配置为仅开发启用，则返回空代理
    if (enabledDevOnly && process.env.NODE_ENV === "production") {
      return [];
    }

    // 获取后端地址：环境变量 > YAML 配置
    const backendUrl = process.env.BACKEND_URL || proxyConfig.backendUrl;
    if (!backendUrl) {
      console.error(
        "❌ 代理目标未配置！请设置环境变量 BACKEND_URL 或在 config.yaml 中配置 proxy.backendUrl。\n" +
          "   示例: BACKEND_URL=http://localhost:8080 npm run dev",
      );
      return [];
    }

    // 去除末尾多余的斜杠
    const cleanBackendUrl = backendUrl.replace(/\/$/, "");
    const source = proxyConfig.source ?? "/api/v1/:path*";
    const destination = `${cleanBackendUrl}${proxyConfig.destinationPattern ?? "/api/v1/:path*"}`;

    // 异步检测后端可达性（不阻塞 rewrites 返回，仅提示）
    checkBackendReachable(cleanBackendUrl).catch(() => {});

    console.log(`🔁 启用代理: ${source} → ${destination}`);
    return [{ source, destination }];
  },
};

export default withNextIntl(nextConfig);
