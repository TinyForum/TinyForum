package do

import (
	"tiny-forum/internal/model/common"
)

type Board struct {
	common.BaseModel
	Name        string   `gorm:"not null;size:50;uniqueIndex" json:"name"`          // 板块名称
	Slug        string   `gorm:"not null;size:50;uniqueIndex" json:"slug"`          // 板块别名
	Description string   `gorm:"size:500" json:"description"`                       // 板块描述
	Icon        string   `gorm:"size:100" json:"icon"`                              // 板块图标
	Cover       string   `gorm:"size:500" json:"cover"`                             // 板块封面
	ParentID    *uint    `gorm:"index;default:null" json:"parent_id"`               // 父板块ID
	SortOrder   int      `gorm:"default:0" json:"sort_order"`                       // 排序顺序
	ViewRole    UserRole `gorm:"type:varchar(20);default:'user'" json:"view_role"`  // 查看权限
	PostRole    UserRole `gorm:"type:varchar(20);default:'user'" json:"post_role"`  // 发帖权限
	ReplyRole   UserRole `gorm:"type:varchar(20);default:'user'" json:"reply_role"` // 回复权限
	PostCount   int      `gorm:"default:0" json:"post_count"`                       // 帖子数量
	ThreadCount int      `gorm:"default:0" json:"thread_count"`                     // 主题数量
	TodayCount  int      `gorm:"default:0" json:"today_count"`                      // 今日发帖数量

	Parent     *Board      `gorm:"foreignKey:ParentID" json:"parent,omitempty" swaggerignore:"true"`   // 父板块
	Children   []Board     `gorm:"foreignKey:ParentID" json:"children,omitempty" swaggerignore:"true"` // 子板块
	Moderators []Moderator `gorm:"foreignKey:BoardID" json:"-"`                                        // 板块管理员
}

const SystemBoardID = 1

var DefaultWorldBoard = &Board{
	Name:        "世界",
	Slug:        "world",
	Description: "世界板块。所有的帖子默认都会被分配到这里。所有用户可见",
	Icon:        "",
	Cover:       "",
	ParentID:    nil,
	SortOrder:   0,
	ViewRole:    "user",
	PostRole:    "user",
	ReplyRole:   "user",

	PostCount:   0,
	ThreadCount: 0,
	TodayCount:  0,

	Parent:     nil,
	Children:   []Board{},
	Moderators: []Moderator{},
}
