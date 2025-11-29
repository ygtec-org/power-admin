// 用户角色分配逻辑
package user

import (
	"context"
	"fmt"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRolesToUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRolesToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesToUserLogic {
	return &AssignRolesToUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRolesToUserLogic) AssignRolesToUser(req *types.AssignRolesToUserReq) (resp *types.AssignRolesToUserResp, err error) {
	userID := req.UserID
	userIDStr := fmt.Sprintf("%d", userID)

	// 清除用户在Casbin中的所有角色关联
	// 获取用户当前所有角色
	oldRoles, err := l.svcCtx.Permission.GetRolesForUser(userIDStr)
	if err == nil {
		// 移除用户所有旧角色
		for _, role := range oldRoles {
			l.svcCtx.Permission.RemoveRoleForUser(userIDStr, role)
		}
	}

	// 为用户分配新角色到Casbin
	for _, roleID := range req.RoleIds {
		roleIDStr := fmt.Sprintf("%d", roleID)
		if err := l.svcCtx.Permission.AddRoleForUser(userIDStr, roleIDStr); err != nil {
			l.Errorf("Failed to add role %s to user %s: %v", roleIDStr, userIDStr, err)
			continue
		}
	}

	resp = &types.AssignRolesToUserResp{
		Success: true,
		Message: "角色分配成功，已更新Casbin规则表",
	}
	return
}
