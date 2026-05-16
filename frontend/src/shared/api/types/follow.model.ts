import { BaseModel } from "./basic.model";
import { UserDO } from "./user.model";

export interface Follow extends BaseModel {
  follower_id: number;
  following_id: number;
  follower?: UserDO;
  following?: UserDO;
}
