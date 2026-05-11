// ---------- 类型定义 ----------
export type OutputType = "standalone" | "export";

export interface ProxyConfig {
  enabled_dev_only?: boolean;
  backend_url?: string;
  source?: string;
  destination_pattern?: string;
}

export interface RawConfigShape {
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
  jwt?: {
    secret: string;
    expire: string;
    refresh_expire: string;
    issuer: string;
  };
}

export interface ConvertedConfig {
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
  jwt?: {
    secret: string;
    expire: string;
    refresh_expire: string;
    issuer: string;
  };
  proxy?: {
    enabledDevOnly?: boolean;
    backendUrl?: string;
    source?: string;
    destinationPattern?: string;
  };
}
