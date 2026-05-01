// components/admin/AnnouncementList.tsx

import {
  AnnouncementDO,
  AnnouncementType,
  AnnouncementStatus,
} from "@/shared/api/types/announcement.model";
import { Loader2, Pin } from "lucide-react";
import { AnnouncementCard } from "./AnnouncementCard";

interface AnnouncementListProps {
  announcements: AnnouncementDO[];
  pinnedAnnouncements: AnnouncementDO[];
  isLoading: boolean;
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
