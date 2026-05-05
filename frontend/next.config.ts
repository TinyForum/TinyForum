// next.config.ts
import { ConvertedConfig } from "@/shared/lib/yaml/config.type";
import {
  loadConvertedConfig,
  getRemotePatterns,
  getAllowedDevOrigins,
  checkBackendReachable,
} from "@/shared/lib/yaml/loadConfig";
import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

// 加载配置（不再在顶层立即执行副作用，但为了 nextConfig 仍需要这里加载）
const convertedConfig: ConvertedConfig = loadConvertedConfig();
const remotePatterns = getRemotePatterns(convertedConfig);
const allowedDevOrigins = getAllowedDevOrigins(convertedConfig);

const baseConfig: Partial<NextConfig> = {
  output: convertedConfig.output ?? "standalone",
  images: {
    dangerouslyAllowSVG: convertedConfig.images?.dangerouslyAllowSVG ?? true,
    remotePatterns,
  },
  experimental: convertedConfig.experimental ?? {
    proxyTimeout: 10000,
  },
};

const customConfig = {
  allowedDevOrigins,
};

const nextConfig: NextConfig = {
  ...baseConfig,
  serverRuntimeConfig: customConfig,
  async rewrites() {
    const proxy = convertedConfig.proxy ?? {};
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
