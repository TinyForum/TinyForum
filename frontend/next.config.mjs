/** @type {import('next').NextConfig} */
import createNextIntlPlugin from "next-intl/plugin";
const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

const nextConfig = {
  output: "standalone",
  images: {
    dangerouslyAllowSVG: true,
    remotePatterns: [
      { protocol: "https", hostname: "api.dicebear.com" },
      { protocol: "https", hostname: "images.unsplash.com" },
      { protocol: "http", hostname: "localhost" },
    ],
  },
  allowedDevOrigins: ["192.168.5.180", "localhost", "*.local-origin.dev"],

  async rewrites() {
    // 开发环境启用代理
    if (process.env.NODE_ENV === "production") return [];

    const backendUrl = process.env.BACKEND_URL;

    return [
      {
        source: "/api/v1/:path*",
        destination: `${backendUrl}/api/v1/:path*`,
      },
    ];
  },

  // 实验性配置
  experimental: {
    proxyTimeout: 10000,
  },
};

export default withNextIntl(nextConfig);
