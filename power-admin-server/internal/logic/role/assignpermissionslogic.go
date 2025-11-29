// 角色权限分配逻辑
package role

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignPermissionsLogic {
	return &AssignPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignPermissionsLogic) AssignPermissions(req *types.AssignPermissionsReq) (resp *types.AssignPermissionsResp, err error) {
	// 清除该角色的所有权限关联
	roleID := req.RoleID

	// 删除该角色的所有权限关联
	if err := l.svcCtx.RoleRepo.RemoveAllPermissions(roleID); err != nil {
		l.Errorf("Failed to remove permissions for role %d: %v", roleID, err)
		return nil, err
	}

	// 为角色分配新的权限
	for _, permID := range req.PermissionIds {
		if err := l.svcCtx.RoleRepo.AddPermission(roleID, permID); err != nil {
			l.Errorf("Failed to add permission %d to role %d: %v", permID, roleID, err)
			return nil, err
		}
	}

	resp = &types.AssignPermissionsResp{
		Success: true,
		Message: "权限分配成功",
	}
	return
}
