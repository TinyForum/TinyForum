// components/admin/AnnouncementsManager.tsx
"use client";

import { useState } from "react";
import toast from "react-hot-toast";
import { Globe, BookOpen, Megaphone } from "lucide-react";

import { AnnouncementForm } from "./AnnouncementForm";
import { AnnouncementList } from "./AnnouncementList";

import { useAdminAnnouncements } from "../hooks/useAdminAnnouncements";
import { useBoard } from "@/features/boards/hooks/useBoard";
import {
  getAnnouncementStatusBadge,
  getAnnouncementStatusText,
  getAnnouncementTypeBadge,
  getAnnouncementTypeText,
} from "@/shared/lib/utils/announcement";
import { AnnouncementFormValues } from "@/shared/type/announcement.type";
import { formValuesToPayload } from "@/shared/lib/helpers/formValuesToPayload";
import { formatDate } from "@/shared/lib/format/formatDate";
import { announcementToFormValues } from "@/shared/lib/helpers/announcementToFormValues";
import { AnnouncementDO } from "@/shared/api/types/announcement.model";

export function AnnouncementsManager({ t }: { t: (key: string) => string }) {
  const [announcementType, setAnnouncementType] = useState<"global" | "board">(
    "global",
  );
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [editingAnnouncement, setEditingAnnouncement] =
    useState<AnnouncementDO | null>(null);

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
  const filteredAnnouncements = announcements.filter((ann: AnnouncementDO) => {
    if (announcementType === "global") {
      if (ann.is_global !== true) return false;
    } else {
      if (ann.is_global !== false) return false;
    }
    return !ann.is_pinned;
  });

  const filteredPinnedAnnouncements = pinnedAnnouncements.filter(
    (ann: AnnouncementDO) => {
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
  const handleEdit = (announcement: AnnouncementDO): void => {
    setEditingAnnouncement(announcement);
    setModalVisible(true);
  };

  // 提交表单
  const handleFormSubmit = async (
    values: AnnouncementFormValues,
  ): Promise<void> => {
    const payload = formValuesToPayload(values);
    console.log(payload);

    let result: AnnouncementDO | null;
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
        getTypeBadge={getAnnouncementTypeBadge}
        getTypeText={getAnnouncementTypeText}
        getStatusBadge={getAnnouncementStatusBadge}
        getStatusText={getAnnouncementStatusText}
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
