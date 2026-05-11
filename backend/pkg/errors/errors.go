package apperrors

// ========== 预定义错误实例 ==========
// 注意：这些是"模板"，业务层若需要附加 Detail 或 Cause，
// 请使用链式方法（如 ErrUserNotFound.WithDetail(...).WithCause(err)），
// 不要直接修改这些全局变量。

var (
	// =========================================================================================== 通用

	ErrUnknown             = New(CodeUnknown, "未知错误")                 // 未知错误
	ErrValidation          = New(CodeValidation, "参数验证失败")            // 参数验证失败
	ErrUnauthorized        = New(CodeUnauthorized, "未授权，请先登录")        // 未授权，请先登录
	ErrForbidden           = New(CodeForbidden, "权限不足")               // 权限不足
	ErrNotFound            = New(CodeNotFound, "资源不存在")               // 资源不存在
	ErrTooManyRequests     = New(CodeTooManyRequests, "请求过于频繁，请稍后再试") // 请求过于频繁，请稍后再试
	ErrInternalError       = New(CodeInternalError, "服务器内部错误")        // 服务器内部错误
	ErrInvalidRequest      = New(CodeInvalidRequest, "无效的请求")         // 无效的请求
	ErrSystemBusy          = New(CodeSystemBusy, "系统繁忙，请稍后再试")        // 系统繁忙，请稍后再试
	ErrInvalidConfirmation = New(CodeInvalidConfirmation, "无效的确认信息")  // 无效的确认信息

	// =========================================================================================== 用户模块

	ErrUserNotFound           = New(CodeUserNotFound, "用户不存在")                      // 用户不存在
	ErrUserExist              = New(CodeUserExist, "用户已存在")                         // 用户已存在
	ErrInvalidEmail           = New(CodeInvalidEmail, "无效的邮箱地址")                    // 无效的邮箱地址
	ErrInvalidPhone           = New(CodeInvalidPhone, "无效的手机号码")                    // 无效的手机号码
	ErrInvalidPassword        = New(CodeInvalidPassword, "无效的密码")                   // 无效的密码
	ErrInvalidCurrentPassword = New(CodeInvalidCurrentPassword, "当前密码不正确")          // 当前密码不正确
	ErrInvalidUsername        = New(CodeInvalidUsername, "无效的用户名")                  // 无效的用户名
	ErrInvalidAvatar          = New(CodeInvalidAvatar, "无效的头像链接")                   // 无效的头像链接
	ErrInvalidNickname        = New(CodeInvalidNickname, "无效的昵称")                   // 无效的昵称
	ErrInvalidUserID          = New(CodeInvalidUserID, "无效的用户ID")                   // 无效的用户ID
	ErrInvalidRole            = New(CodeInvalidRole, "无法更改到此角色类型")                  // 无法更改到此角色类型
	ErrCannotModifySelf       = New(CodeCannotModifySelf, "不能修改自己的信息")              // 不能修改自己的信息
	ErrCannotChangeOwnerRole  = New(CodeCannotChangeOwner, "不能修改超级管理员的角色")          // 不能修改超级管理员的角色
	ErrFollowSelf             = New(CodeFollowSelf, "不能关注自己")                       // 不能关注自己
	ErrAlreadyFollow          = New(CodeAlreadyFollow, "已经关注了该用户")                  // 已经关注了该用户
	ErrNotFollow              = New(CodeNotFollow, "尚未关注该用户")                       // 尚未关注该用户
	ErrScoreNotEnough         = New(CodeScoreNotEnough, "积分不足")                     // 积分不足
	ErrUserBlocked            = New(CodeUserBlocked, "用户已被封禁")                      // 用户已被封禁
	ErrUserDeleted            = New(CodeUserDeleted, "用户已被删除")                      // 用户已被删除
	ErrCannotBlockSelf        = New(CodeCannotModifySelf, "不能封禁自己的账号")              // 不能封禁自己的账号
	ErrCannotModifySuperAdmin = New(CodeCannotChangeOwner, "不能修改超级管理员")             // 不能修改超级管理员
	ErrCannotBlockAdmin       = New(CodeInsufficientPermission, "只有超级管理员才能封禁其他管理员") //	只有超级管理员才能封禁其他管理员

	// 密码校验
	ErrPasswordNotMatch  = New(CodeUnauthorized, "密码不匹配")                  // 密码不匹配
	ErrPasswordTooShort  = New(CodePasswordTooShort, "密码长度至少为 8 位")        // 密码长度至少为6位
	ErrPasswordTooLong   = New(CodeInvalidPassword, "密码长度不能超过32位")         // 密码长度不能超过32位
	ErrPasswordSameAsOld = New(CodePasswordSameAsOld, "新密码不能与旧密码相同")       // 新密码不能与旧密码相同
	ErrWeakPassword      = New(CodeInvalidPassword, "密码强度太弱，请使用更长且更复杂的密码") // 密码强度太弱，请使用更长且更复杂的密码

	// Token / 验证码
	ErrTokenExpired    = New(CodeUnauthorized, "Token已过期") // Token已过期
	ErrInvalidToken    = New(CodeUnauthorized, "无效的Token") // 无效的Token
	ErrRequiredToken   = New(CodeValidation, "需要Token")    // 需要Token
	ErrRequiredCaptcha = New(CodeValidation, "需要验证码")      // 需要验证码
	ErrInvalidCaptcha  = New(CodeValidation, "验证码错误")      // 验证码错误

	// 内容模块 - 帖子
	ErrPostNotFound = New(CodePostNotFound, "帖子不存在")     // 帖子不存在
	ErrPostLocked   = New(CodePostLocked, "帖子已被锁定，无法操作") // 帖子已被锁定，无法操作
	ErrPostPinned   = New(CodePostPinned, "帖子置顶状态冲突")    // 帖子置顶状态冲突
	ErrPostDeleted  = New(CodePostDeleted, "帖子已被删除")     // 帖子已被删除

	// 内容模块 - 板块
	ErrBoardNotFound = New(CodeBoardNotFound, "板块不存在") // 板块不存在

	// 内容模块 - 问答
	ErrAcceptForbidden  = New(CodeAcceptForbidden, "只有发帖人才能采纳答案") // 只有发帖人才能采纳答案
	ErrAnswerNotFound   = New(CodeAnswerNotFound, "回答不存在")        // 回答不存在
	ErrQuestionNotFound = New(CodeQuestionNotFound, "问题不存在")      // 问题不存在

	// 内容模块 - 评论
	ErrCommentNotFound = New(CodeCommentNotFound, "评论不存在") // 评论不存在
	ErrCommentDeleted  = New(CodeCommentDeleted, "评论已被删除") // 评论已被删除

	// 内容模块 - 主题 / 标签 / 通知
	ErrTopicNotFound        = New(CodeTopicNotFound, "主题不存在")        // 主题不存在
	ErrTagNotFound          = New(CodeTagNotFound, "标签不存在")          // 标签不存在
	ErrNotificationNotFound = New(CodeNotificationNotFound, "通知不存在") // 通知不存在

	// 点赞 / 收藏
	ErrLikeAlready     = New(CodeLikeAlready, "已经点过赞了")    // 已经点过赞了
	ErrLikeNotExist    = New(CodeLikeNotExist, "尚未点赞")     // 尚未点赞
	ErrCollectAlready  = New(CodeCollectAlready, "已经收藏过了") // 已经收藏过了
	ErrCollectNotExist = New(CodeCollectNotExist, "尚未收藏")  // 尚未收藏

	// 权限模块 - 版主申请
	ErrModeratorApplyExist    = New(CodeModeratorApplyExist, "已经提交过版主申请，请勿重复提交") // 已经提交过版主申请，请勿重复提交
	ErrModeratorApplyNotFound = New(CodeModeratorApplyNotFound, "版主申请不存在")       // 版主申请不存在
	ErrAlreadyModerator       = New(CodeAlreadyModerator, "已经是版主，无需重复申请")        // 已经是版主，无需重复申请
	ErrInsufficientPermission = New(CodeInsufficientPermission, "权限不足")          // 权限不足

	// 积分模块
	ErrFailedToQueryScore  = New(CodeFailedToQueryScore, "查询积分失败")   // 查询积分失败
	ErrScoreRecordNotFound = New(CodeScoreRecordNotFound, "积分记录不存在") // 积分记录不存在

	// 公告模块
	ErrAnnouncementNotFound = New(CodeAnnouncementNotFound, "公告不存在")           // 公告不存在
	ErrInvalidPublishTime   = New(CodeAnnouncementInvalidTime, "发布时间无效")       // 发布时间无效
	ErrExpiredTimeInvalid   = New(CodeAnnouncementInvalidTime, "过期时间必须晚于发布时间") // 过期时间必须晚于发布时间

	// 统计模块
	ErrStatsNotFound = New(CodeStatsNotFound, "统计数据不存在") // 统计数据不存在

	// 时间线模块
	ErrTimelineEmpty = New(CodeTimelineEmpty, "时间线暂无内容") // 时间线暂无内容

	// 文件插件模块
	ErrFileTooLarge        = New(CodeFileTooLarge, "文件过大")                      // 文件过大
	ErrFileTypeInvalid     = New(CodeFileTypeInvalid, "不支持的文件类型")               // 不支持的文件类型
	ErrUploadFailed        = New(CodeUploadFailed, "文件上传失败")                    // 文件上传失败
	ErrInvalidPluginFormat = New(CodeInvalidPluginFormat, "插件格式无效")             // 插件格式无效
	ErrManifestNotFound    = New(CodeManifestNotFound, "没有找到 manifest.json 文件") // 没有找到 manifest.json 文件
	ErrInvalidManifest     = New(CodeInvalidManifest, "文件 manifest.json 无效")    // 文件 manifest.json 无效
	ErrPluginAlreadyExist  = New(CodePluginAlreadyExist, "插件已存在")               // 插件已存在
	ErrPluginEnabledFirst  = New(CodePluginEnabledFirst, "请先禁用插件")              // 请先禁用插件
)
