// hooks/useTopicList.ts
import { useCallback } from "react";
import { useQuery } from "@tanstack/react-query";
import { toast } from "react-hot-toast";
import { topicApi, CreateTopicPayload } from "@/shared/api/modules/topics";
import { PageData } from "@/shared/api/types/basic.model";
import { Topic } from "@/shared/api/types/topic.model";

// 话题列表查询 Hook（分页）
export function useTopicList(page: number, pageSize: number) {
  const query = useQuery<PageData<Topic>>({
    queryKey: ["topics", "list", page],
    queryFn: async () => {
      const response = await topicApi.list({ page, page_size: pageSize });
      if (response.data.code !== 0) {
        toast.error(response.data.message || "加载失败");
        throw new Error(response.data.message || "加载失败");
      }
      if (!response.data.data) {
        throw new Error("话题数据为空");
      }
      return response.data.data;
    },
  });

  return {
    data: query.data,
    isLoading: query.isLoading,
    refetch: query.refetch,
  };
}

// 创建话题 Hook
export function useCreateTopic() {
  const createTopic = useCallback(
    (data: CreateTopicPayload) => topicApi.create(data),
    [],
  );
  return { createTopic };
}
