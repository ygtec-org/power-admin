package role

import (
	"context"
	"fmt"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenuAndApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleMenuAndApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleMenuAndApiLogic {
	return &GetRoleMenuAndApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleMenuAndApiLogic) GetRoleMenuAndApi(req *types.GetRolePermissionsReq) (resp *types.GetRoleMenuAndApiResp, err error) {
	roleID := req.RoleID
	roleIDStr := fmt.Sprintf("%d", roleID)

	// 获取所有菜单
	menus, _, err := l.svcCtx.MenuRepo.List(0, 0, 0)
	if err != nil {
		l.Errorf("Failed to get menus: %v", err)
		return nil, err
	}

	// 获取所有API
	apis, _, err := l.svcCtx.APIRepo.List(1, 1000)
	if err != nil {
		l.Errorf("Failed to get apis: %v", err)
		return nil, err
	}

	// 从casbin_rule表中查询该角色已有的菜单权限
	// 从casbin_rule表中查询该角色已有的菜单权限是是从casbin版版表或role_menus表
	var roleMenus []models.RoleMenu
	l.svcCtx.DB.Where("role_id = ?", roleID).Find(&roleMenus)
	l.Infof("Role menus for role %d: %v", roleID, roleMenus)

	// 构建MenuID集合以便快速查找
	selectedMenuIdMap := make(map[int64]bool)
	for _, roleMenu := range roleMenus {
		selectedMenuIdMap[roleMenu.MenuID] = true
	}

	// 从casbin_rule表中查询该角色已有的API权限
	// API权限格式: (p, roleId, apiPath, httpMethod)
	var apiRules []models.CasbinRule
	l.svcCtx.DB.Where("ptype = ? AND v0 = ?", "p", roleIDStr).Find(&apiRules)
	l.Infof("API rules for role %s: %v, Total count: %d", roleIDStr, apiRules, len(apiRules))

	// 构建API路径:method的map
	selectedApiPaths := make(map[string]bool)
	for _, rule := range apiRules {
		key := rule.V1 + ":" + rule.V2
		selectedApiPaths[key] = true
		l.Infof("API rule found: v0=%s, v1=%s, v2=%s, key=%s", rule.V0, rule.V1, rule.V2, key)
	}

	// 构建菜单列表和已选菜单ID
	menuList := make([]types.RolePermissionItem, 0)
	selectedMenuIds := make([]int64, 0)
	for _, menu := range menus {
		item := types.RolePermissionItem{
			Id:       menu.ID,
			Name:     menu.MenuName,
			Path:     menu.MenuPath,
			Type:     "menu",
			ParentId: menu.ParentID,
		}
		menuList = append(menuList, item)

		// 如果该菜单在role_menus表中有记录，则添加到已选列表
		if selectedMenuIdMap[menu.ID] {
			selectedMenuIds = append(selectedMenuIds, menu.ID)
		}
	}

	// 构建API列表和已选API ID
	apiList := make([]types.RolePermissionItem, 0)
	selectedApiIds := make([]int64, 0)
	for _, api := range apis {
		item := types.RolePermissionItem{
			Id:     api.ID,
			Name:   api.APIName,
			Path:   api.APIPath,
			Type:   "api",
			Method: api.APIMethod,
			Group:  api.Group,
		}
		apiList = append(apiList, item)

		// 如果该API在casbin_rule中有权限记录，则添加到已选列表
		if selectedApiPaths[api.APIPath+":"+api.APIMethod] {
			selectedApiIds = append(selectedApiIds, api.ID)
		}
	}

	resp = &types.GetRoleMenuAndApiResp{
		Menus:           menuList,
		Apis:            apiList,
		SelectedMenuIds: selectedMenuIds,
		SelectedApiIds:  selectedApiIds,
	}
	l.Infof("Response: SelectedMenuIds=%v, SelectedApiIds=%v", selectedMenuIds, selectedApiIds)
	return
}
