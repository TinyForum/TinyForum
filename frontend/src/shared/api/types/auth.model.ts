import { UserDO } from "./user.model.do";

export interface RegisterPayload {
  username: string;
  email: string;
  password: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface AuthResultVO {
  token: string;
  user: UserDO;
}

export interface LoginResultVO {
  token: string;
  user: UserDO;
}
