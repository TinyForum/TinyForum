"use client";

import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";
import { Modal } from "./Modal";
import Image from "next/image";
import { AnnouncementFormValues } from "@/shared/type/announcement.type";
import {
  AnnouncementType,
  AnnouncementStatus,
} from "@/shared/api/types/announcement.model";

// 表单验证 Schema - 直接使用数字枚举
const announcementSchema = z
  .object({
    title: z.string().min(1, "请输入标题").max(200, "标题不能超过200字"),
    content: z.string().min(1, "请输入内容"),
    summary: z.string().max(500, "摘要不能超过500字").optional(),
    cover: z.string().url("请输入有效的URL").optional().or(z.literal("")),
    type: z.preprocess((val) => {
      if (typeof val === "string") return parseInt(val, 10);
      if (typeof val === "number") return val;
      return undefined;
    }, z.nativeEnum(AnnouncementType)),
    is_pinned: z.boolean().default(false),
    status: z.preprocess(
      (val) => {
        if (typeof val === "string") {
          if (val === "draft") return AnnouncementStatus.Draft;
          if (val === "published") return AnnouncementStatus.Published;
          if (val === "archived") return AnnouncementStatus.Archived;
          return parseInt(val, 10);
        }
        if (typeof val === "number") return val;
        return undefined;
      },
      z
        .nativeEnum(AnnouncementStatus)
        .refine((v) => v !== AnnouncementStatus.All, {
          message: "无效的状态",
        }),
    ),
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
    control,
    setValue,
    formState: { errors },
  } = useForm<AnnouncementFormValues>({
    resolver: zodResolver(announcementSchema),
    defaultValues: {
      type: AnnouncementType.Normal,
      status: AnnouncementStatus.Draft,
      is_pinned: false,
      is_global: true,
      cover: "",
      board_id: null,
      published_at: null,
      expired_at: null,
    },
  });

  const isGlobal = useWatch({ control, name: "is_global" });
  const selectedBoardId = useWatch({ control, name: "board_id" });
  const coverUrl = useWatch({ control, name: "cover" });

  // 当弹窗打开或 defaultValues 变化时，重置表单
  useEffect(() => {
    if (isOpen) {
      reset({
        type: AnnouncementType.Normal,
        status: AnnouncementStatus.Draft,
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
    console.log(values);
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
          {coverUrl && (
            <div className="mt-2">
              <Image
                src={coverUrl}
                alt="封面预览"
                width={128}
                height={128}
                className="object-cover rounded-lg border"
                onError={(e) => {
                  (e.target as HTMLImageElement).style.display = "none";
                }}
              />
            </div>
          )}
        </div>

        {/* 类型和状态 */}
        <div className="grid grid-cols-2 gap-4">
          {/* 类型 - 使用数字枚举值 */}
          <div>
            <label className="label text-sm font-medium">{t("type")}</label>
            <select
              className="select select-bordered w-full"
              {...register("type")}
            >
              <option value={AnnouncementType.Normal}>{t("normal")}</option>
              <option value={AnnouncementType.Important}>
                {t("important")}
              </option>
              <option value={AnnouncementType.Emergency}>
                {t("emergency")}
              </option>
              <option value={AnnouncementType.Event}>{t("event")}</option>
            </select>
            {errors.type && (
              <p className="text-error text-xs mt-1">{errors.type.message}</p>
            )}
          </div>

          {/* 状态 - 使用数字枚举值，排除 All */}
          <div>
            <label className="label text-sm font-medium">{t("status")}</label>
            <select
              className="select select-bordered w-full"
              {...register("status")}
            >
              <option value={AnnouncementStatus.Draft}>{t("draft")}</option>
              <option value={AnnouncementStatus.Published}>
                {t("published")}
              </option>
              <option value={AnnouncementStatus.Archived}>
                {t("archived")}
              </option>
            </select>
            {errors.status && (
              <p className="text-error text-xs mt-1">{errors.status.message}</p>
            )}
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
              <p className="text-xs text-base-content/50">
                {t("global_announcement_desc")}
              </p>
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
              <p className="text-xs text-base-content/50">
                {t("board_announcement_desc")}
              </p>
            </div>
          </label>

          {/* 板块选择（仅非全局时显示） */}
          {!isGlobal && (
            <div className="ml-8 mt-2">
              <label className="label text-sm font-medium">
                {t("select_board")}
              </label>
              <select
                className="select select-bordered w-full"
                value={selectedBoardId ?? ""}
                onChange={handleBoardChange}
                disabled={boardsLoading}
              >
                <option value="">
                  {boardsLoading
                    ? t("loading_boards")
                    : t("please_select_board")}
                </option>
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
                <p className="text-error text-xs mt-1">
                  {errors.board_id.message}
                </p>
              )}
            </div>
          )}
        </div>

        {/* 置顶开关 */}
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            className="toggle toggle-sm"
            {...register("is_pinned")}
          />
          <span className="text-sm">{t("pin")}</span>
        </div>

        {/* 时间 */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="label text-sm font-medium">
              {t("publish_time")}
            </label>
            <input
              type="datetime-local"
              className="input input-bordered w-full"
              {...register("published_at")}
            />
            <p className="text-xs text-base-content/50 mt-1">
              {t("publish_time_hint")}
            </p>
          </div>
          <div>
            <label className="label text-sm font-medium">
              {t("expire_time")}
            </label>
            <input
              type="datetime-local"
              className="input input-bordered w-full"
              {...register("expired_at")}
            />
            <p className="text-xs text-base-content/50 mt-1">
              {t("expire_time_hint")}
            </p>
            {errors.expired_at && (
              <p className="text-error text-xs mt-1">
                {errors.expired_at.message}
              </p>
            )}
          </div>
        </div>

        {/* 按钮 */}
        <div className="flex justify-end gap-2 pt-4">
          <button type="button" className="btn btn-ghost" onClick={onClose}>
            {t("cancel")}
          </button>
          <button
            type="submit"
            className="btn btn-primary"
            disabled={isSubmitting}
          >
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
