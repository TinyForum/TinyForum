package do

import "tiny-forum/internal/model/common"

// PluginStatus 插件运行状态
type PluginStatus string

const (
	PluginStatusActive   PluginStatus = "active"
	PluginStatusInactive PluginStatus = "inactive"
	PluginStatusError    PluginStatus = "error"
	PluginStatusLoading  PluginStatus = "loading"
)

// PluginType 插件类型
type PluginType string

const (
	PluginTypeUI      PluginType = "ui"      // 前端组件注入
	PluginTypeBackend PluginType = "backend" // 纯后端逻辑扩展
	PluginTypeLib     PluginType = "lib"     // 通用工具库
	PluginTypeApp     PluginType = "app"     // 独立子应用
	PluginTypeMiniapp PluginType = "miniapp" // 小程序/H5轻应用
)

// PluginPricingModel 价格模型
type PluginPricingModel string

const (
	PricingFree     PluginPricingModel = "free"
	PricingFreemium PluginPricingModel = "freemium"
	PricingPaid     PluginPricingModel = "paid"
)

// PluginPricing 价格信息
type PluginPricing struct {
	Model       PluginPricingModel `json:"model" gorm:"type:varchar(20)"`
	Price       *float64           `json:"price,omitempty" gorm:"type:decimal(10,2)"`      // 付费价格（元）
	Cycle       string             `json:"cycle,omitempty" gorm:"type:varchar(20)"`        // once/monthly/yearly
	FreeLimit   string             `json:"freeLimit,omitempty" gorm:"type:text"`           // 免费版限制说明
	PurchaseURL string             `json:"purchaseUrl,omitempty" gorm:"type:varchar(255)"` // 购买链接
}

// PluginCategory 插件分类
// PluginCategory 插件分类（行业标准版）
type PluginCategory string

const (
	CategoryCommerce    PluginCategory = "commerce"    // 电商/购物/分销/返利/O2O
	CategoryMarketing   PluginCategory = "marketing"   // 营销推广/活动/投票/SEO
	CategoryContent     PluginCategory = "content"     // 内容管理/知识付费/小说/漫画/问答
	CategoryCommunity   PluginCategory = "community"   // 社区/论坛/社交/评论/私信
	CategoryMedia       PluginCategory = "media"       // 音视频/直播/音乐/图库
	CategoryEducation   PluginCategory = "education"   // 教育培训/在线课程/考试
	CategorySupport     PluginCategory = "support"     // 客服/工单/聊天/问卷/IM
	CategoryPayment     PluginCategory = "payment"     // 支付网关/订阅/发票
	CategorySecurity    PluginCategory = "security"    // 认证/登录/权限/验证/防刷
	CategoryStorage     PluginCategory = "storage"     // 存储/OSS/COS/CDN/备份
	CategoryAnalytics   PluginCategory = "analytics"   // 数据统计/分析/报表
	CategoryDesign      PluginCategory = "design"      // UI组件/主题/样式/布局
	CategoryDevtools    PluginCategory = "devtools"    // API扩展/后端逻辑/开发工具
	CategoryIntegration PluginCategory = "integration" // 第三方集成/Webhook/通知
	CategoryUtility     PluginCategory = "utility"     // 系统工具/增强功能/任务调度
)

// PluginCategoryLabels 分类中文标签
var PluginCategoryLabels = map[PluginCategory]string{
	CategoryCommerce:    "电商/O2O",
	CategoryMarketing:   "营销推广",
	CategoryContent:     "内容/知识付费",
	CategoryCommunity:   "社区/互动",
	CategoryMedia:       "媒体/音视频",
	CategoryEducation:   "教育培训",
	CategorySupport:     "客服/工单/IM",
	CategoryPayment:     "支付/财务",
	CategorySecurity:    "安全/认证",
	CategoryStorage:     "存储/备份",
	CategoryAnalytics:   "统计/分析",
	CategoryDesign:      "UI/主题设计",
	CategoryDevtools:    "开发/API",
	CategoryIntegration: "外部集成",
	CategoryUtility:     "系统工具",
}

// CompatVersion 平台兼容版本
type CompatVersion string

const (
	CompatV0 CompatVersion = "v0"
	CompatV1 CompatVersion = "v1"
	CompatV2 CompatVersion = "v2"
	CompatV3 CompatVersion = "v3"
)

// PluginCompatibility 兼容性信息
type PluginCompatibility struct {
	Versions  []CompatVersion `json:"versions" gorm:"type:json"`                 // 兼容的平台版本
	MinNode   string          `json:"minNode,omitempty" gorm:"type:varchar(50)"` // 最低Node.js版本
	Requires  []string        `json:"requires,omitempty" gorm:"type:json"`       // 依赖的其他插件ID
	Conflicts []string        `json:"conflicts,omitempty" gorm:"type:json"`      // 冲突的插件ID
}

// PluginConfigField 配置字段定义
type PluginConfigField struct {
	Key          string      `json:"key" gorm:"type:varchar(100)"`
	Label        string      `json:"label" gorm:"type:varchar(100)"`
	Type         string      `json:"type" gorm:"type:varchar(20)"` // text/number/boolean/select/textarea/color/url/secret
	DefaultValue interface{} `json:"defaultValue,omitempty" gorm:"type:json"`
	Placeholder  string      `json:"placeholder,omitempty" gorm:"type:varchar(200)"`
	Description  string      `json:"description,omitempty" gorm:"type:text"`
	Required     bool        `json:"required"`
	Options      []struct {
		Label string      `json:"label"`
		Value interface{} `json:"value"`
	} `json:"options,omitempty" gorm:"type:json"`
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// PluginPermission 权限声明
type PluginPermission string

const (
	PermReadUser      PluginPermission = "read:user"
	PermReadPosts     PluginPermission = "read:posts"
	PermWritePosts    PluginPermission = "write:posts"
	PermReadComments  PluginPermission = "read:comments"
	PermWriteComments PluginPermission = "write:comments"
	PermReadSettings  PluginPermission = "read:settings"
	PermWriteSettings PluginPermission = "write:settings"
	PermSendEmail     PluginPermission = "send:email"
	PermSendSMS       PluginPermission = "send:sms"
	PermStorageRead   PluginPermission = "storage:read"
	PermStorageWrite  PluginPermission = "storage:write"
	PermPayment       PluginPermission = "payment"
	PermNetwork       PluginPermission = "network"
)

// PluginMeta 插件元数据（数据库模型）
type PluginMeta struct {
	common.BaseModel
	// 基础标识
	Name        string   `json:"name" gorm:"type:varchar(100);not null;index:idx_name,unique"` // 插件名称
	Version     string   `json:"version" gorm:"type:varchar(50);not null"`                     // 插件版本
	Description string   `json:"description" gorm:"type:text"`                                 // 插件描述
	Summary     string   `json:"summary,omitempty" gorm:"type:varchar(300)"`                   // 一句话简介
	IconURL     string   `json:"iconUrl,omitempty" gorm:"type:varchar(255)"`                   // 插件图标
	Screenshots []string `json:"screenshots,omitempty" gorm:"type:json"`                       // 截图列表
	HomepageURL string   `json:"homepageUrl,omitempty" gorm:"type:varchar(255)"`               // 官网地址

	// 分类与类型
	Type     PluginType     `json:"type" gorm:"type:varchar(20);not null;index"`     // 插件类型（对于服务端）
	Category PluginCategory `json:"category" gorm:"type:varchar(30);not null;index"` // 插件分类（对于业务）
	Tags     []string       `json:"tags,omitempty" gorm:"type:json"`                 // 标签列表

	// 作者信息
	AuthorID    uint   `json:"authorId,omitempty" gorm:"type:bigint(20)"`
	AuthorEmail string `json:"authorEmail,omitempty" gorm:"type:varchar(100)"` // 作者邮箱
	AuthorURL   string `json:"authorUrl,omitempty" gorm:"type:varchar(255)"`   // 作者主页

	// 加载配置
	ScriptURL   string   `json:"scriptUrl" gorm:"type:varchar(500);not null"`    // 前端脚本入口
	ServerEntry string   `json:"serverEntry,omitempty" gorm:"type:varchar(255)"` // 后端服务入口
	Slots       []string `json:"slots,omitempty" gorm:"type:json"`               // 注入的插槽名称列表
	Routes      []string `json:"routes,omitempty" gorm:"type:json"`              // 注册的路由路径

	// 价格与兼容性
	Pricing       PluginPricing       `json:"pricing" gorm:"type:json;serializer:json"`       // 价格信息
	Compatibility PluginCompatibility `json:"compatibility" gorm:"type:json;serializer:json"` // 兼容性信息

	// 权限
	Permissions []PluginPermission `json:"permissions,omitempty" gorm:"type:json"` // 权限声明

	// 运行时（服务端写入，前端只读）
	Enabled      bool         `json:"enabled" gorm:"default:false;index"`                // 是否启用
	Status       PluginStatus `json:"status" gorm:"type:varchar(20);default:'inactive'"` // 插件状态
	InstallCount int          `json:"installCount" gorm:"default:0"`                     // 安装次数
	Rating       float32      `json:"rating" gorm:"type:decimal(2,1);default:0"`         // 评分 0~5

	// 配置
	ConfigSchema []PluginConfigField `json:"configSchema,omitempty" gorm:"type:json;serializer:json"` // 配置字段定义
	Config       map[string]any      `json:"config,omitempty" gorm:"type:json;serializer:json"`       // 配置值
}

// TableName 指定表名
func (PluginMeta) TableName() string {
	return "plugins"
}
