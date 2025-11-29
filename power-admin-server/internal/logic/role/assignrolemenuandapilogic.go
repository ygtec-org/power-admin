package role

import (
	"context"
	"fmt"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleMenuAndApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRoleMenuAndApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleMenuAndApiLogic {
	return &AssignRoleMenuAndApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleMenuAndApiLogic) AssignRoleMenuAndApi(req *types.AssignRoleMenuAndApiReq) (resp *types.AssignRoleMenuAndApiResp, err error) {
	roleID := req.RoleId
	menuIds := req.MenuIds
	apiIds := req.ApiIds

	// 清除该角色在role_menus表中的所有菜单关联
	if err := l.svcCtx.DB.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
		l.Errorf("Failed to clear role menus for role %d: %v", roleID, err)
		return nil, err
	}

	// 清除该角色在casbin_rule表中的所有API规则
	if err := l.clearRoleCasbinRules(roleID); err != nil {
		l.Errorf("Failed to clear casbin rules for role %d: %v", roleID, err)
		return nil, err
	}

	// 为菜单添加role_menus关联
	for _, menuID := range menuIds {
		roleMenu := &models.RoleMenu{
			RoleID: roleID,
			MenuID: menuID,
		}
		if err := l.svcCtx.DB.Create(roleMenu).Error; err != nil {
			l.Errorf("Failed to insert role menu: %v", err)
			continue
		}
	}

	// 为API生成Casbin规则
	// API使用 (p, roleId, apiPath, httpMethod) 格式
	for _, apiID := range apiIds {
		api, err := l.svcCtx.APIRepo.GetByID(apiID)
		if err != nil {
			l.Errorf("Failed to get API %d: %v", apiID, err)
			continue
		}

		if api != nil && api.APIPath != "" {
			// API权限规则：(p, 角色ID, API路径, HTTP方法)
			rule := &models.CasbinRule{
				PType: "p",
				V0:    fmt.Sprintf("%d", roleID), // 角色ID
				V1:    api.APIPath,               // API路径
				V2:    api.APIMethod,             // HTTP方法
			}

			if err := l.svcCtx.DB.Omit("id").Create(rule).Error; err != nil {
				l.Errorf("Failed to insert api casbin rule: %v", err)
				continue
			}
		}
	}

	resp = &types.AssignRoleMenuAndApiResp{
		Success: true,
		Message: "权限分配成功，已更新role_menus和casbin_rule表",
	}
	return
}

// clearRoleCasbinRules 清除角色在casbin_rule表中的所有规则
func (l *AssignRoleMenuAndApiLogic) clearRoleCasbinRules(roleID int64) error {
	return l.svcCtx.DB.Where("ptype = ? AND v0 = ?", "p", fmt.Sprintf("%d", roleID)).Delete(&models.CasbinRule{}).Error
}
