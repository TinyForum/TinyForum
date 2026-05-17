package do

import (
	"tiny-forum/internal/model/common"
)

// FavoriteGroup 收藏夹
type FavoriteGroup struct {
	common.BaseModel
	UserID    int64  `gorm:"not null;index:idx_user_id"`                      // 用户
	Name      string `gorm:"size:100;not null"`                               // 收藏夹名称
	IsDefault bool   `gorm:"default:false;index:idx_user_default,priority:2"` // 是否为默认收藏夹
	IsPrivate bool   `gorm:"default:false"`                                   // 是否为私密收藏夹
	SortOrder int    `gorm:"default:0"`                                       // 排序顺序
}

func (FavoriteGroup) TableName() string { return "favorite_groups" }

// Favorite 收藏记录（使用软删除实现取消）
type Favorite struct {
	common.BaseModel
	UserID     int64  `gorm:"not null;uniqueIndex:idx_fav_unique,priority:1"`         // 收藏者
	TargetID   int64  `gorm:"not null;uniqueIndex:idx_fav_unique,priority:2"`         // 收藏对象
	TargetType string `gorm:"size:50;not null;uniqueIndex:idx_fav_unique,priority:3"` // 收藏对象类型
	GroupID    int64  `gorm:"not null;index"`                                         // 收藏夹
}

func (Favorite) TableName() string { return "favorites" }
