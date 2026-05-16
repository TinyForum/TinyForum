// ==================== 新增：零代码相关 hooks ====================

import { botApi } from "@/shared/api/modules/bot";
import { getErrorMessage } from "@/shared/lib/utils";
import { useState, useCallback, useEffect } from "react";
import toast from "react-hot-toast";
import { Flow, NocodeMetadata } from "../noco.type";

interface UseNocodeMetadataReturn {
  metadata: NocodeMetadata | null;
  loading: boolean;
  error: string | null;
  fetchMetadata: () => Promise<void>;
}

export function useNocodeMetadata(): UseNocodeMetadataReturn {
  const [metadata, setMetadata] = useState<NocodeMetadata | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchMetadata = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await botApi.nocode.getMetadata();
      console.log("机器人：", response.data.data);
      if (response.data.code === 0) {
        setMetadata(response.data.data ?? null);
      } else {
        throw new Error(response.data.message || "获取元数据失败");
      }
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, []);

  // 组件挂载时获取节点元数据
  useEffect(() => {
    fetchMetadata();
  }, [fetchMetadata]);

  return { metadata, loading, error, fetchMetadata };
}

interface UseValidateFlowReturn {
  validate: (flow: Flow) => Promise<{ valid: boolean; errors?: string[] }>;
  loading: boolean;
  error: string | null;
}

export function useValidateFlow(): UseValidateFlowReturn {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const validate = useCallback(async (flow: Flow) => {
    setLoading(true);
    setError(null);
    try {
      const response = await botApi.nocode.validateFlow({ flow });
      if (response.data.code === 0) {
        return response.data.data as { valid: boolean; errors?: string[] };
      } else {
        throw new Error(response.data.message || "流程校验失败");
      }
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
      return { valid: false, errors: [errorMsg] };
    } finally {
      setLoading(false);
    }
  }, []);

  return { validate, loading, error };
}
