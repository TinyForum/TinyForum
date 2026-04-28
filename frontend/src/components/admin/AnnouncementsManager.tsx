// components/admin/AnnouncementsManager.tsx
"use client";

import { useState } from "react";
import { format, parseISO, isValid } from "date-fns";
import { zhCN } from "date-fns/locale";
import toast from "react-hot-toast";
import { Globe, BookOpen, Megaphone } from "lucide-react";
import type {
  Announcement,
  AnnouncementType,
  AnnouncementStatus,
} from "@/lib/api/modules/announcements";
import { useAdminAnnouncements } from "@/hooks/admin/useAdminAnnouncements";
import { useBoard } from "@/hooks/useBoard";
import { AnnouncementForm } from "./AnnouncementForm";
import { AnnouncementList } from "./AnnouncementList";

import { CreateAnnouncementPayload } from "@/lib/api/modules/announcements";
// 类型定义 - 让可选字段真正可选
interface AnnouncementFormValues {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  type: AnnouncementType;
  is_pinned: boolean;
  status: AnnouncementStatus;
  is_global: boolean;
  board_id?: number | null;
  published_at?: string | null;
  expired_at?: string | null;
}

// 日期格式化
const formatDate = (dateStr: string | null): string => {
  if (!dateStr) return "未发布";
  const date = parseISO(dateStr);
  if (!isValid(date)) return "无效日期";
  return format(date, "yyyy-MM-dd HH:mm", { locale: zhCN });
};

// 将 ISO 日期转换为 datetime-local 需要的格式
const toDateTimeLocal = (isoString: string | null | undefined): string => {
  if (!isoString) return "";
  try {
    const date = new Date(isoString);
    if (isNaN(date.getTime())) return "";
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    const hours = String(date.getHours()).padStart(2, "0");
    const minutes = String(date.getMinutes()).padStart(2, "0");
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  } catch {
    return "";
  }
};

// 将 datetime-local 值转换为 ISO 格式
const fromDateTimeLocal = (
  localValue: string | null | undefined,
): string | null => {
  if (!localValue) return null;
  try {
    const date = new Date(localValue);
    if (isNaN(date.getTime())) return null;
    return date.toISOString();
  } catch {
    return null;
  }
};

// 转换公告数据为表单数据
const announcementToFormValues = (
  announcement: Announcement,
): AnnouncementFormValues => {
  return {
    title: announcement.title,
    content: announcement.content,
    summary: announcement.summary || "",
    cover: announcement.cover || "",
    type: announcement.type,
    is_pinned: announcement.is_pinned,
    status: announcement.status,
    is_global: announcement.is_global,
    board_id: announcement.board_id || null,
    published_at: toDateTimeLocal(announcement.published_at),
    expired_at: toDateTimeLocal(announcement.expired_at),
  };
};

// 转换表单提交数据
const formValuesToPayload = (
  values: AnnouncementFormValues,
): CreateAnnouncementPayload => {
  return {
    title: values.title,
    content: values.content,
    summary: values.summary || "",
    cover: values.cover || "",
    type: values.type,
    is_pinned: values.is_pinned,
    is_global: values.is_global,
    board_id: values.is_global ? null : values.board_id || null,
    published_at: fromDateTimeLocal(values.published_at),
    expired_at: fromDateTimeLocal(values.expired_at),
  };
};

export function AnnouncementsManager({ t }: { t: (key: string) => string }) {
  const [announcementType, setAnnouncementType] = useState<"global" | "board">(
    "global",
  );
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [editingAnnouncement, setEditingAnnouncement] =
    useState<Announcement | null>(null);

  const {
    announcements,
    pinnedAnnouncements,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
    publishAnnouncement,
    pinAnnouncement,
    isLoading,
    isSubmitting,
  } = useAdminAnnouncements();

  // 获取板块列表
  const { boards, loading: boardsLoading } = useBoard({
    autoLoad: true,
    page: 1,
    pageSize: 100,
  });

  // 过滤公告
  const filteredAnnouncements = announcements.filter((ann: Announcement) => {
    if (announcementType === "global") {
      if (ann.is_global !== true) return false;
    } else {
      if (ann.is_global !== false) return false;
    }
    return !ann.is_pinned;
  });

  const filteredPinnedAnnouncements = pinnedAnnouncements.filter(
    (ann: Announcement) => {
      if (announcementType === "global") {
        return ann.is_global === true;
      } else {
        return ann.is_global === false;
      }
    },
  );

  // 打开创建表单
  const handleCreate = (): void => {
    setEditingAnnouncement(null);
    setModalVisible(true);
  };

  // 打开编辑表单
  const handleEdit = (announcement: Announcement): void => {
    setEditingAnnouncement(announcement);
    setModalVisible(true);
  };

  // 提交表单
  const handleFormSubmit = async (
    values: AnnouncementFormValues,
  ): Promise<void> => {
    const payload = formValuesToPayload(values);

    let result: Announcement | null;
    if (editingAnnouncement) {
      result = await updateAnnouncement(editingAnnouncement.id, payload);
    } else {
      result = await createAnnouncement(payload);
    }

    if (result) {
      setModalVisible(false);
      setEditingAnnouncement(null);
      toast.success(
        editingAnnouncement ? t("update_success") : t("create_success"),
      );
    }
  };

  // 删除公告
  const handleDelete = async (id: number): Promise<void> => {
    if (confirm(t("confirm_delete"))) {
      const success = await deleteAnnouncement(id);
      if (success) toast.success(t("delete_success"));
    }
  };

  // 置顶/取消置顶
  const handlePin = async (
    id: number,
    currentPinned: boolean,
  ): Promise<void> => {
    const success = await pinAnnouncement(id, !currentPinned);
    if (success)
      toast.success(!currentPinned ? t("pin_success") : t("unpin_success"));
  };

  // 发布公告
  const handlePublish = async (id: number): Promise<void> => {
    if (confirm(t("confirm_publish"))) {
      const success = await publishAnnouncement(id);
      if (success) toast.success(t("publish_success"));
    }
  };

  // 样式辅助函数
  const getTypeBadge = (type: AnnouncementType): string => {
    const styles: Record<AnnouncementType, string> = {
      normal: "badge-info",
      important: "badge-warning",
      emergency: "badge-error",
      event: "badge-success",
    };
    return styles[type] || "badge-info";
  };

  const getTypeText = (type: AnnouncementType): string => {
    const texts: Record<AnnouncementType, string> = {
      normal: t("normal"),
      important: t("important"),
      emergency: t("emergency"),
      event: t("event"),
    };
    return texts[type] || type;
  };

  // 判断公告是否过期
  const isExpired = (expiredAt: string | null): boolean => {
    if (!expiredAt) return false;
    return new Date(expiredAt) < new Date();
  };

  const getStatusBadge = (status: string, expiredAt: string | null): string => {
    if (status === "draft") return "badge-ghost";
    if (status === "published") {
      if (isExpired(expiredAt)) return "badge-error";
      return "badge-success";
    }
    return "badge-ghost";
  };

  const getStatusText = (status: string, expiredAt: string | null): string => {
    if (status === "draft") return t("draft");
    if (status === "published") {
      if (isExpired(expiredAt)) return t("expired");
      return t("published");
    }
    return status;
  };

  // 转换板块数据为表单需要的格式
  const boardOptions = boards.map((board) => ({
    id: board.id,
    name: board.name,
  }));

  return (
    <div className="space-y-6">
      {/* 头部操作栏 */}
      <div className="flex justify-between items-center flex-wrap gap-2">
        <div className="flex gap-2">
          <button
            onClick={() => setAnnouncementType("global")}
            className={`btn btn-sm ${announcementType === "global" ? "btn-primary" : "btn-ghost"}`}
          >
            <Globe className="w-4 h-4" /> {t("global_announcement")}
          </button>
          <button
            onClick={() => setAnnouncementType("board")}
            className={`btn btn-sm ${announcementType === "board" ? "btn-primary" : "btn-ghost"}`}
          >
            <BookOpen className="w-4 h-4" /> {t("board_announcement")}
          </button>
        </div>
        <button className="btn btn-primary btn-sm" onClick={handleCreate}>
          <Megaphone className="w-4 h-4" /> {t("create_announcement")}
        </button>
      </div>

      {/* 公告列表 */}
      <AnnouncementList
        announcements={filteredAnnouncements}
        pinnedAnnouncements={filteredPinnedAnnouncements}
        isLoading={isLoading}
        getTypeBadge={getTypeBadge}
        getTypeText={getTypeText}
        getStatusBadge={getStatusBadge}
        getStatusText={getStatusText}
        formatDate={formatDate}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onPin={handlePin}
        onPublish={handlePublish}
        t={t}
      />

      {/* 表单模态框 */}
      <AnnouncementForm
        isOpen={modalVisible}
        onClose={() => {
          setModalVisible(false);
          setEditingAnnouncement(null);
        }}
        onSubmit={handleFormSubmit}
        defaultValues={
          editingAnnouncement
            ? announcementToFormValues(editingAnnouncement)
            : undefined
        }
        isEditing={!!editingAnnouncement}
        isSubmitting={isSubmitting}
        boards={boardOptions}
        boardsLoading={boardsLoading}
        t={t}
      />
    </div>
  );
}

// 需要导入 CreateAnnouncementPayload 类型
