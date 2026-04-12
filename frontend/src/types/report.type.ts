import { User } from "./user.type";

type ReportStatus = 'pending' | 'resolved' | 'rejected';

export interface Report extends BaseModel {
  reporter_id: number;
  target_id: number;
  target_type: string; // post | comment | user
  reason: string;
  status: ReportStatus;
  handler_id: number | null;
  handle_note: string;

  reporter?: User;
  handler?: User;
}