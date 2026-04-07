'use client';

import { use } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { postApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import CommentSection from '@/components/post/CommentSection';
import Image from 'next/image';
import Link from 'next/link';
import { timeAgo, formatDate } from '@/lib/utils';
import {
  Eye, Heart, HeartOff, Share2, Pencil, Trash2,
  Clock, ArrowLeft, Tag,
} from 'lucide-react';
import toast from 'react-hot-toast';
import { useRouter } from 'next/navigation';

export default function PostDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const postId = Number(id);
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ['post', postId],
    queryFn: () => postApi.getById(postId).then((r) => r.data.data),
  });

  const likeMutation = useMutation({
    mutationFn: () => (data?.liked ? postApi.unlike(postId) : postApi.like(postId)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['post', postId] });
      toast.success(data?.liked ? '已取消点赞' : '点赞成功');
    },
    onError: () => toast.error('操作失败'),
  });

  const deleteMutation = useMutation({
    mutationFn: () => postApi.delete(postId),
    onSuccess: () => {
      toast.success('帖子已删除');
      router.push('/');
    },
    onError: () => toast.error('删除失败'),
  });

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href);
    toast.success('链接已复制到剪贴板');
  };

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto space-y-4">
        <div className="skeleton h-10 w-3/4" />
        <div className="skeleton h-4 w-1/2" />
        <div className="skeleton h-64 w-full" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="text-center py-20">
        <p className="text-base-content/40 mb-4">帖子不存在或已被删除</p>
        <Link href="/" className="btn btn-primary">返回首页</Link>
      </div>
    );
  }

  const { post, liked } = data;
  const isAuthor = user?.id === post.author_id;
  const isAdmin = user?.role === 'admin';

  return (
    <div className="max-w-3xl mx-auto">
      {/* Back button */}
      <button
        onClick={() => router.back()}
        className="btn btn-ghost btn-sm gap-1 mb-4"
      >
        <ArrowLeft className="w-4 h-4" /> 返回
      </button>

      <article className="card bg-base-100 border border-base-300 shadow-sm mb-6">
        <div className="card-body p-6 lg:p-8">
          {/* Post type + tags */}
          <div className="flex items-center flex-wrap gap-2 mb-3">
            <span className={`badge ${
              post.type === 'article' ? 'badge-secondary' :
              post.type === 'topic' ? 'badge-accent' : 'badge-ghost'
            }`}>
              {post.type === 'article' ? '文章' : post.type === 'topic' ? '话题' : '帖子'}
            </span>
            {post.tags?.map((tag) => (
              <Link
                key={tag.id}
                href={`/posts?tag_id=${tag.id}`}
                className="badge badge-sm gap-1"
                style={{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }}
              >
                <Tag className="w-2.5 h-2.5" /> {tag.name}
              </Link>
            ))}
          </div>

          {/* Title */}
          <h1 className="text-2xl lg:text-3xl font-bold text-base-content leading-tight">
            {post.title}
          </h1>

          {/* Author info */}
          <div className="flex items-center gap-3 mt-4 pb-4 border-b border-base-300">
            <Link href={`/users/${post.author?.id}`}>
              <div className="avatar">
                <div className="w-10 h-10 rounded-full">
                  <Image
                    src={post.author?.avatar || `https://api.dicebear.com/8.x/initials/svg?seed=${post.author?.username}`}
                    alt={post.author?.username || ''}
                    width={40}
                    height={40}
                    className="rounded-full"
                  />
                </div>
              </div>
            </Link>
            <div>
              <Link href={`/users/${post.author?.id}`} className="font-medium hover:text-primary transition-colors text-sm">
                {post.author?.username}
              </Link>
              <div className="flex items-center gap-2 text-xs text-base-content/40">
                <Clock className="w-3 h-3" />
                <span title={formatDate(post.created_at)}>{timeAgo(post.created_at)}</span>
                <span>·</span>
                <Eye className="w-3 h-3" />
                <span>{post.view_count} 次阅读</span>
              </div>
            </div>

            {/* Action buttons */}
            <div className="ml-auto flex items-center gap-1">
              {(isAuthor || isAdmin) && (
                <>
                  <Link href={`/posts/${postId}/edit`} className="btn btn-ghost btn-xs gap-1">
                    <Pencil className="w-3.5 h-3.5" /> 编辑
                  </Link>
                  <button
                    className="btn btn-ghost btn-xs text-error gap-1"
                    onClick={() => {
                      if (confirm('确认删除此帖子？')) deleteMutation.mutate();
                    }}
                    disabled={deleteMutation.isPending}
                  >
                    <Trash2 className="w-3.5 h-3.5" /> 删除
                  </button>
                </>
              )}
            </div>
          </div>

          {/* Cover image */}
          {post.cover && (
            <div className="my-4 rounded-xl overflow-hidden">
              <Image
                src={post.cover}
                alt={post.title}
                width={800}
                height={400}
                className="w-full object-cover max-h-72"
              />
            </div>
          )}

          {/* Content */}
          <div
            className="prose-content mt-4 text-base-content/80 leading-relaxed"
            dangerouslySetInnerHTML={{ __html: post.content }}
          />

          {/* Footer actions */}
          <div className="flex items-center gap-3 mt-8 pt-4 border-t border-base-300">
            <button
              className={`btn btn-sm gap-2 ${liked ? 'btn-error' : 'btn-ghost'}`}
              onClick={() => {
                if (!isAuthenticated) {
                  toast.error('请先登录');
                  return;
                }
                likeMutation.mutate();
              }}
              disabled={likeMutation.isPending}
            >
              {liked ? <HeartOff className="w-4 h-4" /> : <Heart className="w-4 h-4" />}
              {post.like_count} 点赞
            </button>
            <button className="btn btn-ghost btn-sm gap-2" onClick={handleShare}>
              <Share2 className="w-4 h-4" /> 分享
            </button>
          </div>
        </div>
      </article>

      {/* Comments */}
      <div className="card bg-base-100 border border-base-300 shadow-sm">
        <div className="card-body p-6 lg:p-8">
          <CommentSection postId={postId} />
        </div>
      </div>
    </div>
  );
}
