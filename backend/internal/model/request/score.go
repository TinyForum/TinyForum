package request

type GetUserScoreRequest struct{
	UserID uint   `form:"user_id"`
}