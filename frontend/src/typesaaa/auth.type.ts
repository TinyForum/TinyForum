export interface SignIn extends BaseModel {
  user_id: number;
  sign_date: Date;
  score: number;
  continued: number;
}
