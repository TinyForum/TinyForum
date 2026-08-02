// src/app/[locale]/posts/[id]/edit/EditPostClient.tsx
"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useForm, Controller, useWatch } from "react-hook-form";
import type { FieldErrors } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import { Save, X } from "lucide-react";
import { useTranslations } from "next-intl";
import { Tag } from "@/shared/api/types/tag.model";
import { PostDetailResponse } from "@/shared/api/types/post.model";
import { usePost, useUpdatePost } from "@/features/post/hooks/usePosts";
import { useTags } from "@/features/tag/hooks/useTags";
import { uploadApi } from "@/shared/api/modules/uploads";
import { ImageUploader } from "@/shared/ui/upload/ImageUploader";
import { RichTextEditor } from "@/shared/ui/editor/richtext/RichTextEditor";

function stripHtml(html: string): string {
  return html.replace(/<[^>]*>/g, "").trim();
}

function isValidImageUrl(val: string): boolean {
  if (!val) return true;
  return /^(\/|https?:\/\/)/.test(val);
}

const schema = z.object({
  title: z.string().min(2, "标题至少2个字符").max(200, "标题最多200个字符"),
  content: z.string().refine(
    (val) => stripHtml(val).length >= 10,
    "内容至少10个字符",
  ),
  summary: z.string().max(500, "摘要最多500个字符").optional(),
  cover: z.string().refine(isValidImageUrl, "请输入有效的图片URL").optional().or(z.literal("")),
  tag_ids: z.array(z.number()),
});

type EditForm = z.infer<typeof schema>;

/** 数据加载层：在 postData 就绪前展示骨架屏 */
export default function EditPostClient({ postId }: { postId: number }) {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const { data: postData } = usePost(postId);
  const { tags: tagsData } = useTags();

  useEffect(() => {
    if (!isAuthenticated) router.push("/auth/login");
  }, [isAuthenticated, router]);

  if (!postData) {
    return (
      <div className="flex justify-center py-20">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  return (
    <EditPostForm
      postId={postId}
      postData={postData}
      tags={tagsData ?? []}
    />
  );
}

/** 表单组件：仅在 postData 就绪时渲染，避免 defaultValues 空值与 useEffect reset 竞态 */
function EditPostForm({
  postId,
  postData,
  tags,
}: {
  postId: number;
  postData: PostDetailResponse;
  tags: Tag[];
}) {
  const router = useRouter();
  const t = useTranslations("Post");
  const [loading, setLoading] = useState(false);
  const { mutateAsync: updatePost } = useUpdatePost();

  const p = postData.post;

  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: { errors },
  } = useForm<EditForm>({
    resolver: zodResolver(schema),
    defaultValues: {
      title: p.creation.title,
      content: p.creation.content,
      summary: p.creation.summary || "",
      cover: p.creation.cover_url || "",
      tag_ids: p.creation.tags?.map((tag: Tag) => tag.id) ?? [],
    },
  });

  const selectedTagIds = useWatch({
    control,
    name: "tag_ids",
    defaultValue: p.creation.tags?.map((tag: Tag) => tag.id) ?? [],
  });

  const toggleTag = (tagId: number) => {
    const current = selectedTagIds ?? [];
    setValue(
      "tag_ids",
      current.includes(tagId)
        ? current.filter((id: number) => id !== tagId)
        : [...current, tagId],
    );
  };

  const onSubmit = async (data: EditForm) => {
    setLoading(true);
    try {
      await updatePost({
        id: postId,
        data: {
          ...data,
          cover: data.cover || undefined,
          summary: data.summary || undefined,
        },
      });
      toast.success(t("update_success"));
      router.push(`/posts/${postId}`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  const onError = useCallback(
    (errors: FieldErrors<EditForm>) => {
      const firstError = Object.values(errors).find((e) => e?.message);
      if (firstError?.message) {
        toast.error(String(firstError.message));
      }
    },
    [],
  );

  const initialImages =
    p.creation.image_urls?.map((url: string) => ({
      url,
      isCover: url === p.creation.cover_url,
    })) ?? [];

  return (
    <div className="max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">{t("edit_post")}</h1>
      <form onSubmit={handleSubmit(onSubmit, onError)} className="space-y-5">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-5 space-y-4">
            {/* 标题 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">{t("post_title")}</span>
              </label>
              <input
                {...register("title")}
                className={`input input-bordered focus:outline-none focus:border-primary ${errors.title ? "input-error" : ""}`}
              />
              {errors.title && (
                <span className="text-error text-sm mt-1">
                  {errors.title.message}
                </span>
              )}
            </div>

            {/* 标签 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">{t("tags")}</span>
              </label>
              <div className="flex flex-wrap gap-2">
                {tags.map((tag: Tag) => {
                  const selected = selectedTagIds?.includes(tag.id);
                  return (
                    <button
                      key={tag.id}
                      type="button"
                      onClick={() => toggleTag(tag.id)}
                      className={`badge badge-lg cursor-pointer transition-all ${selected ? "ring-2" : "opacity-60 hover:opacity-100"}`}
                      style={{
                        backgroundColor: tag.color + "20",
                        color: tag.color,
                        borderColor: tag.color + "40",
                      }}
                    >
                      {selected && <X className="w-3 h-3 mr-1" />}
                      {tag.name}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* 图片管理 */}
            <ImageUploader
              initialImages={initialImages}
              uploadFn={async (file) => {
                const res = await uploadApi.uploadPostFile(postId, file);
                return {
                  url: res.data.data?.url ?? "",
                  file_id: res.data.data?.file_id,
                };
              }}
              maxCount={6}
              supportCover={true}
              layout="grid"
              gridSize={3}
              onChange={(images) => {
                const cover = images.find((img) => img.isCover);
                if (cover) {
                  setValue("cover", cover.url);
                }
              }}
            />

            {/* 摘要 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">{t("summary")}</span>
              </label>
              <textarea
                {...register("summary")}
                rows={2}
                placeholder={t("summary_placeholder")}
                className="textarea textarea-bordered focus:outline-none focus:border-primary resize-none"
              />
              {errors.summary && (
                <span className="text-error text-sm mt-1">
                  {errors.summary.message}
                </span>
              )}
            </div>
          </div>
        </div>

        {/* 正文编辑器 */}
        <div>
          <label className="label pb-2">
            <span className="label-text font-medium text-base">
              {t("post_content")}
            </span>
          </label>
          <Controller
            name="content"
            control={control}
            render={({ field }) => (
              <RichTextEditor
                key={`edit-post-${postId}`}
                value={field.value}
                onChange={field.onChange}
                placeholder={t("post_content_placeholder")}
                maxLength={20000}
                defaultMode="rich"
              />
            )}
          />
          {errors.content && (
            <p className="text-error text-sm mt-1">{errors.content.message}</p>
          )}
        </div>

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
              <Save className="w-4 h-4" />
            )}
            {t("save_changes")}
          </button>
        </div>
      </form>
    </div>
  );
}
