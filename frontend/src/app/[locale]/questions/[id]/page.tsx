'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { questionsApi } from '@/lib/api/questions';
import { useAuthStore } from '@/store/auth';
import { toast } from 'react-hot-toast';
import {
  ArrowLeftIcon,
  EyeIcon,
  ChatBubbleLeftRightIcon,
  CheckBadgeIcon,
  UserCircleIcon,
  CalendarIcon,
  TagIcon,
  HeartIcon,
  ShareIcon,
  FlagIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import type { Post, Question, Comment } from '@/types';
import { ArrowDownIcon, ArrowUpIcon } from 'lucide-react';

// 回答卡片组件
function AnswerCard({ 
  answer, 
  isAccepted, 
  isAuthor, 
  canAccept, 
  onAccept,
  currentUserId,
}: { 
  answer: Comment; 
  isAccepted: boolean; 
  isAuthor: boolean; 
  canAccept: boolean; 
  onAccept: () => void;
  currentUserId?: number;
}) {
  const [userVote, setUserVote] = useState<string>('');
  const [voteCount, setVoteCount] = useState((answer as any).vote_count || 0);
  const [voting, setVoting] = useState(false);

  useEffect(() => {
    if (currentUserId) {
      loadVoteStatus();
    }
  }, [currentUserId]);

  const loadVoteStatus = async () => {
    try {
      const response = await questionsApi.getVoteStatus(answer.id);
      if (response.data.code === 200) {
        setUserVote(response.data.data.vote_type);
      }
    } catch (error) {
      console.error('Failed to load vote status:', error);
    }
  };

  const handleVote = async (voteType: 'up' | 'down') => {
    if (!currentUserId) {
      toast.error('请先登录');
      return;
    }

    if (answer.author_id === currentUserId) {
      toast.error('不能给自己的答案投票');
      return;
    }

    if (voting) return;

    const wasVoted = userVote === voteType;
    const newVoteType = wasVoted ? '' : voteType;
    const delta = wasVoted
      ? (voteType === 'up' ? -1 : 1)
      : (voteType === 'up' ? 1 : -1);
    
    setVoting(true);
    setUserVote(newVoteType);
    // FIXME: 修复投票计数问题
    setVoteCount(prev => prev + delta);
    
    try {
      await questionsApi.voteAnswer(answer.id, voteType);
    } catch (error) {
      setUserVote(userVote);
      setVoteCount(prev => prev - delta);
      toast.error('投票失败');
    } finally {
      setVoting(false);
    }
  };

  return (
    <div className={`bg-white rounded-lg shadow-sm p-6 transition-all ${
      isAccepted ? 'border-2 border-green-500 shadow-md' : ''
    }`}>
      <div className="flex gap-4">
        {/* 投票区域 */}
        <div className="flex flex-col items-center gap-1">
          <button
            onClick={() => handleVote('up')}
            disabled={voting}
            className={`p-1 rounded hover:bg-gray-100 transition-colors ${
              userVote === 'up' ? 'text-orange-500' : 'text-gray-400'
            }`}
          >
            <ArrowUpIcon className="w-6 h-6" />
          </button>
          <span className="font-medium text-gray-700">{voteCount}</span>
          <button
            onClick={() => handleVote('down')}
            disabled={voting}
            className={`p-1 rounded hover:bg-gray-100 transition-colors ${
              userVote === 'down' ? 'text-blue-500' : 'text-gray-400'
            }`}
          >
            <ArrowDownIcon className="w-6 h-6" />
          </button>
        </div>

        {/* 内容区域 */}
        <div className="flex-1">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-3 text-sm text-gray-500">
              <div className="flex items-center gap-1">
                <UserCircleIcon className="w-4 h-4" />
                <Link href={`/users/${answer.author_id}`} className="hover:text-indigo-600">
                  {answer.author?.username}
                </Link>
              </div>
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-4 h-4" />
                {new Date(answer.created_at).toLocaleDateString()}
              </div>
            </div>
            {isAccepted && (
              <div className="flex items-center gap-1 text-green-500 bg-green-50 px-2 py-1 rounded">
                <CheckBadgeIcon className="w-4 h-4" />
                <span className="text-sm font-medium">已采纳</span>
              </div>
            )}
          </div>

          <div
            className="prose max-w-none text-gray-700"
            dangerouslySetInnerHTML={{ __html: answer.content }}
          />

          {/* 采纳按钮 */}
          {canAccept && (
            <div className="mt-4 pt-3 border-t">
              <button
                onClick={onAccept}
                className="flex items-center gap-1 text-green-600 hover:text-green-700 text-sm font-medium transition-colors"
              >
                <CheckBadgeIcon className="w-4 h-4" />
                采纳为答案
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// 回答表单组件
function AnswerForm({ questionId, onSuccess }: { questionId: number; onSuccess: () => void }) {
  const [content, setContent] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!content.trim()) {
      toast.error('请输入回答内容');
      return;
    }

    setSubmitting(true);
    try {
      const response = await questionsApi.createAnswer(questionId, { content });
      if (response.data.code === 200) {
        setContent('');
        onSuccess();
      } else {
        toast.error(response.data.message || '发布失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '发布失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">你的回答</h3>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={6}
        placeholder="写下你的回答...\n\n建议：\n1. 直接回答问题\n2. 提供代码示例\n3. 给出具体步骤"
        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
      />
      <div className="flex justify-end mt-4">
        <button
          onClick={handleSubmit}
          disabled={submitting}
          className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {submitting ? '发布中...' : '发布回答'}
        </button>
      </div>
    </div>
  );
}

// 主页面组件
export default function QuestionDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const [question, setQuestion] = useState<Post | null>(null);
  const [questionInfo, setQuestionInfo] = useState<Question | null>(null);
  const [answers, setAnswers] = useState<Comment[]>([]);
  const [answersTotal, setAnswersTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [answerPage, setAnswerPage] = useState(1);
  const [acceptedAnswerId, setAcceptedAnswerId] = useState<number | null>(null);
  const [liked, setLiked] = useState(false);
  const [likesCount, setLikesCount] = useState(0);

  const id = Number(params.id);
  const pageSize = 20;

  useEffect(() => {
    loadQuestion();
  }, [id, answerPage]);

  const loadQuestion = async () => {
    setLoading(true);
    try {
      const response = await questionsApi.getDetail(id, {
        answer_page: answerPage,
        answer_page_size: pageSize,
      });
      if (response.data.code === 200) {
        const data = response.data.data;
        setQuestion(data.post);
        setQuestionInfo(data.question);
        setAnswers(data.answers);
        setAnswersTotal(data.answers_total);
        setAcceptedAnswerId(data.question?.accepted_answer_id || null);
        setLiked(data.liked);
        setLikesCount(data.post.like_count);
      }
    } catch (error) {
      console.error('Failed to load question:', error);
      toast.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAcceptAnswer = async (answerId: number) => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      router.push('/login');
      return;
    }

    try {
      const response = await questionsApi.acceptAnswer(id, answerId);
      if (response.data.code === 200) {
        toast.success('已采纳答案');
        setAcceptedAnswerId(answerId);
        await loadQuestion();
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleLike = async () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      router.push('/login');
      return;
    }
    // TODO: 调用点赞 API
    setLiked(!liked);
    setLikesCount(prev => liked ? prev - 1 : prev + 1);
  };

  const onAnswerCreated = () => {
    setAnswerPage(1);
    loadQuestion();
    toast.success('回答已发布');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-500">加载中...</div>
      </div>
    );
  }

  if (!question) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-500 mb-4">问题不存在</p>
          <Link href="/questions" className="text-indigo-600 hover:underline">
            返回问答列表
          </Link>
        </div>
      </div>
    );
  }

  const isAuthor = user?.id === question.author_id;
  const hasAccepted = !!acceptedAnswerId;
  const totalPages = Math.ceil(answersTotal / pageSize);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 返回按钮 */}
        <div className="flex items-center justify-between mb-4">
          <Link
            href="/questions"
            className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            返回列表
          </Link>
          <div className="flex gap-2">
            <button className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100">
              <ShareIcon className="w-5 h-5" />
            </button>
            <button className="p-2 text-gray-400 hover:text-red-500 rounded-lg hover:bg-gray-100">
              <FlagIcon className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* 问题卡片 */}
        <div className="bg-white rounded-lg shadow-sm mb-6">
          <div className="p-6">
            <h1 className="text-2xl font-bold text-gray-900 mb-4">{question.title}</h1>
            
            {/* 元信息 */}
            <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500 mb-4 pb-4 border-b">
              <div className="flex items-center gap-1">
                <UserCircleIcon className="w-4 h-4" />
                <Link href={`/users/${question.author_id}`} className="hover:text-indigo-600">
                  {question.author?.username}
                </Link>
              </div>
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-4 h-4" />
                {new Date(question.created_at).toLocaleDateString()}
              </div>
              <div className="flex items-center gap-1">
                <EyeIcon className="w-4 h-4" />
                {question.view_count} 浏览
              </div>
              <div className="flex items-center gap-1">
                <ChatBubbleLeftRightIcon className="w-4 h-4" />
                {answersTotal} 回答
              </div>
              <button
                onClick={handleLike}
                className="flex items-center gap-1 hover:text-red-500 transition-colors"
              >
                {liked ? (
                  <HeartSolidIcon className="w-4 h-4 text-red-500" />
                ) : (
                  <HeartIcon className="w-4 h-4" />
                )}
                {likesCount} 点赞
              </button>
              {questionInfo?.reward_score && questionInfo.reward_score > 0 && (
                <div className="flex items-center gap-1 text-orange-500 bg-orange-50 px-2 py-1 rounded">
                  💰 {questionInfo.reward_score} 积分悬赏
                </div>
              )}
              {hasAccepted && (
                <div className="flex items-center gap-1 text-green-500 bg-green-50 px-2 py-1 rounded">
                  <CheckBadgeIcon className="w-4 h-4" />
                  已解决
                </div>
              )}
            </div>

            {/* 标签 */}
            {question.tags && question.tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-4">
                {question.tags.map((tag) => (
                  <Link
                    key={tag.id}
                    href={`/questions?tag_id=${tag.id}`}
                    className="px-2 py-1 bg-gray-100 text-gray-600 text-sm rounded-md hover:bg-gray-200 transition-colors"
                    style={{ borderLeft: `3px solid ${tag.color || '#6366f1'}` }}
                  >
                    <TagIcon className="w-3 h-3 inline mr-1" />
                    {tag.name}
                  </Link>
                ))}
              </div>
            )}

            {/* 内容 */}
            <div
              className="prose max-w-none"
              dangerouslySetInnerHTML={{ __html: question.content }}
            />
          </div>
        </div>

        {/* 答案列表 */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900">
              {answersTotal} 个回答
            </h2>
            {answersTotal > 0 && (
              <select className="text-sm border rounded-md px-2 py-1">
                <option value="vote">按投票排序</option>
                <option value="newest">按最新排序</option>
                <option value="oldest">按最早排序</option>
              </select>
            )}
          </div>
          
          {answers.length === 0 ? (
            <div className="bg-white rounded-lg shadow-sm p-8 text-center text-gray-500">
              <ChatBubbleLeftRightIcon className="w-12 h-12 mx-auto text-gray-300 mb-3" />
              <p>暂无回答</p>
              <p className="text-sm mt-1">成为第一个回答的人吧！</p>
            </div>
          ) : (
            <div className="space-y-4">
              {answers.map((answer) => (
                <AnswerCard
                  key={answer.id}
                  answer={answer}
                  isAccepted={answer.id === acceptedAnswerId}
                  isAuthor={isAuthor}
                  canAccept={!hasAccepted && isAuthor && answer.id !== acceptedAnswerId}
                  onAccept={() => handleAcceptAnswer(answer.id)}
                  currentUserId={user?.id}
                />
              ))}
            </div>
          )}

          {/* 分页 */}
          {totalPages > 1 && (
            <div className="flex justify-center gap-2 mt-6">
              <button
                onClick={() => setAnswerPage(p => Math.max(1, p - 1))}
                disabled={answerPage === 1}
                className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
              >
                上一页
              </button>
              <span className="px-3 py-1 text-gray-600">
                第 {answerPage} / {totalPages} 页
              </span>
              <button
                onClick={() => setAnswerPage(p => Math.min(totalPages, p + 1))}
                disabled={answerPage >= totalPages}
                className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
              >
                下一页
              </button>
            </div>
          )}
        </div>

        {/* 回答表单 */}
        {isAuthenticated ? (
          <AnswerForm questionId={id} onSuccess={onAnswerCreated} />
        ) : (
          <div className="bg-white rounded-lg shadow-sm p-6 text-center">
            <p className="text-gray-500 mb-3">登录后回答这个问题</p>
            <Link
              href={`/login?redirect=/questions/${id}`}
              className="inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
            >
              登录
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}