package board

import (
	"context"
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *boardService) ApplyModerator(ctx context.Context, req request.ApplyModeratorRequest) error {
	// 1. 检查是否已经是版主
	isMod, err := s.boardRepo.IsModerator(req.UserID, req.BoardID)
	if err != nil {
		return fmt.Errorf("检查版主状态失败: %w", err)
	}
	if isMod {
		return errors.New("你已经是该板块的版主")
	}

	// 2. 检查是否存在待审核的申请
	existing, err := s.boardRepo.FindPendingApplication(req.UserID, req.BoardID)
	if err != nil {
		return fmt.Errorf("查询待审核申请失败: %w", err)
	}
	if existing != nil {
		return errors.New("你已有一条待审核的申请，请等待管理员处理")
	}

	// 3. 构建申请记录（直接使用传入的权限列表，或从旧布尔字段转换）
	//    推荐调用方直接传入 RequestedPermissions，若仍需兼容旧字段可进行转换
	app := &do.ModeratorApplication{
		UserID:               req.UserID,
		BoardID:              req.BoardID,
		Reason:               req.Reason,
		Status:               do.ApplicationPending,
		RequestedPermissions: req.RequestedPermissions, // 核心改动
	}

	// 可选：对请求的权限进行合法性校验（防止非法权限字符串）
	for _, perm := range app.RequestedPermissions {
		if !perm.IsValid() {
			return fmt.Errorf("无效的权限: %s", perm)
		}
	}

	// 4. 保存申请
	if err := s.boardRepo.CreateApplication(app); err != nil {
		return fmt.Errorf("提交申请失败: %w", err)
	}
	return nil
}

func (s *boardService) CancelApplication(applicationID, userID uint) error {
	app, err := s.boardRepo.GetApplicationByID(applicationID)
	if err != nil || app == nil {
		return errors.New("申请不存在")
	}
	if app.UserID != userID {
		return errors.New("无权操作此申请")
	}
	if app.Status != do.ApplicationPending {
		return errors.New("只能撤销待审核的申请")
	}
	app.Status = do.ApplicationCanceled
	return s.boardRepo.UpdateApplication(app)
}

func (s *boardService) GetUserApplications(userID uint, page, pageSize int) ([]do.ModeratorApplication, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.boardRepo.GetApplicationsByUserID(userID, page, pageSize)
}

// ReviewApplication 审核版主申请
func (s *boardService) ReviewApplication(_ context.Context, req request.ReviewApplicationRequest, reviewerID uint) error {
	// 1. 获取申请记录
	app, err := s.boardRepo.GetApplicationByID(req.ApplicationID)
	if err != nil {
		return fmt.Errorf("获取申请失败: %w", err)
	}
	if app == nil {
		return errors.New("申请不存在")
	}
	if app.Status != do.ApplicationPending {
		return errors.New("该申请已被处理")
	}

	// 2. 状态变更（使用状态机方法，封装业务规则）
	if req.Approve {
		if err := app.Approve(reviewerID, req.ReviewNote); err != nil {
			return err
		}
	} else {
		if err := app.Reject(reviewerID, req.ReviewNote); err != nil {
			return err
		}
	}

	// 3. 持久化状态
	if err := s.boardRepo.UpdateApplication(app); err != nil {
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 4. 拒绝分支：仅发送通知
	if !req.Approve {
		s.sendRejectionNotification(app, reviewerID, req.ReviewNote)
		return nil
	}

	// 5. 批准分支：授予版主权限
	if err := s.grantModeratorPermissions(app, req, reviewerID); err != nil {
		return err
	}
	s.sendApprovalNotification(app, reviewerID)
	return nil
}

// grantModeratorPermissions 授予版主权限（处理权限合并与创建）
func (s *boardService) grantModeratorPermissions(app *do.ModeratorApplication, req request.ReviewApplicationRequest, reviewerID uint) error {
	// 避免重复添加（幂等性）
	isMod, _ := s.boardRepo.IsModerator(app.UserID, app.BoardID)
	if isMod {
		s.writeLog(reviewerID, app.BoardID, "approve_application", "user", app.UserID, "用户已是版主，跳过创建")
		return nil
	}

	// 确定最终授予的权限列表
	finalPerms := resolveFinalPermissions(app, req)

	// 创建版主记录
	mod := &do.Moderator{
		UserID:      app.UserID,
		BoardID:     app.BoardID,
		Permissions: finalPerms, // 直接存储权限切片（GORM JSON 字段）
	}
	if err := s.boardRepo.AddModerator(mod); err != nil {
		return fmt.Errorf("创建版主失败: %w", err)
	}

	s.writeLog(reviewerID, app.BoardID, "approve_application", "user", app.UserID, "审批通过并授予权限")
	return nil
}

// resolveFinalPermissions 合并申请时请求的权限和审批人覆盖的权限
func resolveFinalPermissions(app *do.ModeratorApplication, req request.ReviewApplicationRequest) []do.ModeratorPermission {
	// 构建审批人指定的权限集合
	overridePerms := make(map[do.ModeratorPermission]bool)
	if req.CanDeletePost != nil && *req.CanDeletePost {
		overridePerms[do.PerModDeletePost] = true
	}
	if req.CanPinPost != nil && *req.CanPinPost {
		overridePerms[do.PerMoePinPost] = true
	}
	if req.CanEditAnyPost != nil && *req.CanEditAnyPost {
		overridePerms[do.PerModEditAnyPost] = true
	}
	if req.CanManageModerator != nil && *req.CanManageModerator {
		overridePerms[do.PerModManageModerator] = true
	}
	if req.CanBanUser != nil && *req.CanBanUser {
		overridePerms[do.PerModBanUser] = true
	}

	// 最终权限：优先使用审批人的设置，否则采用申请时的请求
	var final []do.ModeratorPermission
	// 从申请时的请求权限出发
	for _, perm := range app.RequestedPermissions {
		if overridePerms[perm] {
			final = append(final, perm) // 审批人明确允许
		} else if !isPermissionOverridden(perm, req) {
			// 如果审批人没有明确禁止（未设置），则保留申请时的权限
			// 注意：这里假设请求结构体中每个权限字段是 *bool，nil 表示未覆盖
			final = append(final, perm)
		}
	}
	// 额外添加审批人单独开启但申请时未请求的权限
	for perm, granted := range overridePerms {
		if granted && !containsPermission(final, perm) {
			final = append(final, perm)
		}
	}
	return final
}

// isPermissionOverridden 判断审批人是否明确覆盖了该权限（即字段不为 nil）
func isPermissionOverridden(perm do.ModeratorPermission, req request.ReviewApplicationRequest) bool {
	switch perm {
	case do.PerModDeletePost:
		return req.CanDeletePost != nil
	case do.PerMoePinPost:
		return req.CanPinPost != nil
	case do.PerModEditAnyPost:
		return req.CanEditAnyPost != nil
	case do.PerModManageModerator:
		return req.CanManageModerator != nil
	case do.PerModBanUser:
		return req.CanBanUser != nil
	}
	return false
}

// containsPermission 辅助函数
func containsPermission(perms []do.ModeratorPermission, target do.ModeratorPermission) bool {
	for _, p := range perms {
		if p == target {
			return true
		}
	}
	return false
}

// sendRejectionNotification 发送拒绝通知
func (s *boardService) sendRejectionNotification(app *do.ModeratorApplication, reviewerID uint, note string) {
	s.notifSvc.Create(
		app.UserID,
		&reviewerID,
		do.NotifySystem,
		fmt.Sprintf("你的版主申请已被拒绝：%s", note),
		&app.BoardID,
		"board",
	)
}

// sendApprovalNotification 发送通过通知
func (s *boardService) sendApprovalNotification(app *do.ModeratorApplication, reviewerID uint) {
	s.notifSvc.Create(
		app.UserID,
		&reviewerID,
		do.NotifySystem,
		"恭喜！你的版主申请已通过",
		&app.BoardID,
		"board",
	)
}

// ListApplications 版主申请列表查询（支持板块筛选、状态筛选、分页）
func (s *boardService) ListApplications(boardID *uint, status do.ApplicationStatus, page, pageSize int) ([]do.ModeratorApplication, int64, error) {
	// 参数规范化
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	// 可选：验证 status 有效性
	if status != "" && !status.IsValid() {
		return nil, 0, fmt.Errorf("无效的申请状态: %s", status)
	}
	return s.boardRepo.ListApplications(boardID, status, page, pageSize)
}
