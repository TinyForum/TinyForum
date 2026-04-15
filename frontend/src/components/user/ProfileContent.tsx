'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { postApi } from '@/lib/api';
import PostCard from '@/components/post/PostCard';
import { useTranslations } from 'next-intl';
import { FileText, MessageSquare, BookOpen, Hash } from 'lucide-react';
import type { PostType } from '@/lib/api/types';
import { ViolationStatus } from '@/components/user/ViolationStatus';

interface ProfileContentProps {
  userId: number;
  isAuthenticated: boolean;
}

const TAB_CONFIG = [
  { key: 'post' as PostType, label: 'the_posts', icon: FileText },
  { key: 'article' as PostType, label: 'the_articles', icon: BookOpen },
  { key: 'question' as PostType, label: 'the_question', icon: MessageSquare },
  { key: 'topic' as PostType, label: 'the_topic', icon: Hash },
];

export function ProfileContent({ userId, isAuthenticated }: ProfileContentProps) {
  const [tab, setTab] = useState<PostType>('post');
  const t = useTranslations("Profile");

  // 获取用户的帖子/文章
  const { data: postsData, isLoading } = useQuery({
    queryKey: ['user-posts', userId, tab],
    queryFn: () =>
      postApi.list({
        author_id: userId,
        type: ['article', 'post', 'question'].includes(tab) ? tab : 'topic',
        page: 1,
        page_size: 20,
      }).then((r) => r.data.data),
  });

  const posts = postsData?.list ?? [];
  const currentTabConfig = TAB_CONFIG.find(t => t.key === tab);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="skeleton h-12 w-full rounded-xl" />
        <div className="skeleton h-40 w-full rounded-xl" />
        <div className="skeleton h-40 w-full rounded-xl" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* 违规状态（已登录时） */}
      {isAuthenticated && <ViolationStatus />}

      {/* Tab 切换栏 */}
      <div className="tabs tabs-boxed bg-base-100 border border-base-300 p-1">
        {TAB_CONFIG.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            className={`tab gap-2 flex-1 sm:flex-initial ${tab === key ? 'tab-active' : ''}`}
            onClick={() => setTab(key)}
          >
            <Icon className="w-4 h-4" />
            <span className="hidden sm:inline">{t(label)}</span>
          </button>
        ))}
      </div>

      {/* 内容列表 */}
      {posts.length === 0 ? (
        <div className="text-center py-16 bg-base-100 rounded-xl border border-base-200">
          {currentTabConfig && <currentTabConfig.icon className="w-12 h-12 mx-auto mb-3 opacity-30" />}
          <p className="text-base-content/40">
            {tab === 'post' ? t("no_post") : 
             tab === 'article' ? t("no_article") :
             tab === 'question' ? t("no_question") : t("no_topic")}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}
    </div>
  );
}