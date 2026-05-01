import {
  AnnouncementStatus,
  AnnouncementType,
} from "@/shared/api/types/announcement.model";

/**
 * 判断公告是否过期
 * @param expiredAt
 * @returns
 */
export const isAnnouncementExpired = (expiredAt: string | null): boolean => {
  if (!expiredAt) return false;
  return new Date(expiredAt) < new Date();
};

/**
 * 获取状态样式类名
 * @param status
 * @param expiredAt
 * @returns
 */
export const getAnnouncementStatusBadge = (
  status: AnnouncementStatus,
  expiredAt: string | null,
): string => {
  if (status === AnnouncementStatus.Draft) return "badge-ghost";
  if (status === AnnouncementStatus.Published) {
    return isAnnouncementExpired(expiredAt) ? "badge-error" : "badge-success";
  }
  return "badge-ghost";
};

/**
 * 获取状态文本
 * @param status
 * @param expiredAt
 * @param t
 * @returns
 */
export const getAnnouncementStatusText = (
  status: AnnouncementStatus,
  expiredAt: string | null,
  t: (key: string) => string,
): string => {
  if (status === AnnouncementStatus.Draft) return t("draft");
  if (status === AnnouncementStatus.Published) {
    return isAnnouncementExpired(expiredAt) ? t("expired") : t("published");
  }
  if (status === AnnouncementStatus.Archived) return t("archived");
  return t("unknown");
};

/**
 * 获取类型样式类名
 * @param type
 * @returns
 */
export const getAnnouncementTypeBadge = (type: AnnouncementType): string => {
  const map: Record<AnnouncementType, string> = {
    [AnnouncementType.Normal]: "badge-info",
    [AnnouncementType.Important]: "badge-warning",
    [AnnouncementType.Emergency]: "badge-error",
    [AnnouncementType.Event]: "badge-success",
  };
  return map[type] || "badge-ghost";
};

/**
 * 获取类型文本
 * @param type
 * @param t
 * @returns
 */
export const getAnnouncementTypeText = (
  type: AnnouncementType,
  t: (key: string) => string,
): string => {
  const map: Record<AnnouncementType, string> = {
    [AnnouncementType.Normal]: t("normal"),
    [AnnouncementType.Important]: t("important"),
    [AnnouncementType.Emergency]: t("emergency"),
    [AnnouncementType.Event]: t("event"),
  };
  return map[type];
};

/**
 * 获取公告状态的颜色样式
 * 适用于 Ant Design 的 Tag 或其他 UI 库的 color 属性
 */
export function getAnnouncementStatusColor(status: AnnouncementStatus): string {
  const colorMap: Record<AnnouncementStatus, string> = {
    [AnnouncementStatus.All]: "default", // 新增
    [AnnouncementStatus.Draft]: "gray",
    [AnnouncementStatus.Published]: "green",
    [AnnouncementStatus.Archived]: "orange",
  };
  return colorMap[status] || "default";
}

/**
 * 格式化公告发布时间
 */
export function formatAnnouncementTime(dateStr: string | null): string {
  if (!dateStr) return "未发布";
  const date = new Date(dateStr);
  return date.toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}
