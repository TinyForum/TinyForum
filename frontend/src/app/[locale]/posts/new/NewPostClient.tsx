// src/app/[locale]/posts/new/page.tsx 客户端组件
"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useForm, Controller, useWatch } from "react-hook-form";
import type { FieldErrors } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import { FileText, Send } from "lucide-react";
import { useTranslations } from "next-intl";
import { useBoard } from "@/features/boards/hooks/useBoard";
import { useTags } from "@/features/tag/hooks/useTags";
import { useCreatePost } from "@/features/post/hooks/usePosts";
import { PostEditor } from "./PostEditor";
import { PostSettings } from "./PostSettings";
import { PostForm, postSchema } from "./newPost.types";
import { PostType } from "@/shared/api/types/post.model";

export default function NewPostClient() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const t = useTranslations("Post");
  const { mutateAsync: createPost } = useCreatePost();

  useEffect(() => {
    if (!isAuthenticated) router.push("/auth/login");
  }, [isAuthenticated, router]);

  // 获取板块列表
  const { boards: boardsData, loading: boardsLoading } = useBoard();

  // 获取标签列表
  const { tags: tagsData } = useTags();

  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: { errors },
  } = useForm<PostForm>({
    resolver: zodResolver(postSchema),
    defaultValues: {
      type: "image_text",
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
  const selectedType = useWatch({ control, name: "type", defaultValue: "image_text" }) as PostType;
  const [imageUrls, setImageUrls] = useState<string[]>([]);

  const handleImageUrlsChange = useCallback((urls: string[]) => {
    setImageUrls((prev) => {
      if (prev.length === urls.length && prev.every((u, i) => u === urls[i])) return prev;
      return urls;
    });
  }, []);

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
        let content = data.content;
        if (imageUrls.length > 0) {
          const imgTags = imageUrls.map((url) => `<img src="${url}" alt="" />`).join("");
          content = imgTags + content;
        }

        const created = await createPost({
          ...data,
          content,
          cover: data.cover || undefined,
          summary: data.summary || undefined,
          video_url: data.video_url || undefined,
        });
        toast.success(t("publish_success"));
        const postId = created?.id;
        router.push(postId ? `/posts/${postId}` : "/posts");
      } catch (err) {
        toast.error(getErrorMessage(err));
      } finally {
        setLoading(false);
      }
    },
    [createPost, router, t, imageUrls],
  );

  const onError = useCallback(
    (errors: FieldErrors<PostForm>) => {
      const firstError = Object.values(errors).find(
        (e) => e?.message,
      );
      if (firstError?.message) {
        toast.error(String(firstError.message));
      }
    },
    [],
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

      <form onSubmit={handleSubmit(onSubmit, onError)}>
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
                  selectedType={selectedType}
                  coverValue={coverValue || ""}
                  onToggleTag={toggleTag}
                  onCoverChange={(url) => setValue("cover", url)}
                  onVideoChange={(videoUrl) => setValue("video_url", videoUrl)}
                  onImageUrlsChange={handleImageUrlsChange}
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
