"use client";

import { useAnnouncementsData } from "@/hooks/admin/useAnnouncementsData";
import { Globe, BookOpen, Megaphone, Edit, Trash2, Pin, Eye, Calendar, X, Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { format, formatISO, isValid, parse, parseISO } from "date-fns";
import { zhCN } from "date-fns/locale";
import toast from "react-hot-toast";
import type { Announcement, AnnouncementType } from "@/lib/api/modules/announcements";

// ==================== 表单验证 Schema ====================
const announcementSchema = z.object({
  title: z.string().min(1, "请输入标题").max(200, "标题不能超过200字"),
  content: z.string().min(1, "请输入内容"),
  summary: z.string().max(500, "摘要不能超过500字").optional(),
  type: z.enum(["normal", "important", "emergency", "event"]),
  is_pinned: z.boolean().default(false),
  status: z.enum(["draft", "published", "expired"]).default("draft"),
  is_global: z.boolean().default(true),
  published_at: z.string().nullable().optional(),
  expired_at: z.string().nullable().optional(),
}).refine((data) => {
  if (data.published_at && data.expired_at) {
    return new Date(data.published_at) < new Date(data.expired_at);
  }
  return true;
}, {
  message: "过期时间必须晚于发布时间",
  path: ["expired_at"],
});

type AnnouncementFormValues = z.infer<typeof announcementSchema>;

// ==================== 日期格式化 ====================
const formatDate = (dateStr: string | null): string => {
  if (!dateStr) return "未发布";
  const date = parseISO(dateStr);
  if (!isValid(date)) return "无效日期";
  return format(date, "yyyy-MM-dd HH:mm", { locale: zhCN });
};

// ==================== 模态框组件 ====================
function Modal({
  isOpen,
  onClose,
  title,
  children,
}: {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
}) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={onClose} />
      <div className="relative bg-base-100 rounded-lg shadow-xl w-full max-w-lg max-h-[90vh] overflow-auto">
        <div className="sticky top-0 bg-base-100 border-b border-base-300 px-6 py-4 flex justify-between items-center">
          <h3 className="text-lg font-semibold">{title}</h3>
          <button onClick={onClose} className="btn btn-sm btn-ghost btn-square">
            <X className="w-4 h-4" />
          </button>
        </div>
        <div className="p-6">{children}</div>
      </div>
    </div>
  );
}

// ==================== 公告管理组件 ====================
export function AnnouncementsManager({ t }: { t: (key: string) => string }) {
  const [announcementType, setAnnouncementType] = useState<"global" | "board">("global");
  const [modalVisible, setModalVisible] = useState(false);
  const [editingAnnouncement, setEditingAnnouncement] = useState<Announcement | null>(null);

  const {
    announcements,
    pinnedAnnouncements,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
    publishAnnouncement,
    pinAnnouncement,
    isLoading,
    // isSubmitting,
  } = useAnnouncementsData(true);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors, isSubmitting: isFormSubmitting },
  } = useForm<AnnouncementFormValues>({
    resolver: zodResolver(announcementSchema),
    defaultValues: {
      type: "normal",
      status: "draft",
      is_pinned: false,
      is_global: announcementType === "global",
      published_at: null,
      expired_at: null,
    },
  });

  const isGlobal = watch("is_global");

  // 根据类型过滤公告
  const filteredAnnouncements = announcements.filter((ann) => {
    if (announcementType === "global") {
      return ann.is_global === true;
    } else {
      return ann.is_global === false;
    }
  });

  const filteredPinnedAnnouncements = pinnedAnnouncements.filter((ann) => {
    if (announcementType === "global") {
      return ann.is_global === true;
    } else {
      return ann.is_global === false;
    }
  });

  // 打开创建模态框
  const handleCreate = () => {
    setEditingAnnouncement(null);
    reset({
      title: "",
      content: "",
      summary: "",
      type: "normal",
      is_pinned: false,
      is_global: announcementType === "global",
      published_at: null,
      expired_at: null,
    });
    setModalVisible(true);
  };

  // 打开编辑模态框
  const handleEdit = (announcement: Announcement) => {
    setEditingAnnouncement(announcement);
    reset({
      title: announcement.title,
      content: announcement.content,
      summary: announcement.summary || "",
      type: announcement.type,
      is_pinned: announcement.is_pinned,
      is_global: announcement.is_global,
      published_at: announcement.published_at,
      expired_at: announcement.expired_at,
    });
    setModalVisible(true);
  };

const formatDateTimeLocal = (value: string | null | undefined): string | null => {
  if (!value) return null;
  const date = parse(value, "yyyy-MM-dd'T'HH:mm", new Date());
  return formatISO(date);
};

  // 提交表单
  const onSubmit = async (values: AnnouncementFormValues) => {
    const payload = {
      ...values,
      published_at: formatDateTimeLocal(values.published_at),
      expired_at: formatDateTimeLocal(values.expired_at),
    };

    let result;
    if (editingAnnouncement) {
      result = await updateAnnouncement(editingAnnouncement.id, payload);
    } else {
      result = await createAnnouncement(payload);
    }

    if (result) {
      setModalVisible(false);
      reset();
      toast.success(editingAnnouncement ? t("update_success") : t("create_success"));
    }
  };

  // 删除公告
  const handleDelete = async (id: number) => {
    if (confirm(t("confirm_delete"))) {
      const success = await deleteAnnouncement(id);
      if (success) {
        toast.success(t("delete_success"));
      }
    }
  };

  // 置顶/取消置顶
  const handlePin = async (id: number, currentPinned: boolean) => {
    const success = await pinAnnouncement(id, !currentPinned);
    if (success) {
      toast.success(!currentPinned ? t("pin_success") : t("unpin_success"));
    }
  };

  // 发布公告
  const handlePublish = async (id: number) => {
    if (confirm(t("confirm_publish"))) {
      const success = await publishAnnouncement(id);
      if (success) {
        toast.success(t("publish_success"));
      }
    }
  };

  // 获取类型标签样式
  const getTypeBadge = (type: AnnouncementType) => {
    const styles = {
      normal: "badge-info",
      important: "badge-warning",
      emergency: "badge-error",
      event: "badge-success",
    };
    return styles[type] || "badge-info";
  };

  const getTypeText = (type: AnnouncementType) => {
    const texts = {
      normal: t("normal"),
      important: t("important"),
      emergency: t("emergency"),
      event: t("event"),
    };
    return texts[type] || type;
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      draft: "badge-ghost",
      published: "badge-success",
      archived: "badge-error",
    };
    return styles[status] || "badge-ghost";
  };

  const getStatusText = (status: string) => {
    const texts: Record<string, string> = {
      draft: t("draft"),
      published: t("published"),
      archived: t("archived"),
    };
    return texts[status] || status;
  };

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 头部操作栏 */}
      <div className="flex justify-between items-center">
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

      {/* 置顶公告区域 */}
      {filteredPinnedAnnouncements.length > 0 && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium text-base-content/70">
            <Pin className="w-4 h-4" />
            <span>{t("pinned_announcements")}</span>
          </div>
          <div className="space-y-2">
            {filteredPinnedAnnouncements.map((ann) => (
              <div key={ann.id} className="card bg-primary/5 border border-primary/20">
                <div className="card-body p-4">
                  <AnnouncementCard
                    announcement={ann}
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
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 公告列表 */}
      <div className="space-y-3">
        {filteredAnnouncements.length === 0 ? (
          <div className="text-center py-8 text-base-content/50">{t("no_announcements")}</div>
        ) : (
          filteredAnnouncements.map((ann) => (
            <div key={ann.id} className="card bg-base-100 border border-base-300">
              <div className="card-body p-4">
                <AnnouncementCard
                  announcement={ann}
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
              </div>
            </div>
          ))
        )}
      </div>

      {/* 创建/编辑模态框 */}
      <Modal
        isOpen={modalVisible}
        onClose={() => setModalVisible(false)}
        title={editingAnnouncement ? t("edit_announcement") : t("create_announcement")}
      >
       <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
  {/* 标题 */}
  <div>
    <label className="label text-sm font-medium">{t("title")}</label>
    <input
      type="text"
      className="input input-bordered w-full"
      placeholder={t("title_placeholder")}
      {...register("title")}
    />
    {errors.title && <p className="text-error text-xs mt-1">{errors.title.message}</p>}
  </div>

  {/* 内容 */}
  <div>
    <label className="label text-sm font-medium">{t("content")}</label>
    <textarea
      className="textarea textarea-bordered w-full h-32"
      placeholder={t("content_placeholder")}
      {...register("content")}
    />
    {errors.content && <p className="text-error text-xs mt-1">{errors.content.message}</p>}
  </div>

  {/* 摘要 */}
  <div>
    <label className="label text-sm font-medium">{t("summary")}</label>
    <textarea
      className="textarea textarea-bordered w-full h-20"
      placeholder={t("summary_placeholder")}
      {...register("summary")}
    />
    {errors.summary && <p className="text-error text-xs mt-1">{errors.summary.message}</p>}
  </div>

  {/* 第一行：类型和状态 */}
  <div className="grid grid-cols-2 gap-4">
    <div>
      <label className="label text-sm font-medium">{t("type")}</label>
      <select className="select select-bordered w-full" {...register("type")}>
        <option value="normal">{t("normal")}</option>
        <option value="important">{t("important")}</option>
        <option value="emergency">{t("emergency")}</option>
        <option value="event">{t("event")}</option>
      </select>
    </div>
    <div>
      <label className="label text-sm font-medium">{t("status")}</label>
      <select className="select select-bordered w-full" {...register("status")}>
        <option value="draft">{t("draft")}</option>
        <option value="published">{t("published")}</option>
        <option value="expired">{t("archived")}</option>
      </select>
    </div>
  </div>

  {/* 开关 */}
  <div className="flex gap-4">
    <label className="flex items-center gap-2">
      <input type="checkbox" className="toggle toggle-sm" {...register("is_pinned")} />
      <span className="text-sm">{t("pin")}</span>
    </label>
    <label className="flex items-center gap-2">
      <input
        type="checkbox"
        className="toggle toggle-sm"
        {...register("is_global")}
        disabled={announcementType === "global"}
      />
      <span className="text-sm">{t("global")}</span>
    </label>
  </div>

  {/* 时间 */}
  <div className="grid grid-cols-2 gap-4">
    <div>
      <label className="label text-sm font-medium">{t("publish_time")}</label>
      <input
        type="datetime-local"
        className="input input-bordered w-full"
        {...register("published_at")}
      />
    </div>
    <div>
      <label className="label text-sm font-medium">{t("expire_time")}</label>
      <input
        type="datetime-local"
        className="input input-bordered w-full"
        {...register("expired_at")}
      />
      {errors.expired_at && (
        <p className="text-error text-xs mt-1">{errors.expired_at.message}</p>
      )}
    </div>
  </div>

  {/* 按钮 */}
  <div className="flex justify-end gap-2 pt-4">
    <button type="button" className="btn btn-ghost" onClick={() => setModalVisible(false)}>
      {t("cancel")}
    </button>
    <button type="submit" className="btn btn-primary" disabled={isFormSubmitting}>
      {isFormSubmitting ? (
        <Loader2 className="w-4 h-4 animate-spin" />
      ) : editingAnnouncement ? (
        t("update")
      ) : (
        t("create")
      )}
    </button>
  </div>
</form>
      </Modal>
    </div>
  );
}

// ==================== 公告卡片组件 ====================
function AnnouncementCard({
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
}: {
  announcement: Announcement;
  getTypeBadge: (type: AnnouncementType) => string;
  getTypeText: (type: AnnouncementType) => string;
  getStatusBadge: (status: string) => string;
  getStatusText: (status: string) => string;
  formatDate: (date: string | null) => string;
  onEdit: (ann: Announcement) => void;
  onDelete: (id: number) => void;
  onPin: (id: number, currentPinned: boolean) => void;
  onPublish: (id: number) => void;
  t: (key: string) => string;
}) {
  return (
    <div className="flex justify-between items-start">
      <div className="flex-1">
        <div className="flex items-center gap-2 mb-2 flex-wrap">
          <span className={`badge badge-sm ${getTypeBadge(announcement.type)}`}>
            {getTypeText(announcement.type)}
          </span>
          {announcement.is_pinned && (
            <span className="badge badge-sm badge-warning">
              <Pin className="w-3 h-3 mr-1" />
              {t("pinned")}
            </span>
          )}
          <span className={`badge badge-sm ${getStatusBadge(announcement.status)}`}>
            {getStatusText(announcement.status)}
          </span>
          <span className="text-xs text-base-content/50 flex items-center gap-1">
            <Eye className="w-3 h-3" />
            {announcement.view_count}
          </span>
        </div>
        <h3 className="font-semibold mb-1">{announcement.title}</h3>
        <p className="text-sm text-base-content/70 line-clamp-2">{announcement.content}</p>
        <div className="flex items-center gap-4 mt-2 text-xs text-base-content/50">
          <span className="flex items-center gap-1">
            <Calendar className="w-3 h-3" />
            {formatDate(announcement.published_at)}
          </span>
          {announcement.expired_at && (
            <span className="flex items-center gap-1">
              <X className="w-3 h-3" />
              {t("expires")}: {formatDate(announcement.expired_at)}
            </span>
          )}
        </div>
      </div>
      <div className="flex gap-1">
        <button
          className="btn btn-ghost btn-xs"
          onClick={() => onPin(announcement.id, announcement.is_pinned)}
          title={announcement.is_pinned ? t("unpin") : t("pin")}
        >
          <Pin className={`w-3 h-3 ${announcement.is_pinned ? "text-warning" : ""}`} />
        </button>
        {announcement.status !== "published" && (
          <button
            className="btn btn-ghost btn-xs text-success"
            onClick={() => onPublish(announcement.id)}
            title={t("publish")}
          >
            {t("publish")}
          </button>
        )}
        <button className="btn btn-ghost btn-xs" onClick={() => onEdit(announcement)} title={t("edit")}>
          <Edit className="w-3 h-3" />
        </button>
        <button
          className="btn btn-ghost btn-xs text-error"
          onClick={() => onDelete(announcement.id)}
          title={t("delete")}
        >
          <Trash2 className="w-3 h-3" />
        </button>
      </div>
    </div>
  );
}