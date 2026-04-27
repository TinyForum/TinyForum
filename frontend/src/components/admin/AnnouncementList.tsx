// components/admin/AnnouncementList.tsx
import {
  Announcement,
  AnnouncementType,
} from "@/lib/api/modules/announcements";
import { AnnouncementCard } from "./AnnouncementCard";
import { Loader2, Pin } from "lucide-react";

interface AnnouncementListProps {
  announcements: Announcement[];
  pinnedAnnouncements: Announcement[];
  isLoading: boolean;
  getTypeBadge: (type: AnnouncementType) => string;
  getTypeText: (type: AnnouncementType) => string;
  getStatusBadge: (status: string, expiredAt: string | null) => string; // 修改：增加第二个参数
  getStatusText: (status: string, expiredAt: string | null) => string; // 修改：增加第二个参数
  formatDate: (date: string | null) => string;
  onEdit: (ann: Announcement) => void;
  onDelete: (id: number) => void;
  onPin: (id: number, currentPinned: boolean) => void;
  onPublish: (id: number) => void;
  t: (key: string) => string;
}

export function AnnouncementList({
  announcements,
  pinnedAnnouncements,
  isLoading,
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
}: AnnouncementListProps) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 置顶公告区域 */}
      {pinnedAnnouncements.length > 0 && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium text-base-content/70">
            <Pin className="w-4 h-4" />
            <span>{t("pinned_announcements")}</span>
          </div>
          <div className="space-y-2">
            {pinnedAnnouncements.map((ann) => (
              <AnnouncementCard
                key={ann.id}
                announcement={ann}
                getTypeBadge={getTypeBadge}
                getTypeText={getTypeText}
                getStatusBadge={getStatusBadge}
                getStatusText={getStatusText}
                formatDate={formatDate}
                onEdit={onEdit}
                onDelete={onDelete}
                onPin={onPin}
                onPublish={onPublish}
                t={t}
              />
            ))}
          </div>
        </div>
      )}

      {/* 普通公告列表 */}
      <div className="space-y-3">
        {announcements.length === 0 ? (
          <div className="text-center py-8 text-base-content/50">
            {t("no_announcements")}
          </div>
        ) : (
          announcements.map((ann) => (
            <AnnouncementCard
              key={ann.id}
              announcement={ann}
              getTypeBadge={getTypeBadge}
              getTypeText={getTypeText}
              getStatusBadge={getStatusBadge}
              getStatusText={getStatusText}
              formatDate={formatDate}
              onEdit={onEdit}
              onDelete={onDelete}
              onPin={onPin}
              onPublish={onPublish}
              t={t}
            />
          ))
        )}
      </div>
    </div>
  );
}
