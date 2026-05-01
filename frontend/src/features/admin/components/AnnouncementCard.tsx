import {
  AnnouncementDO,
  AnnouncementType,
  AnnouncementStatus,
} from "@/shared/api/types/announcement.model";
import {
  Pin,
  Eye,
  Calendar,
  Edit,
  Trash2,
  Clock,
  CheckCircle,
} from "lucide-react";

// ==================== 类型定义 ====================
interface AnnouncementCardProps {
  announcement: AnnouncementDO;
  getTypeBadge: (type: AnnouncementType, t: (key: string) => string) => string;
  getTypeText: (type: AnnouncementType, t: (key: string) => string) => string;
  getStatusBadge: (
    status: AnnouncementStatus,
    expiredAt: string | null,
  ) => string;
  getStatusText: (
    status: AnnouncementStatus,
    expiredAt: string | null,
    t: (key: string) => string,
  ) => string;
  formatDate: (date: string | null) => string;
  onEdit: (ann: AnnouncementDO) => void;
  onDelete: (id: number) => void;
  onPin: (id: number, currentPinned: boolean) => void;
  onPublish: (id: number) => void;
  t: (key: string) => string;
}

// ==================== 辅助函数 ====================
const isExpired = (expiredAt: string | null): boolean => {
  if (!expiredAt) return false;
  return new Date(expiredAt) < new Date();
};

const isDraft = (status: AnnouncementStatus): boolean => {
  return status === AnnouncementStatus.Draft;
};

// ==================== 公告卡片组件 ====================
export function AnnouncementCard({
  announcement,
  getTypeBadge,
  getTypeText,
  getStatusBadge,
  getStatusText,
  formatDate,
  onEdit,
  onDelete,
  onPin,
  onPublish,
  t,
}: AnnouncementCardProps) {
  const expired = isExpired(announcement.expired_at);
  const draft = isDraft(announcement.status);

  // 状态判断
  const isDimmed = draft || expired; // 草稿或过期变灰
  const hasLineThrough = expired; // 过期添加删除线
  const showPublishButton = draft; // 草稿显示发布按钮

  // 卡片样式
  const cardClasses = `card bg-base-100 border ${isDimmed ? "opacity-60" : ""} ${
    expired ? "border-error/30" : "border-base-300"
  }`;

  return (
    <div className={cardClasses}>
      <div className="card-body p-4">
        <div className="flex justify-between items-start">
          {/* 左侧内容 */}
          <div className="flex-1">
            {/* 标签行 */}
            <div className="flex items-center gap-2 mb-2 flex-wrap">
              {/* 类型标签 */}
              <span
                className={`badge badge-sm ${getTypeBadge(announcement.type, t)}`}
              >
                {getTypeText(announcement.type, t)}
              </span>

              {/* 置顶标签 */}
              {announcement.is_pinned && (
                <span className="badge badge-sm badge-warning">
                  <Pin className="w-3 h-3 mr-1" />
                  {t("pinned")}
                </span>
              )}

              {/* 状态标签 - 传入 expired_at 用于判断过期显示 */}
              <span
                className={`badge badge-sm ${getStatusBadge(announcement.status, announcement.expired_at)}`}
              >
                {getStatusText(announcement.status, announcement.expired_at, t)}
              </span>

              {/* 浏览量 */}
              <span className="text-xs text-base-content/50 flex items-center gap-1">
                <Eye className="w-3 h-3" />
                {announcement.view_count}
              </span>
            </div>

            {/* 标题 */}
            <h3
              className={`font-semibold mb-1 ${hasLineThrough ? "line-through" : ""}`}
            >
              {announcement.title}
            </h3>

            {/* 内容摘要 */}
            <p
              className={`text-sm text-base-content/70 line-clamp-2 ${hasLineThrough ? "line-through" : ""}`}
            >
              {announcement.content}
            </p>

            {/* 时间信息 */}
            <div className="flex items-center gap-4 mt-2 text-xs text-base-content/50">
              <span className="flex items-center gap-1">
                <Calendar className="w-3 h-3" />
                {t("publish_time")}: {formatDate(announcement.published_at)}
              </span>
              {announcement.expired_at && (
                <span className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  {t("expires_at")}: {formatDate(announcement.expired_at)}
                </span>
              )}
            </div>
          </div>

          {/* 右侧操作按钮 */}
          <div className="flex gap-1 ml-4">
            {/* 置顶/取消置顶 */}
            <button
              className="btn btn-ghost btn-xs"
              onClick={() => onPin(announcement.id, announcement.is_pinned)}
              title={announcement.is_pinned ? t("unpin") : t("pin")}
            >
              <Pin
                className={`w-3 h-3 ${announcement.is_pinned ? "text-warning" : ""}`}
              />
            </button>

            {/* 发布按钮（仅草稿显示） */}
            {showPublishButton && (
              <button
                className="btn btn-ghost btn-xs text-success"
                onClick={() => onPublish(announcement.id)}
                title={t("publish")}
              >
                <CheckCircle className="w-3 h-3 mr-1" />
                {t("publish")}
              </button>
            )}

            {/* 编辑按钮 */}
            <button
              className="btn btn-ghost btn-xs"
              onClick={() => onEdit(announcement)}
              title={t("edit")}
            >
              <Edit className="w-3 h-3" />
            </button>

            {/* 删除按钮 */}
            <button
              className="btn btn-ghost btn-xs text-error"
              onClick={() => onDelete(announcement.id)}
              title={t("delete")}
            >
              <Trash2 className="w-3 h-3" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
