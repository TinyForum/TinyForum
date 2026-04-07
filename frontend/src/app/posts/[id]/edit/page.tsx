'use client';

import { use, useEffect } from 'react';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { postApi, tagApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import RichEditor from '@/components/post/RichEditor';
import toast from 'react-hot-toast';
import { getErrorMessage } from '@/lib/utils';
import { Save, X } from 'lucide-react';

const schema = z.object({
  title: z.string().min(2).max(200),
  content: z.string().min(10),
  summary: z.string().max(500).optional(),
  cover: z.string().url().optional().or(z.literal('')),
  tag_ids: z.array(z.number()),
});

type EditForm = z.infer<typeof schema>;

export default function EditPostPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const postId = Number(id);
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);

  const { data: postData } = useQuery({
    queryKey: ['post', postId],
    queryFn: () => postApi.getById(postId).then((r) => r.data.data),
  });

  const { data: tags } = useQuery({
    queryKey: ['tags'],
    queryFn: () => tagApi.list().then((r) => r.data.data),
  });

  const { register, handleSubmit, control, watch, setValue, reset, formState: { errors } } = useForm<EditForm>({
    resolver: zodResolver(schema),
    defaultValues: { tag_ids: [] },
  });

  useEffect(() => {
    if (postData?.post) {
      const p = postData.post;
      reset({
        title: p.title,
        content: p.content,
        summary: p.summary || '',
        cover: p.cover || '',
        tag_ids: p.tags?.map((t) => t.id) ?? [],
      });
    }
  }, [postData, reset]);

  useEffect(() => {
    if (!isAuthenticated) router.push('/auth/login');
  }, [isAuthenticated, router]);

  const selectedTagIds = watch('tag_ids');

  const toggleTag = (tagId: number) => {
    const current = selectedTagIds ?? [];
    setValue(
      'tag_ids',
      current.includes(tagId) ? current.filter((id) => id !== tagId) : [...current, tagId]
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
      toast.success('更新成功');
      router.push(`/posts/${postId}`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (!postData) {
    return <div className="flex justify-center py-20"><span className="loading loading-spinner loading-lg text-primary" /></div>;
  }

  return (
    <div className="max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">编辑帖子</h1>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-5 space-y-4">
            <div className="form-control">
              <label className="label pb-1"><span className="label-text font-medium">标题</span></label>
              <input {...register('title')} className={`input input-bordered focus:outline-none focus:border-primary ${errors.title ? 'input-error' : ''}`} />
              {errors.title && <span className="text-error text-sm mt-1">{errors.title.message}</span>}
            </div>

            <div className="form-control">
              <label className="label pb-1"><span className="label-text font-medium">标签</span></label>
              <div className="flex flex-wrap gap-2">
                {(tags ?? []).map((tag) => {
                  const selected = selectedTagIds?.includes(tag.id);
                  return (
                    <button key={tag.id} type="button" onClick={() => toggleTag(tag.id)}
                      className={`badge badge-lg cursor-pointer transition-all ${selected ? 'ring-2' : 'opacity-60'}`}
                      style={{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }}>
                      {selected && <X className="w-3 h-3 mr-1" />}{tag.name}
                    </button>
                  );
                })}
              </div>
            </div>

            <div className="form-control">
              <label className="label pb-1"><span className="label-text font-medium">封面图片URL</span></label>
              <input {...register('cover')} type="text" className="input input-bordered focus:outline-none focus:border-primary" />
            </div>

            <div className="form-control">
              <label className="label pb-1"><span className="label-text font-medium">摘要</span></label>
              <textarea {...register('summary')} rows={2} className="textarea textarea-bordered focus:outline-none focus:border-primary resize-none" />
            </div>
          </div>
        </div>

        <div>
          <label className="label pb-2"><span className="label-text font-medium text-base">正文内容</span></label>
          <Controller
            name="content"
            control={control}
            render={({ field }) => (
              <RichEditor content={field.value} onChange={field.onChange} />
            )}
          />
          {errors.content && <p className="text-error text-sm mt-1">{errors.content.message}</p>}
        </div>

        <div className="flex gap-3 justify-end">
          <button type="button" className="btn btn-ghost" onClick={() => router.back()}>取消</button>
          <button type="submit" className="btn btn-primary gap-2" disabled={loading}>
            {loading ? <span className="loading loading-spinner loading-sm" /> : <Save className="w-4 h-4" />}
            保存修改
          </button>
        </div>
      </form>
    </div>
  );
}
