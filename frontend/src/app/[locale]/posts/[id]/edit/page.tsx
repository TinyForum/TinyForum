"use client";

import { useEffect, useState, use } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { useForm, Controller, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Post, postApi, Tag, tagApi } from "@/shared/api";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { getErrorMessage } from "@/shared/lib/utils";
import { Save, X } from "lucide-react";
import RichEditor from "@/layout/common/RichEditor";
import { ApiResponse } from "@/shared/api/types/basic.model";

const schema = z.object({
  title: z.string().min(2, "标题至少2个字符").max(200, "标题最多200个字符"),
  content: z.string().min(10, "内容至少10个字符"),
  summary: z.string().max(500, "摘要最多500个字符").optional(),
  cover: z.string().url("请输入有效的图片URL").optional().or(z.literal("")),
  tag_ids: z.array(z.number()),
});

type EditForm = z.infer<typeof schema>;

interface PostDetailResponse {
  post: Post;
}

export default function EditPostPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  // 使用 React.use 解析 Promise 参数（Next.js 15 特性）
  const { id } = use(params);
  const postId = Number(id);
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);

  const { data: postData } = useQuery({
    queryKey: ["post", postId],
    queryFn: () =>
      postApi
        .getById(postId)
        .then((r: { data: ApiResponse<PostDetailResponse> }) => r.data.data),
  });

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
    reset,
    formState: { errors },
  } = useForm<EditForm>({
    resolver: zodResolver(schema),
    defaultValues: {
      tag_ids: [],
      title: "",
      content: "",
      summary: "",
      cover: "",
    },
  });

  useEffect(() => {
    if (postData?.post) {
      const p = postData.post;
      reset({
        title: p.title,
        content: p.content,
        summary: p.summary || "",
        cover: p.cover || "",
        tag_ids: p.tags?.map((t: Tag) => t.id) ?? [],
      });
    }
  }, [postData, reset]);

  useEffect(() => {
    if (!isAuthenticated) router.push("/auth/login");
  }, [isAuthenticated, router]);

  // 使用 useWatch 替代 watch
  const selectedTagIds = useWatch({
    control,
    name: "tag_ids",
    defaultValue: [],
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
      await postApi.update(postId, {
        ...data,
        cover: data.cover || undefined,
        summary: data.summary || undefined,
      });
      toast.success("更新成功");
      router.push(`/posts/${postId}`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!postData) {
    return (
      <div className="flex justify-center py-20">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  const tags = tagsData ?? [];

  return (
    <div className="max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">编辑帖子</h1>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-5 space-y-4">
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">标题</span>
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

            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">标签</span>
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

            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">封面图片URL</span>
              </label>
              <input
                {...register("cover")}
                type="text"
                placeholder="https://example.com/image.jpg"
                className="input input-bordered focus:outline-none focus:border-primary"
              />
              {errors.cover && (
                <span className="text-error text-sm mt-1">
                  {errors.cover.message}
                </span>
              )}
            </div>

            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">摘要</span>
              </label>
              <textarea
                {...register("summary")}
                rows={2}
                placeholder="简要描述文章内容..."
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

        <div>
          <label className="label pb-2">
            <span className="label-text font-medium text-base">正文内容</span>
          </label>
          <Controller
            name="content"
            control={control}
            render={({ field }) => (
              <RichEditor
                content={field.value}
                onChange={field.onChange}
                placeholder="请输入帖子内容..."
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
            取消
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
            保存修改
          </button>
        </div>
      </form>
    </div>
  );
}
