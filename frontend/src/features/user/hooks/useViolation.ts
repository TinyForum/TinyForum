import { useState, useCallback } from "react";
import {
  userViolationApi,
  ViolationVO,
} from "@/shared/api/modules/user/violation";

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

interface UseViolationReturn {
  violations: ViolationVO[]; // ✅ 始终为数组，不是 null
  loadViolations: () => Promise<void>;
  isLoading: boolean;
  error: string | null;
  fetchViolationDetail: (id: string) => Promise<ViolationVO | null>; // ✅ 可能返回 null
  submitAppeal: (id: string, reason: string) => Promise<boolean>;
  isAppealing: boolean;
  appealError: string | null;
}

export function useUserViolation(): UseViolationReturn {
  const [violations, setViolations] = useState<ViolationVO[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isAppealing, setIsAppealing] = useState(false);
  const [appealError, setAppealError] = useState<string | null>(null);

  const loadViolations = useCallback(async () => {
    console.log("loadViolations start");
    setIsLoading(true);
    setError(null);
    try {
      const response = await userViolationApi.listUserViolations();
      if (response.status === 200 && response.data.code === 0) {
        setViolations(response.data.data ?? []);
      } else {
        throw new Error(response.data.message || "获取违规列表失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      setError(
        errorObj.response?.data?.message ||
          errorObj.message ||
          "获取违规列表失败",
      );
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchViolationDetail = useCallback(
    async (id: string): Promise<ViolationVO | null> => {
      try {
        const response = await userViolationApi.getViolationDetail(id);
        if (response.status === 200 && response.data.code === 0) {
          return response.data.data || null;
        } else {
          throw new Error(response.data.message || "获取详情失败");
        }
      } catch (err: unknown) {
        const errorObj = err as ErrorResponse;
        const errorMsg =
          errorObj.response?.data?.message ||
          errorObj.message ||
          "获取详情失败";
        setError(errorMsg);
        return null; // ✅ 明确返回 null 表示失败
      }
    },
    [],
  );

  const submitAppeal = useCallback(
    async (id: string, reason: string): Promise<boolean> => {
      setIsAppealing(true);
      setAppealError(null);
      try {
        const response = await userViolationApi.appeal(id, reason);
        if (response.status === 200 && response.data.code === 0) {
          await loadViolations();
          return true;
        } else {
          throw new Error(response.data.message || "申诉提交失败");
        }
      } catch (err: unknown) {
        const errorObj = err as ErrorResponse;
        setAppealError(
          errorObj.response?.data?.message ||
            errorObj.message ||
            "申诉提交失败",
        );
        return false;
      } finally {
        setIsAppealing(false);
      }
    },
    [loadViolations],
  );

  return {
    violations,
    loadViolations,
    isLoading,
    error,
    fetchViolationDetail,
    submitAppeal,
    isAppealing,
    appealError,
  };
}
