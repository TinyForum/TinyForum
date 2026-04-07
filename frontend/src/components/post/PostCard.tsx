'use client';

import Link from 'next/link';
import Image from 'next/image';
import { Post } from '@/types';
import { timeAgo, truncate } from '@/lib/utils';
import { Eye, Heart, MessageSquare, Pin, Tag } from 'lucide-react';

interface PostCardProps {
  post: Post;
  commentCount?: number;
}

export default function PostCard({ post, commentCount }: PostCardProps) {
  return (
    <div className={`card bg-base-100 shadow-sm hover:shadow-md transition-all duration-200 hover:-translate-y-0.5 border border-base-300 ${post.pin_top ? 'border-l-4 border-l-primary' : ''}`}>
      <div className="card-body p-4">
        {/* Top row */}
        <div className="flex items-start gap-3">
          {/* Author avatar */}
          <Link href={`/users/${post.author?.id}`} className="flex-none">
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

          <div className="flex-1 min-w-0">
            {/* Title */}
            <div className="flex items-center gap-2 flex-wrap">
              {post.pin_top && (
                <span className="badge badge-primary badge-sm gap-1">
                  <Pin className="w-3 h-3" /> 置顶
                </span>
              )}
              <span className={`badge badge-sm ${
                post.type === 'article' ? 'badge-secondary' :
                post.type === 'topic' ? 'badge-accent' : 'badge-ghost'
              }`}>
                {post.type === 'article' ? '文章' : post.type === 'topic' ? '话题' : '帖子'}
              </span>
            </div>

            <Link href={`/posts/${post.id}`} className="group">
              <h2 className="text-base font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-2 mt-1">
                {post.title}
              </h2>
            </Link>

            {post.summary && (
              <p className="text-sm text-base-content/60 mt-1 line-clamp-2">
                {truncate(post.summary, 120)}
              </p>
            )}
          </div>

          {/* Cover image */}
          {post.cover && (
            <Link href={`/posts/${post.id}`} className="flex-none hidden sm:block">
              <div className="w-20 h-16 rounded-lg overflow-hidden">
                <Image
                  src={post.cover}
                  alt={post.title}
                  width={80}
                  height={64}
                  className="object-cover w-full h-full"
                />
              </div>
            </Link>
          )}
        </div>

        {/* Bottom row */}
        <div className="flex items-center justify-between mt-3 flex-wrap gap-2">
          <div className="flex items-center gap-3 text-xs text-base-content/50">
            <Link href={`/users/${post.author?.id}`} className="hover:text-primary transition-colors font-medium">
              {post.author?.username}
            </Link>
            <span>{timeAgo(post.created_at)}</span>
            {post.tags && post.tags.length > 0 && (
              <div className="flex items-center gap-1">
                <Tag className="w-3 h-3" />
                {post.tags.slice(0, 2).map((tag) => (
                  <Link
                    key={tag.id}
                    href={`/posts?tag_id=${tag.id}`}
                    className="badge badge-sm"
                    style={{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }}
                  >
                    {tag.name}
                  </Link>
                ))}
              </div>
            )}
          </div>

          <div className="flex items-center gap-3 text-xs text-base-content/50">
            <span className="flex items-center gap-1">
              <Eye className="w-3.5 h-3.5" /> {post.view_count}
            </span>
            <span className="flex items-center gap-1">
              <Heart className="w-3.5 h-3.5" /> {post.like_count}
            </span>
            {commentCount !== undefined && (
              <span className="flex items-center gap-1">
                <MessageSquare className="w-3.5 h-3.5" /> {commentCount}
              </span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
