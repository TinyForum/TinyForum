// ==================== 新增：零代码相关 hooks ====================

import { useQuery, useMutation } from "@tanstack/react-query";
import { botApi } from "@/shared/api/modules/bot";
import { getErrorMessage } from "@/shared/lib/utils";
import { useState } from "react";
import toast from "react-hot-toast";
import { NocodeMetadata } from "../noco.type";
import { Flow } from "@/shared/api/types/bot.model";
import { botKeys } from "./useBotKeys";

interface UseNocodeMetadataReturn {
  metadata: NocodeMetadata | null;
  loading: boolean;
  error: string | null;
  fetchMetadata: () => Promise<void>;
}

export function useNocodeMetadata(): UseNocodeMetadataReturn {
  // 查询：零代码节点元数据
  const query = useQuery({
    queryKey: botKeys.metadata(),
    queryFn: async () => {
      const res = await botApi.nocode.getMetadata();
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取元数据失败");
      }
      if (!res.data.data) {
        throw new Error("节点元数据为空");
      }
      return res.data.data;
    },
  });

  const metadata = query.data ?? null;
  const loading = query.isLoading;
  const error = query.error ? getErrorMessage(query.error) : null;

  const fetchMetadata = async (): Promise<void> => {
    await query.refetch();
  };

  return { metadata, loading, error, fetchMetadata };
}

interface UseValidateFlowReturn {
  validate: (flow: Flow) => Promise<{ valid: boolean; errors?: string[] }>;
  loading: boolean;
  error: string | null;
}

export function useValidateFlow(): UseValidateFlowReturn {
  const [error, setError] = useState<string | null>(null);

  // 变更：校验流程
  const validateMutation = useMutation({
    mutationFn: async (
      flow: Flow,
    ): Promise<{ valid: boolean; errors?: string[] }> => {
      const res = await botApi.nocode.validateFlow(flow);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "流程校验失败");
      }
      if (!res.data.data) {
        throw new Error("流程校验返回数据为空");
      }
      return res.data.data;
    },
    onError: (err: Error) => {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    },
  });

  const validate = async (
    flow: Flow,
  ): Promise<{ valid: boolean; errors?: string[] }> => {
    setError(null);
    try {
      return await validateMutation.mutateAsync(flow);
    } catch (err) {
      const errorMsg = getErrorMessage(err);
      return { valid: false, errors: [errorMsg] };
    }
  };

  const loading = validateMutation.isPending;

  return { validate, loading, error };
}
