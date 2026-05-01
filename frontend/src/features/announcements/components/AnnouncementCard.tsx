import { AnnouncementDO } from "@/shared/api/types/announcement.model";
import {
  getAnnouncementTypeBadge,
  getAnnouncementStatusBadge,
  getAnnouncementTypeText,
  getAnnouncementStatusText,
} from "@/shared/lib/utils/announcement";
import {
  MegaphoneIcon,
  PinIcon,
  CalendarIcon,
  EyeIcon,
  FileTextIcon,
} from "lucide-react";
import { useTranslations } from "next-intl";
import Link from "next/link";

export function AnnouncementCard({
  announcement,
}: {
  announcement: AnnouncementDO;
}) {
  const t = useTranslations("Announcement"); // 获取翻译函数

  // 使用工具函数获取类型样式和文本
  const typeBadgeClass = getAnnouncementTypeBadge(announcement.type);
  const typeLabel = getAnnouncementTypeText(announcement.type, t);

  // 获取状态样式和文本（根据过期时间判断是否显示“已过期”）
  const statusBadgeClass = getAnnouncementStatusBadge(
    announcement.status,
    announcement.expired_at,
  );
  const statusText = getAnnouncementStatusText(
    announcement.status,
    announcement.expired_at,
    t,
  );

  // 格式化时间
  const formatDate = (dateStr: string | null) => {
    if (!dateStr) return t("pending_publish");
    return new Date(dateStr).toLocaleDateString("zh-CN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });
  };

  return (
    <Link href={`/announcements/${announcement.id}`}>
      <div className="group bg-base-100 rounded-xl shadow-sm hover:shadow-md transition-all duration-300 p-5 border border-base-200 hover:border-primary/20 cursor-pointer">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-3 flex-wrap">
              <MegaphoneIcon className="w-4 h-4 text-primary" />
              {/* 类型标签 */}
              <span
                className={`text-xs px-2.5 py-0.5 rounded-full font-medium ${typeBadgeClass}`}
              >
                {typeLabel}
              </span>
              {/* 置顶标识 */}
              {announcement.is_pinned && (
                <span className="text-xs px-2.5 py-0.5 rounded-full bg-warning/10 text-warning flex items-center gap-1 font-medium">
                  <PinIcon className="w-3 h-3" />
                  {t("pinned")}
                </span>
              )}
              {/* 状态标签（草稿、已发布、已归档、已过期） */}
              <span
                className={`text-xs px-2.5 py-0.5 rounded-full font-medium ${statusBadgeClass}`}
              >
                {statusText}
              </span>
            </div>
            <h3 className="font-semibold text-base-content text-lg mb-2 line-clamp-1 group-hover:text-primary transition-colors">
              {announcement.title}
            </h3>
            <p className="text-sm text-base-content/60 line-clamp-2 mb-4">
              {announcement.summary ||
                announcement.content?.replace(/<[^>]*>/g, "").slice(0, 150)}
            </p>
            <div className="flex items-center gap-4 text-xs text-base-content/40">
              <div className="flex items-center gap-1.5">
                <CalendarIcon className="w-3.5 h-3.5" />
                {formatDate(
                  announcement.published_at || announcement.created_at,
                )}
              </div>
              <div className="flex items-center gap-1.5">
                <EyeIcon className="w-3.5 h-3.5" />
                {announcement.view_count || 0} {t("views")}
              </div>
              {announcement.board && (
                <div className="flex items-center gap-1.5">
                  <FileTextIcon className="w-3.5 h-3.5" />
                  {announcement.board.name}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}
