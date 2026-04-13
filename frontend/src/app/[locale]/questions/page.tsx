'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import {
  MagnifyingGlassIcon,
  ChatBubbleLeftRightIcon,
  CheckBadgeIcon,
  FireIcon,
  ArrowsRightLeftIcon,
  PlusIcon,
  UserCircleIcon,
  EyeIcon,
  TagIcon,
  SparklesIcon,
} from '@heroicons/react/24/outline';
import { useAuthStore } from '@/store/auth';
import { postApi } from '@/lib/api';
import { Post } from '@/lib/api/types';

type FilterType = 'all' | 'unanswered' | 'answered';
type SortType = 'latest' | 'hot' | 'score';

export default function QuestionsPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { isAuthenticated } = useAuthStore();

  const [questions, setQuestions] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [filter, setFilter] = useState<FilterType>((searchParams.get('filter') as FilterType) || 'all');
  const [sort, setSort] = useState<SortType>('latest');
  const [keyword, setKeyword] = useState(searchParams.get('keyword') || '');

  const pageSize = 15;

  useEffect(() => {
    loadQuestions();

  },[]);

  useEffect(() => {
    loadQuestions();
  }, [page, filter, sort, keyword]);

  const loadQuestions = async () => {
    setLoading(true);
    try {
      const response = await postApi.getQuestions({
        page,
        page_size: pageSize,
        filter: filter === 'all' ? undefined : filter,
        keyword: keyword || undefined,
      });
      
      
      if (response.status === 200) {
      // console.log(response.data.data.list);
        setQuestions(response.data.data.list);
        setTotal(response.data.data.total);
      }
    } catch (error) {
      console.error('Failed to load questions:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    loadQuestions();
  };

  const formatTime = (time: string) => {
    const date = new Date(time);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    if (days === 0) return '今天';
    if (days === 1) return '昨天';
    if (days < 7) return `${days}天前`;
    return date.toLocaleDateString();
  };

  const getAnswerCount = (question: Post) => {
    return (question as any).answer_count || (question as any).answers_total || 0;
  };

  const getRewardScore = (question: Post) => {
    return (question as any).reward_score || 0;
  };

  const getIsAccepted = (question: Post) => {
    return !!(question as any).accepted_answer_id;
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Hero Section */}
        <div className="text-center mb-10">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-red-500 to-red-600 shadow-lg mb-4">
            <SparklesIcon className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-red-600 to-red-500 bg-clip-text text-transparent">
            问答社区
          </h1>
          <p className="text-base-content/60 mt-2">
            提问、回答、分享知识
          </p>
        </div>

        {/* Search Bar */}
        <div className="mb-8">
          <div className="card bg-base-100 shadow-md border border-base-200 p-2">
            <form onSubmit={handleSearch} className="flex gap-2">
              <div className="flex-1 relative">
                <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                <input
                  type="text"
                  value={keyword}
                  onChange={(e) => setKeyword(e.target.value)}
                  placeholder="搜索问题..."
                  className="w-full pl-11 pr-4 py-2.5 bg-transparent text-base-content placeholder-base-content/40 focus:outline-none"
                />
              </div>
              <button
                type="submit"
                className="btn btn-primary min-w-[80px]"
              >
                搜索
              </button>
            </form>
          </div>
        </div>

        {/* Main Content */}
        <div className="card bg-base-100 shadow-md border border-base-200">
          {/* Header with Filters */}
          <div className="card-body p-0">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-4 border-b border-base-200">
              <div className="flex flex-wrap gap-2">
                {[
                  { value: 'all', label: '全部', icon: ChatBubbleLeftRightIcon, color: 'text-blue-500' },
                  { value: 'unanswered', label: '未回答', icon: FireIcon, color: 'text-orange-500' },
                  { value: 'answered', label: '已解决', icon: CheckBadgeIcon, color: 'text-green-500' },
                ].map((item) => {
                  const Icon = item.icon;
                  const isActive = filter === item.value;
                  return (
                    <button
                      key={item.value}
                      onClick={() => {
                        setFilter(item.value as FilterType);
                        setPage(1);
                      }}
                      className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all ${
                        isActive
                          ? 'bg-primary/10 text-primary border border-primary/20'
                          : 'text-base-content/60 hover:text-base-content hover:bg-base-200/50'
                      }`}
                    >
                      <Icon className={`w-4 h-4 ${isActive ? item.color : ''}`} />
                      {item.label}
                    </button>
                  );
                })}
              </div>
              
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2 text-sm">
                  <ArrowsRightLeftIcon className="w-4 h-4 text-base-content/40" />
                  <select
                    value={sort}
                    onChange={(e) => setSort(e.target.value as SortType)}
                    className="select select-bordered select-sm bg-transparent focus:outline-none"
                  >
                    <option value="latest">最新</option>
                    <option value="hot">最热</option>
                    <option value="score">悬赏最高</option>
                  </select>
                </div>
                
                {isAuthenticated && (
                  <Link
                    href="/questions/ask"
                    className="btn btn-primary btn-sm"
                  >
                    <PlusIcon className="w-4 h-4" />
                    提问
                  </Link>
                )}
              </div>
            </div>

            {/* Questions List */}
            <div className="divide-y divide-base-200">
              {loading ? (
                <div className="space-y-3 p-4">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <div key={i} className="card bg-base-100 border border-base-200">
                      <div className="card-body p-5">
                        <div className="h-5 bg-base-200 rounded w-3/4 mb-2 animate-pulse" />
                        <div className="h-4 bg-base-200 rounded w-1/2 mb-3 animate-pulse" />
                        <div className="flex gap-4">
                          <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                          <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : questions.length === 0 ? (
                <div className="text-center py-16">
                  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-50 dark:bg-red-900/20 flex items-center justify-center">
                    <ChatBubbleLeftRightIcon className="w-10 h-10 text-red-400" />
                  </div>
                  <h3 className="text-lg font-semibold text-base-content mb-2">暂无问题</h3>
                  <p className="text-base-content/60 mb-4">还没有人提问，成为第一个提问者吧</p>
                  {isAuthenticated && (
                    <Link href="/questions/ask" className="btn btn-primary">
                      <PlusIcon className="w-4 h-4" />
                      我要提问
                    </Link>
                  )}
                </div>
              ) : (
                questions.map((question, index) => (
                  <Link
                    key={question.id}
                    href={`/questions/${question.id}`}
                    className="block group hover:bg-base-200/50 transition-colors"
                  >
                    <div className="p-5">
                      <div className="flex items-start gap-4">
                        {/* 左侧统计信息 */}
                        <div className="shrink-0 text-center min-w-[60px]">
                          <div className="text-sm font-semibold text-base-content">
                            {getAnswerCount(question)}
                          </div>
                          <div className="text-xs text-base-content/40">回答</div>
                          <div className="text-sm font-semibold text-base-content mt-1">
                            {question.view_count || 0}
                          </div>
                          <div className="text-xs text-base-content/40">浏览</div>
                        </div>
                        
                        {/* 问题内容 */}
                        <div className="flex-1 min-w-0">
                          <div className="flex items-start justify-between gap-4">
                            <h3 className="text-base font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-1">
                              {question.title}
                            </h3>
                            {getIsAccepted(question) && (
                              <span className="badge badge-success gap-1 shrink-0">
                                <CheckBadgeIcon className="w-3 h-3" />
                                已解决
                              </span>
                            )}
                          </div>
                          
                          {getRewardScore(question) > 0 && (
                            <div className="flex items-center gap-1 mt-1">
                              <span className="badge badge-warning badge-sm gap-1">
                                💰 {getRewardScore(question)} 积分悬赏
                              </span>
                            </div>
                          )}
                          
                          <p className="text-sm text-base-content/60 line-clamp-2 mt-2">
                            {question.summary || question.content?.replace(/<[^>]*>/g, '').slice(0, 150) || '暂无内容'}
                          </p>
                          
                          <div className="flex flex-wrap items-center gap-4 mt-3 text-xs text-base-content/50">
                            <span className="flex items-center gap-1">
                              <UserCircleIcon className="w-3.5 h-3.5" />
                              {question.author?.username || '匿名用户'}
                            </span>
                            <span>{formatTime(question.created_at)}</span>
                            
                            {/* Tags */}
                            {question.tags && question.tags.length > 0 && (
                              <div className="flex flex-wrap gap-1">
                                {question.tags.slice(0, 3).map((tag) => (
                                  <span
                                    key={tag.id}
                                    className="px-2 py-0.5 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-xs rounded-md inline-flex items-center gap-0.5"
                                  >
                                    <TagIcon className="w-3 h-3" />
                                    {tag.name}
                                  </span>
                                ))}
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </Link>
                ))
              )}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center gap-2 p-4 border-t border-base-200">
                <div className="join">
                  <button
                    onClick={() => setPage(p => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="join-item btn btn-sm"
                  >
                    «
                  </button>
                  <button className="join-item btn btn-sm btn-active">
                    {page}
                  </button>
                  <button
                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                    disabled={page >= totalPages}
                    className="join-item btn btn-sm"
                  >
                    »
                  </button>
                </div>
                <span className="text-sm text-base-content/60 flex items-center">
                  共 {total} 条
                </span>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}