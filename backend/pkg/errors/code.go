package apperrors

// ========== 错误码常量 ==========
// 规则：模块前缀 + 4位序号
// 通用(1xxxx) 用户(2xxxx) 内容(3xxxx) 权限(4xxxx) 积分(5xxxx) 公告(6xxxx) 统计(7xxxx) 时间线(8xxxx) 文件(9xxxx)
const (
	// 通用 (10000-10099)
	CodeUnknown             = 10000 // 未知错误
	CodeValidation          = 10001 // 参数校验错误
	CodeUnauthorized        = 10002 // 未授权
	CodeForbidden           = 10003 // 禁止访问
	CodeNotFound            = 10004 // 未找到
	CodeTooManyRequests     = 10005 // 请求过多
	CodeInternalError       = 10006 // 内部错误
	CodeInvalidRequest      = 10007 // 无效请求
	CodeSystemBusy          = 10008 // 系统繁忙
	CodeInvalidConfirmation = 10009 // 无效确认

	// 用户模块 (20000-20999)
	CodeUserNotFound           = 20001 // 用户不存在
	CodeUserExist              = 20002 // 用户已存在
	CodeInvalidEmail           = 20010 // 无效邮箱
	CodeInvalidPhone           = 20011 // 无效手机号
	CodeInvalidPassword        = 20012 // 无效密码
	CodeInvalidUsername        = 20013 // 无效用户名
	CodeInvalidAvatar          = 20014 // 无效头像
	CodeInvalidNickname        = 20015 // 无效昵称
	CodeInvalidUserID          = 20016 // 无效用户ID
	CodeInvalidCurrentPassword = 20017 // 无效当前密码
	CodeInvalidRole            = 20020 // 无效角色
	CodeCannotModifySelf       = 20021 // 不能修改自己
	CodeCannotChangeOwner      = 20022 // 不能修改拥有者
	CodeFollowSelf             = 20040 // 不能关注自己
	CodeAlreadyFollow          = 20041 // 已关注
	CodeNotFollow              = 20042 // 未关注
	CodeScoreNotEnough         = 20060 // 积分不足
	CodeUserBlocked            = 20080 // 用户被禁用
	CodeUserDeleted            = 20081 // 用户被删除

	// 内容模块 (30000-30999)
	CodePostNotFound         = 30001 // 帖子不存在
	CodePostLocked           = 30002 // 帖子已锁定
	CodePostPinned           = 30003 // 帖子已置顶
	CodePostDeleted          = 30004 // 帖子已删除
	CodeBoardNotFound        = 30020 // 版块不存在
	CodeQuestionNotFound     = 30040 // 问题不存在
	CodeAnswerNotFound       = 30041 // 答案不存在
	CodeCommentNotFound      = 30060 // 评论不存在
	CodeCommentDeleted       = 30061 // 评论已删除
	CodeTopicNotFound        = 30080 // 话题不存在
	CodeTagNotFound          = 30081 // 标签不存在
	CodeNotificationNotFound = 30082 // 通知不存在
	CodeLikeAlready          = 30100 // 点赞已存在
	CodeLikeNotExist         = 30101 // 点赞不存在
	CodeCollectAlready       = 30102 // 收藏已存在
	CodeCollectNotExist      = 30103 // 收藏不存在
	CodePasswordTooShort     = 30120 // 密码太短
	CodePasswordSameAsOld    = 30121 // 新密码与旧密码相同

	// 权限模块 (40000-40999)
	//
	CodeInsufficientPermission = 40001 // 权限不足
	CodeAcceptForbidden        = 40002 // 不能接受
	CodeModeratorApplyExist    = 40003 // 申请已存在
	CodeModeratorApplyNotFound = 40004 // 申请不存在
	CodeAlreadyModerator       = 40005 // 已经是版主

	// 积分模块 (50000-50999)
	CodeFailedToQueryScore  = 50001 // 查询积分失败
	CodeScoreRecordNotFound = 50002 // 积分记录不存在

	// 公告模块 (60000-60999)
	CodeAnnouncementNotFound    = 60001 // 公告不存在
	CodeAnnouncementInvalidTime = 60002 // 公告时间无效

	// 统计模块 (70000-70999)
	CodeStatsNotFound = 70001 // 统计不存在

	// 时间线模块 (80000-80999)
	CodeTimelineEmpty = 80001 // 时间线为空

	// 文件插件模块 (90000-90999)
	CodeFileTooLarge        = 90001 // 文件过大
	CodeFileTypeInvalid     = 90002 // 文件类型无效
	CodeUploadFailed        = 90003 // 上传失败
	CodeInvalidPluginFormat = 90004 // 无效插件格式
	CodeManifestNotFound    = 90005 // Manifest 文件不存在
	CodeInvalidManifest     = 90006 // 无效的 Manifest 文件
	CodePluginAlreadyExist  = 90007 // 插件已存在
	CodePluginEnabledFirst  = 90008 // 启用插件前需要先禁用
)
