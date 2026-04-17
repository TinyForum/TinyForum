// components/admin/AnnouncementForm.tsx
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";
import { Modal } from "./Modal";

// 表单验证 Schema
const announcementSchema = z
  .object({
    title: z.string().min(1, "请输入标题").max(200, "标题不能超过200字"),
    content: z.string().min(1, "请输入内容"),
    summary: z.string().max(500, "摘要不能超过500字").optional(),
    cover: z.string().url("请输入有效的URL").optional().or(z.literal("")),
    type: z.enum(["normal", "important", "emergency", "event"]),
    is_pinned: z.boolean().default(false),
    status: z.enum(["draft", "published", "archived"]).default("draft"),
    is_global: z.boolean().default(true),
    board_id: z.number().nullable().optional(),
    published_at: z.string().nullable().optional(),
    expired_at: z.string().nullable().optional(),
  })
  .refine(
    (data) => {
      if (data.published_at && data.expired_at) {
        return new Date(data.published_at) < new Date(data.expired_at);
      }
      return true;
    },
    {
      message: "过期时间必须晚于发布时间",
      path: ["expired_at"],
    },
  );

type AnnouncementFormValues = z.infer<typeof announcementSchema>;

interface AnnouncementFormProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (values: AnnouncementFormValues) => Promise<void>;
  defaultValues?: Partial<AnnouncementFormValues>;
  isEditing?: boolean;
  isSubmitting?: boolean;
  boards?: Array<{ id: number; name: string }>;
  boardsLoading?: boolean;
  t: (key: string) => string;
}

export function AnnouncementForm({
  isOpen,
  onClose,
  onSubmit,
  defaultValues,
  isEditing = false,
  isSubmitting = false,
  boards = [],
  boardsLoading = false,
  t,
}: AnnouncementFormProps) {
  const {
    register,
    handleSubmit,
    reset,
    watch,
    setValue,
    formState: { errors },
  } = useForm<AnnouncementFormValues>({
    resolver: zodResolver(announcementSchema),
    defaultValues: {
      type: "normal",
      status: "draft",
      is_pinned: false,
      is_global: true,
      cover: "",
      board_id: null,
      published_at: null,
      expired_at: null,
    },
  });

  const isGlobal = watch("is_global");
  const selectedBoardId = watch("board_id");

  // 当弹窗打开或 defaultValues 变化时，重置表单
  useEffect(() => {
    if (isOpen) {
      reset({
        type: "normal",
        status: "draft",
        is_pinned: false,
        is_global: true,
        cover: "",
        board_id: null,
        published_at: null,
        expired_at: null,
        ...defaultValues,
      });
    }
  }, [isOpen, defaultValues, reset]);

  // 当 is_global 变化时，清空 board_id
  const handleGlobalChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const checked = e.target.checked;
    setValue("is_global", checked);
    if (checked) {
      setValue("board_id", null);
    }
  };

  // 选择板块时，自动设置为非全局
  const handleBoardChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    if (value) {
      setValue("board_id", parseInt(value));
      setValue("is_global", false);
    } else {
      setValue("board_id", null);
    }
  };

  const handleFormSubmit = async (values: AnnouncementFormValues) => {
    await onSubmit(values);
    // 提交成功后关闭弹窗，表单重置会在下次打开时由 useEffect 处理
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? t("edit_announcement") : t("create_announcement")}
    >
      <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
        {/* 标题 */}
        <div>
          <label className="label text-sm font-medium">
            {t("title")} <span className="text-error">*</span>
          </label>
          <input
            type="text"
            className="input input-bordered w-full"
            placeholder={t("title_placeholder")}
            {...register("title")}
          />
          {errors.title && (
            <p className="text-error text-xs mt-1">{errors.title.message}</p>
          )}
        </div>

        {/* 内容 */}
        <div>
          <label className="label text-sm font-medium">
            {t("content")} <span className="text-error">*</span>
          </label>
          <textarea
            className="textarea textarea-bordered w-full h-32"
            placeholder={t("content_placeholder")}
            {...register("content")}
          />
          {errors.content && (
            <p className="text-error text-xs mt-1">{errors.content.message}</p>
          )}
        </div>

        {/* 摘要 */}
        <div>
          <label className="label text-sm font-medium">{t("summary")}</label>
          <textarea
            className="textarea textarea-bordered w-full h-20"
            placeholder={t("summary_placeholder")}
            {...register("summary")}
          />
          {errors.summary && (
            <p className="text-error text-xs mt-1">{errors.summary.message}</p>
          )}
        </div>

        {/* 封面图片 */}
        <div>
          <label className="label text-sm font-medium">{t("cover")}</label>
          <input
            type="text"
            className="input input-bordered w-full"
            placeholder={t("cover_placeholder")}
            {...register("cover")}
          />
          {errors.cover && (
            <p className="text-error text-xs mt-1">{errors.cover.message}</p>
          )}
          {watch("cover") && (
            <div className="mt-2">
              <img
                src={watch("cover")}
                alt="封面预览"
                className="w-32 h-32 object-cover rounded-lg border"
                onError={(e) => {
                  (e.target as HTMLImageElement).style.display = "none";
                }}
              />
            </div>
          )}
        </div>

        {/* 类型和状态 */}
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

        {/* 全局/板块切换 */}
        <div className="space-y-3">
          <label className="label text-sm font-medium">{t("scope")}</label>
          
          {/* 全局公告选项 */}
          <label className="flex items-center gap-3 p-3 rounded-lg border cursor-pointer hover:bg-base-200">
            <input
              type="radio"
              className="radio radio-primary"
              checked={isGlobal === true}
              onChange={() => {
                setValue("is_global", true);
                setValue("board_id", null);
              }}
            />
            <div className="flex-1">
              <span className="font-medium">{t("global_announcement")}</span>
              <p className="text-xs text-base-content/50">{t("global_announcement_desc")}</p>
            </div>
          </label>

          {/* 板块公告选项 */}
          <label className="flex items-center gap-3 p-3 rounded-lg border cursor-pointer hover:bg-base-200">
            <input
              type="radio"
              className="radio radio-primary"
              checked={isGlobal === false}
              onChange={() => setValue("is_global", false)}
            />
            <div className="flex-1">
              <span className="font-medium">{t("board_announcement")}</span>
              <p className="text-xs text-base-content/50">{t("board_announcement_desc")}</p>
            </div>
          </label>

          {/* 板块选择（仅非全局时显示） */}
          {!isGlobal && (
            <div className="ml-8 mt-2">
              <label className="label text-sm font-medium">{t("select_board")}</label>
              <select
                className="select select-bordered w-full"
                value={selectedBoardId || ""}
                onChange={handleBoardChange}
                disabled={boardsLoading}
              >
                <option value="">{boardsLoading ? t("loading_boards") : t("please_select_board")}</option>
                {boards.map((board) => (
                  <option key={board.id} value={board.id}>
                    {board.name}
                  </option>
                ))}
              </select>
              {boardsLoading && (
                <p className="text-xs text-base-content/50 mt-1 flex items-center gap-1">
                  <Loader2 className="w-3 h-3 animate-spin" />
                  {t("loading_boards")}
                </p>
              )}
              {errors.board_id && (
                <p className="text-error text-xs mt-1">{errors.board_id.message}</p>
              )}
            </div>
          )}
        </div>

        {/* 置顶开关 */}
        <div className="flex items-center gap-2">
          <input type="checkbox" className="toggle toggle-sm" {...register("is_pinned")} />
          <span className="text-sm">{t("pin")}</span>
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
            <p className="text-xs text-base-content/50 mt-1">{t("publish_time_hint")}</p>
          </div>
          <div>
            <label className="label text-sm font-medium">{t("expire_time")}</label>
            <input
              type="datetime-local"
              className="input input-bordered w-full"
              {...register("expired_at")}
            />
            <p className="text-xs text-base-content/50 mt-1">{t("expire_time_hint")}</p>
            {errors.expired_at && (
              <p className="text-error text-xs mt-1">{errors.expired_at.message}</p>
            )}
          </div>
        </div>

        {/* 按钮 */}
        <div className="flex justify-end gap-2 pt-4">
          <button type="button" className="btn btn-ghost" onClick={onClose}>
            {t("cancel")}
          </button>
          <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
            {isSubmitting ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : isEditing ? (
              t("update")
            ) : (
              t("create")
            )}
          </button>
        </div>
      </form>
    </Modal>
  );
}