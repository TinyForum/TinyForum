import { Board } from "./board.type";
import { User } from "./user.type";

export interface Moderator extends BaseModel {
  userId: number;
  boardId: number;
  permissions: Record<string, any> | null;
  canDeletePost: boolean;
  canPinPost: boolean;
  canEditAnyPost: boolean;
  canManageModerator: boolean;
  canBanUser: boolean;

  user?: User;
  board?: Board;
}
