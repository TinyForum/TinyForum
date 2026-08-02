export const timelineKeys = {
  all: ["timeline-events"] as const,
  following: (page: number, pageSize: number) =>
    [...timelineKeys.all, "following", page, pageSize] as const,
};
