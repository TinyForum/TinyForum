import apiClient from "../api-client";
import type {
  ApiResponse,
  AuthResult,
  User,
  Post,
  Comment,
  Tag,
  Notification,
  PageData,
  PostType,
  Board,
  Topic,
  TimelineEvent,
  Question,
  AnswerVoteResult,
  Moderator,
  Subscription,
  TopicFollow,
  TopicPost,
  ModeratorApplication,
} from "@/types";

// ─── Auth ────────────────────────────────────────────────────────────────────

export const authApi = {
  register: (data: { username: string; email: string; password: string }) =>
    apiClient.post<ApiResponse<AuthResult>>("/auth/register", data),

  login: (data: { email: string; password: string }) =>
    apiClient.post<ApiResponse<AuthResult>>("/auth/login", data),

  me: () => apiClient.get<ApiResponse<User>>("/auth/me"),
};

// ─── Posts ───────────────────────────────────────────────────────────────────

export const postApi = {
  list: (params: {
    page?: number;
    page_size?: number;
    keyword?: string;
    sort_by?: string;
    type?: PostType;
    author_id?: number;
    tag_id?: number;
    board_id?: number;
  }) => apiClient.get<ApiResponse<PageData<Post>>>("/posts", { params }),

  getById: (id: number) =>
    apiClient.get<ApiResponse<{ post: Post; liked: boolean }>>(`/posts/${id}`),

  create: (data: {
    title: string;
    content: string;
    summary?: string;
    cover?: string;
    type?: PostType;
    board_id?: number;
    tag_ids?: number[];
  }) => apiClient.post<ApiResponse<Post>>("/posts", data),

  update: (
    id: number,
    data: {
      title?: string;
      content?: string;
      summary?: string;
      cover?: string;
      tag_ids?: number[];
    },
  ) => apiClient.put<ApiResponse<Post>>(`/posts/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/posts/${id}`),

  like: (id: number) => apiClient.post<ApiResponse<null>>(`/posts/${id}/like`),

  unlike: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}/like`),

  // 问答相关
  getQuestions: (params?: {
    page?: number;
    page_size?: number;
    filter?: "all" | "unanswered" | "answered";
  }) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/posts/questions", { params }),

  createQuestion: (data: {
    title: string;
    content: string;
    summary?: string;
    cover?: string;
    board_id?: number;
    tag_ids?: number[];
    reward_score?: number;
  }) => apiClient.post<ApiResponse<Post>>("/posts/question", data),

  getQuestionDetail: (
    id: number,
    params?: { answer_page?: number; answer_page_size?: number },
  ) =>
    apiClient.get<
      ApiResponse<{
        post: Post;
        liked: boolean;
        question: Question;
        answers: Comment[];
        answers_total: number;
        answer_page: number;
        answer_page_size: number;
      }>
    >(`/posts/question/${id}`, { params }),

  acceptAnswer: (id: number, data: { comment_id: number }) =>
    apiClient.post<ApiResponse<null>>(`/posts/question/${id}/accept`, data),

  createAnswer: (id: number, data: { content: string }) =>
    apiClient.post<ApiResponse<Comment>>(`/posts/question/${id}/answer`, data),
};

// ─── Comments ────────────────────────────────────────────────────────────────

export const commentApi = {
  listByPost: (
    postId: number,
    params?: { page?: number; page_size?: number },
  ) =>
    apiClient.get<ApiResponse<PageData<Comment>>>(`/comments/post/${postId}`, {
      params,
    }),

  create: (data: { post_id: number; content: string; parent_id?: number }) =>
    apiClient.post<ApiResponse<Comment>>("/comments", data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/comments/${id}`),

  // 答案投票
  voteAnswer: (id: number, data: { vote_type: "up" | "down" }) =>
    apiClient.post<ApiResponse<AnswerVoteResult>>(`/comments/${id}/vote`, data),

  getAnswerVoteStatus: (id: number) =>
    apiClient.get<
      ApiResponse<{
        has_voted: boolean;
        vote_type: "up" | "down" | "";
        vote_count: number;
      }>
    >(`/comments/${id}/vote`),

  markAsAnswer: (id: number, data: { is_answer: boolean }) =>
    apiClient.put<ApiResponse<null>>(`/comments/${id}/answer`, data),

  getAnswers: (
    postId: number,
    params?: {
      page?: number;
      page_size?: number;
      sort?: "vote" | "newest" | "oldest";
    },
  ) =>
    apiClient.get<ApiResponse<PageData<Comment>>>(
      `/comments/post/${postId}/answers`,
      { params },
    ),

  acceptAnswer: (id: number, postId: number) =>
    apiClient.post<ApiResponse<null>>(
      `/comments/${id}/accept?post_id=${postId}`,
    ),
};

// ─── Tags ────────────────────────────────────────────────────────────────────

export const tagApi = {
  list: () => apiClient.get<ApiResponse<Tag[]>>("/tags"),

  create: (data: { name: string; description?: string; color?: string }) =>
    apiClient.post<ApiResponse<Tag>>("/tags", data),

  update: (
    id: number,
    data: { name?: string; description?: string; color?: string },
  ) => apiClient.put<ApiResponse<Tag>>(`/tags/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/tags/${id}`),
};

// ─── Users ───────────────────────────────────────────────────────────────────

export const userApi = {
  getProfile: (id: number) => apiClient.get<ApiResponse<User>>(`/users/${id}`),

  updateProfile: (data: { bio?: string; avatar?: string }) =>
    apiClient.put<ApiResponse<User>>("/users/profile", data),

  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/users/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/users/${id}/follow`),

  leaderboard: (limit?: number) =>
    apiClient.get<ApiResponse<User[]>>("/users/leaderboard", {
      params: { limit },
    }),
};

// ─── Notifications ───────────────────────────────────────────────────────────

export const notificationApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Notification>>>("/notifications", {
      params,
    }),

  unreadCount: () =>
    apiClient.get<ApiResponse<{ count: number }>>(
      "/notifications/unread-count",
    ),

  markAllRead: () =>
    apiClient.post<ApiResponse<null>>("/notifications/read-all"),
};

// ─── Boards (板块) ───────────────────────────────────────────────────────────

export const boardApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Board>>>("/boards", { params }),

  getTree: () => apiClient.get<ApiResponse<Board[]>>("/boards/tree"),

  getByName: (id: string) => apiClient.get<ApiResponse<Board>>(`/boards/${id}`),

  getPosts: (id: string, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Post>>>(`/boards/${id}/posts`, {
      params,
    }),
  // 获取申请状态
  applyForModerator: (id: number, reason: { reason: string }) =>
    apiClient.post<ApiResponse<null>>(`/boards/${id}/moderator-apply`, reason),
  // 申请成为版主
  checkApplicationStatus: (boardId: number) =>
    apiClient.get<
      ApiResponse<{ has_applied: boolean; application?: ModeratorApplication }>
    >(`/boards/${boardId}/application-status`),
  // 获取我的申请
  getMyApplications: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<ModeratorApplication>>>(
      "/boards/my-applications",
      { params },
    ),
  checkModeratorStatus: (boardId: number) =>
    apiClient.get<
      ApiResponse<{ is_moderator: boolean; moderator?: Moderator }>
    >(`/boards/${boardId}/moderator-status`),
  create: (data: {
    name: string;
    slug: string;
    description?: string;
    icon?: string;
    cover?: string;
    parent_id?: number;
    sort_order?: number;
    view_role?: string;
    post_role?: string;
    reply_role?: string;
  }) => apiClient.post<ApiResponse<Board>>("/boards", data),

  update: (
    id: number,
    data: Partial<{
      name: string;
      slug: string;
      description: string;
      icon: string;
      cover: string;
      parent_id: number;
      sort_order: number;
      view_role: string;
      post_role: string;
      reply_role: string;
    }>,
  ) => apiClient.put<ApiResponse<Board>>(`/boards/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/boards/${id}`),

  // 版主管理
  getModerators: (boardId: number) =>
    apiClient.get<ApiResponse<Moderator[]>>(`/boards/${boardId}/moderators`),

  addModerator: (
    boardId: number,
    data: {
      user_id: number;
      can_delete_post?: boolean;
      can_pin_post?: boolean;
      can_edit_any_post?: boolean;
      can_manage_moderator?: boolean;
      can_ban_user?: boolean;
    },
  ) => apiClient.post<ApiResponse<null>>(`/boards/${boardId}/moderators`, data),

  removeModerator: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<null>>(
      `/boards/${boardId}/moderators/${userId}`,
    ),

  // 禁言管理
  banUser: (
    boardId: number,
    data: {
      user_id: number;
      reason: string;
      expires_at?: string;
    },
  ) => apiClient.post<ApiResponse<null>>(`/boards/${boardId}/bans`, data),

  unbanUser: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<null>>(`/boards/${boardId}/bans/${userId}`),

  // 帖子管理（版主）
  deletePost: (boardId: number, postId: number) =>
    apiClient.delete<ApiResponse<null>>(`/boards/${boardId}/posts/${postId}`),

  pinPost: (boardId: number, postId: number, data: { pin_in_board: boolean }) =>
    apiClient.put<ApiResponse<null>>(
      `/boards/${boardId}/posts/${postId}/pin`,
      data,
    ),
};

// ─── Timeline (时间线) ───────────────────────────────────────────────────────

export const timelineApi = {
  getHomeTimeline: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TimelineEvent>>>("/timeline", {
      params,
    }),

  getFollowingTimeline: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TimelineEvent>>>("/timeline/following", {
      params,
    }),

  subscribe: (userId: number) =>
    apiClient.post<ApiResponse<null>>(`/timeline/subscribe/${userId}`),

  unsubscribe: (userId: number) =>
    apiClient.delete<ApiResponse<null>>(`/timeline/subscribe/${userId}`),

  getSubscriptions: () =>
    apiClient.get<ApiResponse<Subscription[]>>("/timeline/subscriptions"),
};

// ─── Topics (专题) ───────────────────────────────────────────────────────────

export const topicApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Topic>>>("/topics", { params }),

  getById: (id: number) => apiClient.get<ApiResponse<Topic>>(`/topics/${id}`),

  getPosts: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TopicPost>>>(`/topics/${id}/posts`, {
      params,
    }),

  getFollowers: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TopicFollow>>>(
      `/topics/${id}/followers`,
      { params },
    ),

  isFollowing: (id: number) =>
    apiClient.get<ApiResponse<{ is_following: boolean }>>(
      `/topics/${id}/is-following`,
    ),

  create: (data: {
    title: string;
    description?: string;
    cover?: string;
    is_public?: boolean;
  }) => apiClient.post<ApiResponse<Topic>>("/topics", data),

  update: (
    id: number,
    data: {
      title?: string;
      description?: string;
      cover?: string;
      is_public?: boolean;
    },
  ) => apiClient.put<ApiResponse<Topic>>(`/topics/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/topics/${id}`),

  addPost: (id: number, data: { post_id: number; sort_order?: number }) =>
    apiClient.post<ApiResponse<null>>(`/topics/${id}/posts`, data),

  removePost: (id: number, postId: number) =>
    apiClient.delete<ApiResponse<null>>(`/topics/${id}/posts/${postId}`),

  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/topics/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/topics/${id}/follow`),
};

// ─── Admin ───────────────────────────────────────────────────────────────────

export const adminApi = {
  listUsers: (params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }) => apiClient.get<ApiResponse<PageData<User>>>("/admin/users", { params }),

  setUserActive: (id: number, active: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, { active }),

  listPosts: (params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }) => apiClient.get<ApiResponse<PageData<Post>>>("/admin/posts", { params }),

  togglePin: (id: number) =>
    apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/pin`),

  togglePinInBoard: (id: number, data: { pin_in_board: boolean }) =>
    apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/pin-board`, data),

  // 板块管理
  listBoards: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Board>>>("/admin/boards", { params }),

  updateBoardSort: (id: number, data: { sort_order: number }) =>
    apiClient.put<ApiResponse<null>>(`/admin/boards/${id}/sort`, data),
};
