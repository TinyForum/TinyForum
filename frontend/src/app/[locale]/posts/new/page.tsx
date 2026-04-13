'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { postApi, tagApi, boardApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import RichEditor from '@/components/post/RichEditor';
import toast from 'react-hot-toast';
import { getErrorMessage } from '@/lib/utils';
import { FileText, Send, X, FolderOpen } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import type { Board } from '@/lib/api/types';

const postSchema = z.object({
  title: z.string().min(2, '标题至少2个字符').max(200, '标题最多200个字符'),
  content: z.string().min(10, '内容至少10个字符'),
  summary: z.string().max(500).optional(),
  cover: z.string().url('请输入有效的图片URL').optional().or(z.literal('')),
  type: z.enum(['post', 'article', 'topic']),
  board_id: z.number().min(1, '请选择板块'),  // 添加板块验证
  tag_ids: z.array(z.number()).max(5, '最多选择5个标签'),
});

type PostForm = z.infer<typeof postSchema>;

export default function NewPostPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const t = useTranslations('posts');

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/auth/login');
    }
  }, [isAuthenticated, router]);

  // 获取板块列表
  const { data: boards, isLoading: boardsLoading } = useQuery({
    queryKey: ['boards'],
    queryFn: () => boardApi.list().then((r) => r.data.data.list || []),
  });

  // 获取标签列表
  const { data: tags } = useQuery({
    queryKey: ['tags'],
    queryFn: () => tagApi.list().then((r) => r.data.data),
  });

  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useForm<PostForm>({
    resolver: zodResolver(postSchema),
    defaultValues: {
      type: 'post',
      board_id: undefined,
      tag_ids: [],
    },
  });

  const selectedTagIds = watch('tag_ids');
  const selectedBoardId = watch('board_id');

  const toggleTag = (tagId: number) => {
    const current = selectedTagIds ?? [];
    if (current.includes(tagId)) {
      setValue('tag_ids', current.filter((id) => id !== tagId));
    } else if (current.length < 5) {
      setValue('tag_ids', [...current, tagId]);
    } else {
      toast.error(t("select_up_to_tags"));
    }
  };

  const onSubmit = async (data: PostForm) => {
    // 额外验证板块
    if (!data.board_id) {
      toast.error('请选择板块');
      return;
    }

    setLoading(true);
    const requestBody = {
 ...data,
        board_id: data.board_id,  // 确保传递 board_id
        cover: data.cover || undefined,
        summary: data.summary || undefined,
    }
 console.log(requestBody);
    try {
      const response = await postApi.create(requestBody);
     
      toast.success(t("publish_success"));
      console.log(response);
      router.push(`/posts/${response.data.data.id}`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!isAuthenticated) return null;

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center gap-3 mb-6">
        <FileText className="w-6 h-6 text-primary" />
        <h1 className="text-2xl font-bold">{t("publish_new_post")}</h1>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-5 space-y-4">
            {/* 板块选择 - 新增 */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  <FolderOpen className="w-4 h-4 inline mr-1" />
                  选择板块 <span className="text-error">*</span>
                </span>
              </label>
              <select
                {...register('board_id', { required: '请选择板块', valueAsNumber: true })}
                className={`select select-bordered w-full focus:outline-none focus:border-primary ${
                  errors.board_id ? 'select-error' : ''
                }`}
                defaultValue=""
              >
                <option value="" disabled>请选择板块</option>
                {boardsLoading ? (
                  <option disabled>加载中...</option>
                ) : (
                  (boards ?? []).map((board: Board) => (
                    <option key={board.id} value={board.id}>
                      {board.name} {board.description ? `- ${board.description}` : ''}
                    </option>
                  ))
                )}
              </select>
              {errors.board_id && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error">{errors.board_id.message}</span>
                </label>
              )}
              {!selectedBoardId && !errors.board_id && (
                <label className="label pt-1">
                  <span className="label-text-alt text-base-content/40">
                    请选择帖子所属的板块
                  </span>
                </label>
              )}
            </div>

            {/* Type */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">{t("post_type")}</span>
              </label>
              <div className="flex gap-2">
                {[
                  { value: 'post', label: t("post"), desc: t("post_desc") },
                  { value: 'article', label: t("article"), desc: t("article_desc") },
                  { value: 'topic', label: t("topic"), desc: t("topic_desc") },
                ].map((typeOption) => (
                  <label key={typeOption.value} className="flex-1 cursor-pointer">
                    <input 
                      {...register('type')} 
                      type="radio" 
                      value={typeOption.value} 
                      className="hidden peer" 
                    />
                    <div className="border-2 border-base-300 rounded-xl p-3 text-center peer-checked:border-primary peer-checked:bg-primary/5 transition-all">
                      <div className="font-medium text-sm">{typeOption.label}</div>
                      <div className="text-xs text-base-content/40 mt-0.5">{typeOption.desc}</div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Title */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  {t("post_title")} <span className="text-error">*</span>
                </span>
              </label>
              <input
                {...register('title')}
                type="text"
                placeholder={t("post_title_placeholder")}
                className={`input input-bordered focus:outline-none focus:border-primary ${errors.title ? 'input-error' : ''}`}
              />
              {errors.title && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error">{errors.title.message}</span>
                </label>
              )}
            </div>

            {/* Tags */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  {t("tags")} 
                  <span className="text-base-content/40 text-xs ml-2">{t("select_up_to_tags")}</span>
                </span>
              </label>
              <div className="flex flex-wrap gap-2">
                {(tags ?? []).map((tag) => {
                  const selected = selectedTagIds?.includes(tag.id);
                  return (
                    <button
                      key={tag.id}
                      type="button"
                      onClick={() => toggleTag(tag.id)}
                      className={`badge badge-lg cursor-pointer transition-all ${
                        selected ? 'ring-2' : 'opacity-60 hover:opacity-100'
                      }`}
                      style={{
                        backgroundColor: selected ? tag.color + '30' : tag.color + '15',
                        color: tag.color,
                        borderColor: tag.color + '60',
                      }}
                    >
                      {selected && <X className="w-3 h-3 mr-1" />}
                      {tag.name}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Cover image URL */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  {t("cover_image")} 
                  <span className="text-base-content/40 text-xs ml-2">{t("cover_image_desc")}</span>
                </span>
              </label>
              <input
                {...register('cover')}
                type="text"
                placeholder="https://example.com/image.jpg"
                className={`input input-bordered focus:outline-none focus:border-primary ${errors.cover ? 'input-error' : ''}`}
              />
              {errors.cover && (
                <label className="label pt-1">
                  <span className="label-text-alt text-error">{errors.cover.message}</span>
                </label>
              )}
            </div>

            {/* Summary */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">
                  {t("summary")} 
                  <span className="text-base-content/40 text-xs ml-2">{t("summary_desc")}</span>
                </span>
              </label>
              <textarea
                {...register('summary')}
                rows={2}
                placeholder={t("summary_placeholder")}
                className="textarea textarea-bordered focus:outline-none focus:border-primary resize-none"
              />
            </div>
          </div>
        </div>

        {/* Content editor */}
        <div>
          <label className="label pb-2">
            <span className="label-text font-medium text-base">
              {t("post_content")}<span className="text-error">*</span>
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

        {/* Submit */}
        <div className="flex gap-3 justify-end">
          <button type="button" className="btn btn-ghost" onClick={() => router.back()}>
            {t("cancel")}
          </button>
          <button type="submit" className="btn btn-primary gap-2" disabled={loading}>
            {loading ? (
              <span className="loading loading-spinner loading-sm" />
            ) : (
              <Send className="w-4 h-4" />
            )}
            {t("publish_post")}
          </button>
        </div>
      </form>
    </div>
  );
}