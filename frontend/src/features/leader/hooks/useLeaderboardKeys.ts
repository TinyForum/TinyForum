export const leaderboardKeys = {
  all: ["leaderboard"] as const,
  list: (params?: { limit?: number; fields?: string }) =>
    [...leaderboardKeys.all, "list", params?.limit, params?.fields] as const,
};
