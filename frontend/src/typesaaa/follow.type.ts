import { User } from "./user.type";
export interface Follow extends BaseModel {
  follower_id: number;
  following_id: number;

  follower?: User;
  following?: User;
}
