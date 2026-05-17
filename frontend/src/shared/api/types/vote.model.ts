export interface VoteAnswerPayload {
  vote_type: "up" | "down";
}

export interface VoteStatusResponse {
  down_count: number;
  up_count: number;
  user_vote: number; // 1: 赞同, -1: 反对, 0: 未投票
}

export type VoteType = "up" | "down";

export interface AnswerVoteResult {
  vote_count: number;
  user_vote?: VoteType;
}
