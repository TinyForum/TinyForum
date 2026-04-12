import apiClient from '../api-client';
import type { 
  ApiResponse, 
  PageData, 
  Post, 
  Comment, 
  Question,
  AnswerVoteResult 
} from '@/types';

// 问答相关 API
export const questionsApi = {
  // 获取问答列表
  list: (params?: {
    page?: number;
    page_size?: number;
    filter?: 'all' | 'unanswered' | 'answered';
    keyword?: string;
    tag_id?: number;
    board_id?: number;
  }) => apiClient.get<ApiResponse<PageData<Post>>>('/posts/questions', { params }),

  // 获取问答详情
  getDetail: (id: number, params?: { answer_page?: number; answer_page_size?: number }) =>
    apiClient.get<ApiResponse<{
      post: Post;
      liked: boolean;
      question: Question;
      answers: Comment[];
      answers_total: number;
      answer_page: number;
      answer_page_size: number;
    }>>(`/posts/question/${id}`, { params }),

  // 创建问答
  create: (data: {
    title: string;
    content: string;
    summary?: string;
    cover?: string;
    board_id?: number;
    tag_ids?: number[];
    reward_score?: number;
  }) => apiClient.post<ApiResponse<Post>>('/posts/question', data),

  // 创建回答
  createAnswer: (postId: number, data: { content: string }) =>
    apiClient.post<ApiResponse<Comment>>(`/posts/question/${postId}/answer`, data),

  // 采纳答案
  acceptAnswer: (postId: number, commentId: number) =>
    apiClient.post<ApiResponse<null>>(`/posts/question/${postId}/accept`, { comment_id: commentId }),

  // 投票
  voteAnswer: (commentId: number, voteType: 'up' | 'down') =>
    apiClient.post<ApiResponse<AnswerVoteResult>>(`/comments/${commentId}/vote`, { vote_type: voteType }),

  // 获取投票状态
  getVoteStatus: (commentId: number) =>
    apiClient.get<ApiResponse<{ has_voted: boolean; vote_type: string; vote_count: number }>>(
      `/comments/${commentId}/vote`
    ),
};