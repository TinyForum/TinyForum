package user

// avatarURL 生成默认头像URL
func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}
