'use client';

import { useState, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { commentApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import CommentItem from './CommentItem';
import toast from 'react-hot-toast';
import { Send, MessageSquare } from 'lucide-react';
import Link from 'next/link';

interface CommentSectionProps {
  postId: number;
}

export default function CommentSection({ postId }: CommentSectionProps) {
  const { isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [content, setContent] = useState('');
  const [replyTo, setReplyTo] = useState<{ id: number; username: string } | null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['comments', postId],
    queryFn: () => commentApi.listByPost(postId, { page: 1, page_size: 50 }).then((r) => r.data.data),
  });

  const createMutation = useMutation({
    mutationFn: (vars: { content: string; parent_id?: number }) =>
      commentApi.create({ post_id: postId, ...vars }),
    onSuccess: () => {
      toast.success('评论成功');
      setContent('');
      setReplyTo(null);
      queryClient.invalidateQueries({ queryKey: ['comments', postId] });
    },
    onError: () => toast.error('评论失败，请稍后重试'),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;
    createMutation.mutate({
      content: content.trim(),
      parent_id: replyTo?.id,
    });
  };

  const handleReply = (parentId: number, username: string) => {
    setReplyTo({ id: parentId, username });
    textareaRef.current?.focus();
  };

  const comments = data?.list ?? [];
  const total = data?.total ?? 0;

  return (
    <div className="mt-8">
      <h3 className="text-lg font-bold flex items-center gap-2 mb-6">
        <MessageSquare className="w-5 h-5 text-primary" />
        评论 <span className="text-base-content/40 font-normal text-base">({total})</span>
      </h3>

      {/* Comment input */}
      {isAuthenticated ? (
        <form onSubmit={handleSubmit} className="mb-8">
          {replyTo && (
            <div className="flex items-center gap-2 mb-2 text-sm text-base-content/60 bg-base-200 px-3 py-1.5 rounded-lg">
              <span>回复 <strong className="text-primary">@{replyTo.username}</strong></span>
              <button
                type="button"
                className="ml-auto text-error hover:text-error/80 text-xs"
                onClick={() => setReplyTo(null)}
              >
                取消
              </button>
            </div>
          )}
          <div className="flex gap-3">
            <textarea
              ref={textareaRef}
              className="textarea textarea-bordered flex-1 resize-none focus:outline-none focus:border-primary"
              rows={3}
              placeholder={replyTo ? `回复 @${replyTo.username}...` : '写下你的评论...'}
              value={content}
              onChange={(e) => setContent(e.target.value)}
              maxLength={2000}
            />
          </div>
          <div className="flex items-center justify-between mt-2">
            <span className="text-xs text-base-content/40">{content.length}/2000</span>
            <button
              type="submit"
              className="btn btn-primary btn-sm gap-1"
              disabled={!content.trim() || createMutation.isPending}
            >
              {createMutation.isPending ? (
                <span className="loading loading-spinner loading-xs" />
              ) : (
                <Send className="w-4 h-4" />
              )}
              发布评论
            </button>
          </div>
        </form>
      ) : (
        <div className="alert mb-6">
          <span>
            请 <Link href="/auth/login" className="link link-primary">登录</Link> 后发表评论
          </span>
        </div>
      )}

      {/* Comments list */}
      {isLoading ? (
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner loading-md text-primary" />
        </div>
      ) : comments.length === 0 ? (
        <div className="text-center py-12 text-base-content/40">
          <MessageSquare className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>还没有评论，来发表第一条吧</p>
        </div>
      ) : (
        <div className="space-y-6">
          {comments.map((comment) => (
            <CommentItem
              key={comment.id}
              comment={comment}
              postId={postId}
              onReply={handleReply}
            />
          ))}
        </div>
      )}
    </div>
  );
}
