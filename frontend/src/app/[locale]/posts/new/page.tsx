"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useForm, Controller, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { postApi, tagApi, boardApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import RichEditor from "@/components/post/RichEditor";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/lib/utils";
import { FileText, Send, X, FolderOpen } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import type { Board, Tag, ApiResponse } from "@/lib/api/types";

// 增强的校验规则：草稿状态下内容最短可为 0，其他状态至少 10 字符
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
    // 动态内容校验：非草稿状态内容必须 >= 10 字符
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

export default function NewPostPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const t = useTranslations("posts");

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/auth/login");
    }
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
    },
  });

  // 使用 useWatch 替代 watch
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

  const toggleTag = (tagId: number) => {
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
  };

  const onSubmit = async (data: PostForm) => {
    if (!data.board_id) {
      toast.error("请选择板块");
      return;
    }

    setLoading(true);
    const requestBody = {
      ...data,
      board_id: data.board_id,
      cover: data.cover || undefined,
      summary: data.summary || undefined,
      status: data.status,
    };
    console.log(requestBody);

    try {
      const response = await postApi.create(requestBody);
      toast.success(t("publish_success"));
      console.log(response);
      const postId = response.data.data?.id;
      if (postId) {
        router.push(`/posts/${postId}`);
      } else {
        router.push("/posts");
      }
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!isAuthenticated) return null;

  const boards = boardsData ?? [];
  const tags = tagsData ?? [];

  // 状态选项配置
  const statusOptions = [
    { value: "draft", label: "草稿", desc: "仅自己可见，可继续编辑" },
    { value: "published", label: "发布", desc: "直接公开" },
    { value: "pending", label: "待审核", desc: "提交后等待管理员审核" },
    { value: "hidden", label: "隐藏", desc: "不公开，仅管理员可见" },
  ];

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center gap-3 mb-6">
        <FileText className="w-6 h-6 text-primary" />
        <h1 className="text-2xl font-bold">{t("publish_new_post")}</h1>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-5 space-y-4">
            {/* 板块选择 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  <FolderOpen className="w-4 h-4 inline mr-1" />
                  选择板块 <span className="text-error">*</span>
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
                  请选择板块
                </option>
                {boardsLoading ? (
                  <option disabled>加载中...</option>
                ) : (
                  boards.map((board: Board) => (
                    <option key={board.id} value={board.id}>
                      {board.name}{" "}
                      {board.description ? `- ${board.description}` : ""}
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
                  {
                    value: "article",
                    label: t("article"),
                    desc: t("article_desc"),
                  },
                  { value: "topic", label: t("topic"), desc: t("topic_desc") },
                ].map((typeOption) => (
                  <label
                    key={typeOption.value}
                    className="flex-1 cursor-pointer"
                  >
                    <input
                      {...register("type")}
                      type="radio"
                      value={typeOption.value}
                      className="hidden peer"
                    />
                    <div className="border-2 border-base-300 rounded-xl p-3 text-center peer-checked:border-primary peer-checked:bg-primary/5 transition-all">
                      <div className="font-medium text-sm">
                        {typeOption.label}
                      </div>
                      <div className="text-xs text-base-content/40 mt-0.5">
                        {typeOption.desc}
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* 文章状态 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  状态 <span className="text-error">*</span>
                </span>
              </label>
              <select
                {...register("status")}
                className="select select-bordered w-full focus:outline-none focus:border-primary"
              >
                {statusOptions.map((opt) => (
                  <option key={opt.value} value={opt.value}>
                    {opt.label} - {opt.desc}
                  </option>
                ))}
              </select>
              {selectedStatus === "draft" && (
                <label className="label pt-1">
                  <span className="label-text-alt text-info">
                    草稿状态下内容长度不受限制，方便随时保存
                  </span>
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
                {tags.map((tag: Tag) => {
                  const selected = selectedTagIds?.includes(tag.id);
                  return (
                    <button
                      key={tag.id}
                      type="button"
                      onClick={() => toggleTag(tag.id)}
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

            {/* 封面图 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  {t("cover_image")}
                  <span className="text-base-content/40 text-xs ml-2">
                    {t("cover_image_desc")}
                  </span>
                </span>
              </label>
              <input
                {...register("cover")}
                type="text"
                placeholder="https://example.com/image.jpg"
                className={`input input-bordered focus:outline-none focus:border-primary ${
                  errors.cover ? "input-error" : ""
                }`}
              />
              {errors.cover && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error">
                    {errors.cover.message}
                  </span>
                </label>
              )}
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
        </div>

        {/* 富文本编辑器 */}
        <div>
          <label className="label pb-2">
            <span className="label-text font-medium text-base">
              {t("post_content")}
              {selectedStatus !== "draft" && (
                <span className="text-error">*</span>
              )}
            </span>
          </label>
          <Controller
            name="content"
            control={control}
            render={({ field }) => (
              <RichEditor
                content={field.value}
                onChange={field.onChange}
                placeholder={t("post_content_placeholder")}
              />
            )}
          />
          {errors.content && (
            <p className="text-error text-sm mt-1">{errors.content.message}</p>
          )}
        </div>

        {/* 提交按钮 */}
        <div className="flex gap-3 justify-end">
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
            {selectedStatus === "draft" ? "保存草稿" : t("publish_post")}
          </button>
        </div>
      </form>
    </div>
  );
}
