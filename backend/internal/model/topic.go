package model

type Topic struct {
	BaseModel
	Title         string `gorm:"not null;size:100" json:"title"`
	Description   string `gorm:"size:500" json:"description"`
	Cover         string `gorm:"size:500" json:"cover"`
	CreatorID     uint   `gorm:"not null;index" json:"creator_id"`
	IsPublic      bool   `gorm:"default:true" json:"is_public"`
	PostCount     int    `gorm:"default:0" json:"post_count"`
	FollowerCount int    `gorm:"default:0" json:"follower_count"`

	Creator   User          `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Posts     []TopicPost   `gorm:"foreignKey:TopicID" json:"-"`
	Followers []TopicFollow `gorm:"foreignKey:TopicID" json:"-"`
}

type TopicPost struct {
	BaseModel
	TopicID   uint `gorm:"not null;uniqueIndex:idx_topic_post" json:"topic_id"`
	PostID    uint `gorm:"not null;uniqueIndex:idx_topic_post" json:"post_id"`
	SortOrder int  `gorm:"default:0" json:"sort_order"`
	AddedBy   uint `json:"added_by"`

	Topic Topic `gorm:"foreignKey:TopicID" json:"-"`
	Post  Post  `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

type TopicFollow struct {
	BaseModel
	UserID  uint `gorm:"not null;uniqueIndex:idx_user_topic" json:"user_id"`
	TopicID uint `gorm:"not null;uniqueIndex:idx_user_topic" json:"topic_id"`

	User  User  `gorm:"foreignKey:UserID" json:"-"`
	Topic Topic `gorm:"foreignKey:TopicID" json:"-"`
}
