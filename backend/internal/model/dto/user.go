package dto

// type UserPosts struct {
// 	title

// }
// SimpleAuthor 精简的作者信息
type SimpleAuthor struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}
