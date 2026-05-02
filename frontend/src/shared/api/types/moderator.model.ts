import { User } from "@/shared/type/admin.types";
import { Board } from "./board.model";

export interface Moderator {
  id: number;
  user_id: number;
  board_id: number;
  user?: User;
  board?: Board;
  can_delete_post: boolean;
  can_pin_post: boolean;
  can_edit_any_post: boolean;
  can_manage_moderator: boolean;
  can_ban_user: boolean;
}
