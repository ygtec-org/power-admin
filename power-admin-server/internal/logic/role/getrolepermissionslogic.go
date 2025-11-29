// 获取角色权限逻辑
package role

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionsLogic {
	return &GetRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePermissionsLogic) GetRolePermissions(req *types.GetRolePermissionsReq) (resp *types.GetRolePermissionsResp, err error) {
	roleID := req.RoleID

	// 获取该角色的所有权限
	permissions, err := l.svcCtx.RoleRepo.GetPermissions(roleID)
	if err != nil {
		l.Errorf("Failed to get permissions for role %d: %v", roleID, err)
		return nil, err
	}

	data := make([]types.PermissionInfo, 0)
	if permissions != nil {
		for _, p := range permissions {
			data = append(data, types.PermissionInfo{
				Id:          p.ID,
				Name:        p.Name,
				Resource:    p.Resource,
				Action:      p.Action,
				Description: p.Description,
				Status:      p.Status,
			})
		}
	}

	resp = &types.GetRolePermissionsResp{
		Data: data,
	}
	return
}
