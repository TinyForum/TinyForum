// /hooks/admin/useScoreData.ts
import { scoreApi } from "@/lib/api/modules/score";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export function useScoreData(userId?: string | undefined) {
  const queryClient = useQueryClient();

  // 获取所有用户积分列表
  const { data: scoreRecords, isLoading: isLoadingRecords } = useQuery({
    queryKey: ["score", "list", userId],
    queryFn: () =>
      scoreApi.getAllUserScore({ id: userId ? Number(userId) : undefined }),
    enabled: true,
  });

  // 获取当前用户自己的积分
  const { data: myScore, isLoading: isLoadingMyScore } = useQuery({
    queryKey: ["score", "me"],
    queryFn: () => scoreApi.getUserScore(),
  });

  // 设置用户积分
  const setScoreMutation = useMutation({
    mutationFn: ({
      userId,
      score,
      reason,
    }: {
      userId: number;
      score: number;
      reason: string;
    }) => scoreApi.setUserScore(userId, score, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["score"] });
    },
  });

  // 增加用户积分
  const addScoreMutation = useMutation({
    mutationFn: ({
      userId,
      increment,
      reason,
    }: {
      userId: number;
      increment: number;
      reason: string;
    }) => scoreApi.addUserScore(userId, increment, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["score"] });
    },
  });

  // 扣除用户积分
  const subtractScoreMutation = useMutation({
    mutationFn: ({
      userId,
      decrement,
      reason,
    }: {
      userId: number;
      decrement: number;
      reason: string;
    }) => scoreApi.subtractUserScore(userId, decrement, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["score"] });
    },
  });

  return {
    // 查询数据
    scoreRecords: scoreRecords?.data,
    myScore: myScore?.data,
    isLoadingRecords,
    isLoadingMyScore,

    // 同步方法
    setScore: setScoreMutation.mutate,
    addScore: addScoreMutation.mutate,
    subtractScore: subtractScoreMutation.mutate,

    // 异步方法（返回 Promise）
    setScoreAsync: setScoreMutation.mutateAsync,
    addScoreAsync: addScoreMutation.mutateAsync,
    subtractScoreAsync: subtractScoreMutation.mutateAsync,

    // 状态
    isSettingScore: setScoreMutation.isPending,
    isAddingScore: addScoreMutation.isPending,
    isSubtractingScore: subtractScoreMutation.isPending,

    // 错误
    setScoreError: setScoreMutation.error,
    addScoreError: addScoreMutation.error,
    subtractScoreError: subtractScoreMutation.error,
  };
}
