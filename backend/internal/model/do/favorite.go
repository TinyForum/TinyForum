package do

import "tiny-forum/internal/model/common"

// FavoriteStatus 收藏状态
type FavoriteStatus int8

const (
	FavoriteStatusActive   FavoriteStatus = 1
	FavoriteStatusCanceled FavoriteStatus = 2
)

type FavoriteGroup struct {
	common.BaseModel
	UserID    int64  `gorm:"not null;index:idx_fg_user_id"`
	Name      string `gorm:"size:100;not null"`
	IsDefault bool   `gorm:"default:false"` // 标识默认收藏夹
	IsPrivate bool   `gorm:"default:false"`
	SortOrder int    `gorm:"default:0"`
}

func (FavoriteGroup) TableName() string { return "favorite_groups" }

type Favorite struct {
	common.BaseModel
	UserID     int64          `gorm:"not null;uniqueIndex:idx_fav_unique,priority:1"`
	TargetID   int64          `gorm:"not null;uniqueIndex:idx_fav_unique,priority:2"`
	TargetType string         `gorm:"size:50;not null;uniqueIndex:idx_fav_unique,priority:3"`
	GroupID    int64          `gorm:"not null;uniqueIndex:idx_fav_unique,priority:4"` // 非 NULL
	Status     FavoriteStatus `gorm:"default:1;index:idx_fav_status"`
}

func (Favorite) TableName() string { return "favorites" }
