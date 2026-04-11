'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { postApi, tagApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import RichEditor from '@/components/post/RichEditor';
import toast from 'react-hot-toast';
import { getErrorMessage } from '@/lib/utils';
import { FileText, Send, X } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';

const postSchema = z.object({
  title: z.string().min(2, '标题至少2个字符').max(200, '标题最多200个字符'),
  content: z.string().min(10, '内容至少10个字符'),
  summary: z.string().max(500).optional(),
  cover: z.string().url('请输入有效的图片URL').optional().or(z.literal('')),
  type: z.enum(['post', 'article', 'topic']),
  tag_ids: z.array(z.number()).max(5, '最多选择5个标签'),
});

type PostForm = z.infer<typeof postSchema>;

export default function NewPostPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/auth/login');
    }
  }, [isAuthenticated, router]);

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
      tag_ids: [],
    },
  });

  const selectedTagIds = watch('tag_ids');

  const toggleTag = (tagId: number) => {
    const current = selectedTagIds ?? [];
    if (current.includes(tagId)) {
      setValue('tag_ids', current.filter((id) => id !== tagId));
    } else if (current.length < 5) {
      setValue('tag_ids', [...current, tagId]);
    } else {
      toast.error('最多选择5个标签');
    }
  };

  const onSubmit = async (data: PostForm) => {
    setLoading(true);
    try {
      const res = await postApi.create({
        ...data,
        cover: data.cover || undefined,
        summary: data.summary || undefined,
      });
      toast.success('发布成功！');
      router.push(`/posts/${res.data.data.id}`);
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
        <h1 className="text-2xl font-bold">发布新帖子</h1>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-5 space-y-4">
            {/* Type */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">帖子类型</span>
              </label>
              <div className="flex gap-2">
                {[
                  { value: 'post', label: '帖子', desc: '普通讨论帖' },
                  { value: 'article', label: '文章', desc: '长篇技术文章' },
                  { value: 'topic', label: '话题', desc: '热门话题讨论' },
                ].map((t) => (
                  <label key={t.value} className="flex-1 cursor-pointer">
                    <input {...register('type')} type="radio" value={t.value} className="hidden peer" />
                    <div className="border-2 border-base-300 rounded-xl p-3 text-center peer-checked:border-primary peer-checked:bg-primary/5 transition-all">
                      <div className="font-medium text-sm">{t.label}</div>
                      <div className="text-xs text-base-content/40 mt-0.5">{t.desc}</div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Title */}
            <div className="form-control">
              <label className="label pb-1">
                <span className="label-text font-medium">标题 <span className="text-error">*</span></span>
              </label>
              <input
                {...register('title')}
                type="text"
                placeholder="请输入帖子标题"
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
                <span className="label-text font-medium">标签 <span className="text-base-content/40 text-xs">(最多5个)</span></span>
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
                        // ringColor: tag.color,
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
                <span className="label-text font-medium">封面图片 <span className="text-base-content/40 text-xs">(可选，填写图片URL)</span></span>
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
                <span className="label-text font-medium">摘要 <span className="text-base-content/40 text-xs">(可选，显示在列表页)</span></span>
              </label>
              <textarea
                {...register('summary')}
                rows={2}
                placeholder="简短描述帖子内容..."
                className="textarea textarea-bordered focus:outline-none focus:border-primary resize-none"
              />
            </div>
          </div>
        </div>

        {/* Content editor */}
        <div>
          <label className="label pb-2">
            <span className="label-text font-medium text-base">正文内容 <span className="text-error">*</span></span>
          </label>
          <Controller
            name="content"
            control={control}
            render={({ field }) => (
              <RichEditor
                content={field.value}
                onChange={field.onChange}
                placeholder="开始写作..."
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
            取消
          </button>
          <button type="submit" className="btn btn-primary gap-2" disabled={loading}>
            {loading ? (
              <span className="loading loading-spinner loading-sm" />
            ) : (
              <Send className="w-4 h-4" />
            )}
            发布帖子
          </button>
        </div>
      </form>
    </div>
  );
}
