

// // 新增板块相关类型
// export interface Board {
//   id: number;
//   name: string;
//   slug: string;
//   description: string;
//   icon: string;
//   cover: string;
//   parent_id: number | null;
//   sort_order: number;
//   view_role: 'user' | 'admin';
//   post_role: 'user' | 'admin';
//   reply_role: 'user' | 'admin';
//   post_count: number;
//   thread_count: number;
//   today_count: number;
//   created_at: string;
//   updated_at: string;
//   parent?: Board;
//   children?: Board[];
// }

// export interface Moderator {
//   id: number;
//   user_id: number;
//   board_id: number;
//   can_delete_post: boolean;
//   can_pin_post: boolean;
//   can_edit_any_post: boolean;
//   can_manage_moderator: boolean;
//   can_ban_user: boolean;
//   created_at: string;
//   user?: User;
//   board?: Board;
// }

// export interface BoardBan {
//   id: number;
//   user_id: number;
//   board_id: number;
//   banned_by: number;
//   reason: string;
//   expires_at: string | null;
//   created_at: string;
//   user?: User;
//   board?: Board;
//   banner?: User;
// }

// // 问答相关类型
// export interface Question {
//   id: number;
//   post_id: number;
//   accepted_answer_id: number | null;
//   reward_score: number;
//   answer_count: number;
//   created_at: string;
//   updated_at: string;
//   post?: Post;
//   accepted_answer?: Comment;
// }

// export interface AnswerVoteResult {
//   vote_type: 'up' | 'down' | '';
//   vote_count: number;
//   action: 'added' | 'removed' | 'updated';
// }

// // 时间线相关类型
// export interface TimelineEvent {
//   id: number;
//   user_id: number;
//   actor_id: number;
//   action: 'create_post' | 'create_comment' | 'like_post' | 'like_comment' | 'follow_user' | 'accept_answer' | 'sign_in';
//   target_id: number;
//   target_type: string;
//   payload: any;
//   score: number;
//   created_at: string;
//   user?: User;
//   actor?: User;
// }

// export interface Subscription {
//   id: number;
//   subscriber_id: number;
//   target_user_id: number;
//   target_type: string;
//   target_id: number;
//   is_active: boolean;
//   created_at: string;
// }

// // 专题相关类型
// export interface Topic {
//   id: number;
//   title: string;
//   description: string;
//   cover: string;
//   creator_id: number;
//   is_public: boolean;
//   post_count: number;
//   follower_count: number;
//   created_at: string;
//   updated_at: string;
//   creator?: User;
// }

// export interface TopicPost {
//   id: number;
//   topic_id: number;
//   post_id: number;
//   sort_order: number;
//   added_by: number;
//   created_at: string;
//   post?: Post;
//   topic?: Topic;
// }

// export interface TopicFollow {
//   id: number;
//   user_id: number;
//   topic_id: number;
//   created_at: string;
//   user?: User;
//   topic?: Topic;
// }

// // 扩展 Post 类型
// export interface ExtendedPost extends Post {
//   board_id?: number;
//   board?: Board;
//   is_question?: boolean;
//   question?: Question;
//   pin_in_board?: boolean;
// }

// // 扩展 Comment 类型
// export interface ExtendedComment extends Comment {
//   is_answer?: boolean;
//   is_accepted?: boolean;
//   vote_count?: number;
// }

// // 版主申请
// export interface ModeratorApplication {
//   id: number;
//   user_id: number;
//   board_id: number;
//   reason: string;
//   status: 'pending' | 'approved' | 'rejected';
//   handled_by: number | null;
//   handle_note: string;
//   created_at: string;
//   updated_at: string;
//   user?: User;
//   board?: Board;
//   handler?: User;
// }