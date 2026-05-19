"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import { useRouter } from "next/navigation";
import {
  useForm,
  Controller,
  useWatch,
  UseFormRegister,
  FieldErrors,
} from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import { FileText, Send, X, FolderOpen, Eye, Edit } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { Board } from "@/shared/api/types/board.model";
import { ApiResponse } from "@/shared/api/types/basic.model";
import { boardApi } from "@/shared/api/modules/boards";
import { postApi } from "@/shared/api/modules/posts";
import { tagApi } from "@/shared/api/modules/tags";
import { Tag } from "@/shared/api/types/tag.model";
import { ImageUploader, ImageItem } from "@/shared/ui/editor/ImageUploader";
import { uploadApi } from "@/shared/api/modules/uploads";
import DOMPurify from "dompurify";
import { RichTextEditor } from "@/shared/ui/editor/richtext/RichTextEditor";
// import { RichTextEditor } from "@/shared/ui/editor/RichTextEditor";

// ---------- 表单验证 ----------
const postSchema = z
  .object({
    title: z.string().min(2, "标题至少2个字符").max(200, "标题最多200个字符"),
    content: z.string(),
    summary: z.string().max(500).optional(),
    cover: z.string().url("请输入有效的图片URL").optional().or(z.literal("")),
    type: z.enum(["post", "article", "topic"]),
    board_id: z.number().min(1, "请选择板块"),
    tag_ids: z.array(z.number()).max(5, "最多选择5个标签"),
    status: z
      .enum(["draft", "published", "pending", "hidden"])
      .default("published"),
  })
  .superRefine((data, ctx) => {
    if (
      data.status !== "draft" &&
      (!data.content || data.content.length < 10)
    ) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["content"],
        message: "内容至少10个字符",
      });
    }
  });

type PostForm = z.infer<typeof postSchema>;

interface BoardListResponse {
  list: Board[];
  total: number;
  page: number;
  page_size: number;
}

// ---------- 左侧：帖子设置组件（类型优化）----------
interface PostSettingsProps {
  register: UseFormRegister<PostForm>;
  errors: FieldErrors<PostForm>;
  boards: Board[];
  tags: Tag[];
  boardsLoading: boolean;
  selectedTagIds: number[];
  selectedStatus: string;
  coverValue: string;
  onToggleTag: (tagId: number) => void;
  onCoverChange: (url: string) => void;
  t: (key: string) => string;
}

function PostSettings({
  register,
  errors,
  boards,
  tags,
  boardsLoading,
  selectedTagIds,
  selectedStatus,
  coverValue,
  onToggleTag,
  onCoverChange,
  t,
}: PostSettingsProps) {
  const statusOptions = [
    { value: "draft", label: t("status_draft"), desc: t("status_draft_desc") },
    {
      value: "published",
      label: t("status_published"),
      desc: t("status_published_desc"),
    },
    {
      value: "pending",
      label: t("status_pending"),
      desc: t("status_pending_desc"),
    },
    {
      value: "hidden",
      label: t("status_hidden"),
      desc: t("status_hidden_desc"),
    },
  ];

  // 封面上传函数（无需 postId）
  const handleUploadCover = async (file: File): Promise<{ url: string }> => {
    try {
      // 使用通用附件上传接口（需后端支持，返回 { data: { url } }）
      const res = await uploadApi.uploadPluginFile(file, "post_cover");
      return { url: res.data.data };
    } catch (error) {
      toast.error(t("cover_upload_failed"));
      throw error;
    }
  };

  return (
    <div className="space-y-5">
      {/* 板块选择 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            <FolderOpen className="w-4 h-4 inline mr-1" />
            {t("board")} <span className="text-error">*</span>
          </span>
        </label>
        <select
          {...register("board_id", {
            required: "请选择板块",
            valueAsNumber: true,
          })}
          className={`select select-bordered w-full focus:outline-none focus:border-primary ${
            errors.board_id ? "select-error" : ""
          }`}
          defaultValue=""
        >
          <option value="" disabled>
            {t("select_board")}
          </option>
          {boardsLoading ? (
            <option disabled>{t("loading")}</option>
          ) : (
            boards.map((board) => (
              <option key={board.id} value={board.id}>
                {board.name} {board.description ? `- ${board.description}` : ""}
              </option>
            ))
          )}
        </select>
        {errors.board_id && (
          <label className="label pt-1">
            <span className="label-text-alt text-error">
              {errors.board_id.message}
            </span>
          </label>
        )}
      </div>

      {/* 帖子类型 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">{t("post_type")}</span>
        </label>
        <div className="flex gap-2">
          {[
            { value: "post", label: t("post"), desc: t("post_desc") },
            { value: "article", label: t("article"), desc: t("article_desc") },
            { value: "topic", label: t("topic"), desc: t("topic_desc") },
          ].map((typeOption) => (
            <label key={typeOption.value} className="flex-1 cursor-pointer">
              <input
                {...register("type")}
                type="radio"
                value={typeOption.value}
                className="hidden peer"
              />
              <div className="border-2 border-base-300 rounded-xl p-3 text-center peer-checked:border-primary peer-checked:bg-primary/5 transition-all">
                <div className="font-medium text-sm">{typeOption.label}</div>
                <div className="text-xs text-base-content/40 mt-0.5">
                  {typeOption.desc}
                </div>
              </div>
            </label>
          ))}
        </div>
      </div>

      {/* 状态 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("status")} <span className="text-error">*</span>
          </span>
        </label>
        <select
          {...register("status")}
          className="select select-bordered w-full"
        >
          {statusOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label} - {opt.desc}
            </option>
          ))}
        </select>
        {selectedStatus === "draft" && (
          <label className="label pt-1">
            <span className="label-text-alt text-info">{t("draft_hint")}</span>
          </label>
        )}
      </div>

      {/* 标题 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("post_title")} <span className="text-error">*</span>
          </span>
        </label>
        <input
          {...register("title")}
          type="text"
          placeholder={t("post_title_placeholder")}
          className={`input input-bordered focus:outline-none focus:border-primary ${
            errors.title ? "input-error" : ""
          }`}
        />
        {errors.title && (
          <label className="label pt-1">
            <span className="label-text-alt text-error">
              {errors.title.message}
            </span>
          </label>
        )}
      </div>

      {/* 标签 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("tags")}
            <span className="text-base-content/40 text-xs ml-2">
              {t("select_up_to_tags")}
            </span>
          </span>
        </label>
        <div className="flex flex-wrap gap-2">
          {tags.map((tag) => {
            const selected = selectedTagIds?.includes(tag.id);
            return (
              <button
                key={tag.id}
                type="button"
                onClick={() => onToggleTag(tag.id)}
                className={`badge badge-lg cursor-pointer transition-all ${
                  selected ? "ring-2" : "opacity-60 hover:opacity-100"
                }`}
                style={{
                  backgroundColor: selected
                    ? tag.color + "30"
                    : tag.color + "15",
                  color: tag.color,
                  borderColor: tag.color + "60",
                }}
              >
                {selected && <X className="w-3 h-3 mr-1" />}
                {tag.name}
              </button>
            );
          })}
        </div>
      </div>

      {/* 封面图 - 绑定表单 cover 字段 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">{t("cover")}</span>
        </label>
        <ImageUploader
          initialImages={coverValue ? [{ url: coverValue }] : []}
          uploadFn={handleUploadCover}
          maxCount={1}
          supportCover={false}
          layout="grid"
          gridSize={2}
          onChange={(images: ImageItem[]) => {
            const coverUrl = images.length > 0 ? images[0].url : "";
            onCoverChange(coverUrl);
          }}
        />
        <label className="label">
          <span className="label-text-alt text-base-content/40">
            {t("cover_hint")}
          </span>
        </label>
      </div>

      {/* 摘要 */}
      <div className="form-control">
        <label className="label pb-1">
          <span className="label-text font-medium">
            {t("summary")}
            <span className="text-base-content/40 text-xs ml-2">
              {t("summary_desc")}
            </span>
          </span>
        </label>
        <textarea
          {...register("summary")}
          rows={2}
          placeholder={t("summary_placeholder")}
          className="textarea textarea-bordered focus:outline-none focus:border-primary resize-none"
        />
      </div>
    </div>
  );
}

// ---------- 右侧：编辑器 + 预览组件 ----------
interface PostEditorProps {
  content: string;
  setContent: (value: string) => void;
  placeholder?: string;
}

function PostEditor({ content, setContent }: PostEditorProps) {
  const [mode, setMode] = useState<"edit" | "preview">("edit");
  const sanitizedHtml = useMemo(() => DOMPurify.sanitize(content), [content]);

  return (
    <div className="border border-base-300 rounded-lg bg-base-100 shadow-sm flex flex-col h-full">
      <div className="flex justify-between items-center border-b border-base-300 p-2 bg-base-200 rounded-t-lg">
        <div className="flex gap-1">
          <button
            type="button"
            onClick={() => setMode("edit")}
            className={`btn btn-xs gap-1 ${mode === "edit" ? "btn-primary" : "btn-ghost"}`}
          >
            <Edit className="w-3 h-3" />
            编辑
          </button>
          <button
            type="button"
            onClick={() => setMode("preview")}
            className={`btn btn-xs gap-1 ${mode === "preview" ? "btn-primary" : "btn-ghost"}`}
          >
            <Eye className="w-3 h-3" />
            预览
          </button>
        </div>
        {mode === "preview" && (
          <span className="text-xs text-base-content/50">纯预览模式</span>
        )}
      </div>

      <div className="flex-1 min-h-[400px]">
        {mode === "edit" ? (
          // 添加 key={mode} 强制重新创建编辑器实例，避免 DOM 冲突
          <RichTextEditor
            key="rich-editor"
            value={content}
            onChange={setContent}
            placeholder="撰写帖子内容..."
            maxLength={20000}
            defaultMode="rich"
          />
        ) : (
          <div className="prose prose-sm max-w-none p-4 overflow-auto h-full">
            {content ? (
              <div dangerouslySetInnerHTML={{ __html: sanitizedHtml }} />
            ) : (
              <p className="text-base-content/40">暂无内容</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

// ---------- 主页面 ----------
export default function NewPostPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const t = useTranslations("Post");

  useEffect(() => {
    if (!isAuthenticated) router.push("/auth/login");
  }, [isAuthenticated, router]);

  // 获取板块列表
  const { data: boardsData, isLoading: boardsLoading } = useQuery({
    queryKey: ["boards"],
    queryFn: () =>
      boardApi
        .list()
        .then(
          (r: { data: ApiResponse<BoardListResponse> }) =>
            r.data.data?.list || [],
        ),
  });

  // 获取标签列表
  const { data: tagsData } = useQuery({
    queryKey: ["tags"],
    queryFn: () =>
      tagApi
        .list()
        .then((r: { data: ApiResponse<Tag[]> }) => r.data.data || []),
  });

  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: { errors },
  } = useForm<PostForm>({
    resolver: zodResolver(postSchema),
    defaultValues: {
      type: "post",
      board_id: undefined,
      tag_ids: [],
      status: "published",
      content: "",
      cover: "",
    },
  });

  const selectedTagIds = useWatch({
    control,
    name: "tag_ids",
    defaultValue: [],
  });
  const selectedStatus = useWatch({
    control,
    name: "status",
    defaultValue: "published",
  });
  const coverValue = useWatch({ control, name: "cover", defaultValue: "" });

  const toggleTag = useCallback(
    (tagId: number) => {
      const current = selectedTagIds ?? [];
      if (current.includes(tagId)) {
        setValue(
          "tag_ids",
          current.filter((id: number) => id !== tagId),
        );
      } else if (current.length < 5) {
        setValue("tag_ids", [...current, tagId]);
      } else {
        toast.error(t("select_up_to_tags"));
      }
    },
    [selectedTagIds, setValue, t],
  );

  const onSubmit = useCallback(
    async (data: PostForm) => {
      if (!data.board_id) {
        toast.error(t("select_board_error"));
        return;
      }

      setLoading(true);
      try {
        const response = await postApi.create({
          ...data,
          cover: data.cover || undefined,
          summary: data.summary || undefined,
        });
        toast.success(t("publish_success"));
        const postId = response.data.data?.id;
        router.push(postId ? `/posts/${postId}` : "/posts");
      } catch (err) {
        toast.error(getErrorMessage(err));
      } finally {
        setLoading(false);
      }
    },
    [router, t],
  );

  if (!isAuthenticated) return null;

  const boards = boardsData ?? [];
  const tags = tagsData ?? [];

  return (
    <div className="max-w-7xl mx-auto px-4 py-6">
      <div className="flex items-center gap-3 mb-6">
        <FileText className="w-6 h-6 text-primary" />
        <h1 className="text-2xl font-bold">{t("publish_new_post")}</h1>
      </div>

      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="flex flex-col lg:flex-row gap-6">
          {/* 左侧：设置区 */}
          <div className="lg:w-80 flex-shrink-0">
            <div className="card bg-base-100 border border-base-300 shadow-sm sticky top-20">
              <div className="card-body p-5">
                <PostSettings
                  register={register}
                  errors={errors}
                  boards={boards}
                  tags={tags}
                  boardsLoading={boardsLoading}
                  selectedTagIds={selectedTagIds}
                  selectedStatus={selectedStatus}
                  coverValue={coverValue || ""}
                  onToggleTag={toggleTag}
                  onCoverChange={(url) => setValue("cover", url)}
                  t={t}
                />
              </div>
            </div>
          </div>

          {/* 右侧：编辑器区 */}
          <div className="flex-1 min-w-0">
            <Controller
              name="content"
              control={control}
              render={({ field }) => (
                <PostEditor
                  content={field.value}
                  setContent={field.onChange}
                  placeholder={t("post_content_placeholder")}
                />
              )}
            />
            {errors.content && (
              <p className="text-error text-sm mt-2">
                {errors.content.message}
              </p>
            )}

            <div className="flex gap-3 justify-end mt-6">
              <button
                type="button"
                className="btn btn-ghost"
                onClick={() => router.back()}
              >
                {t("cancel")}
              </button>
              <button
                type="submit"
                className="btn btn-primary gap-2"
                disabled={loading}
              >
                {loading ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  <Send className="w-4 h-4" />
                )}
                {selectedStatus === "draft"
                  ? t("save_draft")
                  : t("publish_post")}
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
}
