// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"fmt"
	"strconv"

	"power-admin-server/common/constant"
	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuListLogic {
	return &MenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuListLogic) MenuList(req *types.MenuListReq) (resp *types.MenuListResp, err error) {
	// 从上下文中获取用户ID
	userIDStr := l.ctx.Value(constant.AdminUserKey)
	if userIDStr == nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	// 获取用户的角色
	userRoles, err := l.svcCtx.RoleRepo.GetRolesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 如果用户没有角色，返回空菜单
	if len(userRoles) == 0 {
		resp = &types.MenuListResp{
			Data:  []types.MenuItem{},
			Total: 0,
		}
		return
	}

	// 检查是否是超级管理员（role_id = 1）
	isSuperAdmin := false
	var menus []models.Menu

	for _, role := range userRoles {
		if role.ID == 1 {
			isSuperAdmin = true
			break
		}
	}

	if isSuperAdmin {
		// 超级管理员返回所有菜单
		allMenus, _, err := l.svcCtx.MenuRepo.All(0)
		if err != nil {
			return nil, fmt.Errorf("failed to get all menus: %w", err)
		}
		menus = allMenus
	} else {
		// 普通用户：提取角色ID列表
		roleIDs := make([]int64, len(userRoles))
		for i, role := range userRoles {
			roleIDs[i] = role.ID
		}

		// 获取用户角色绑定的菜单
		var err2 error
		menus, err2 = l.svcCtx.MenuRepo.GetMenusByRoleIDs(roleIDs)
		if err2 != nil {
			return nil, fmt.Errorf("failed to get menus: %w", err2)
		}
	}

	// 转换为 MenuItem 并构建树形结构
	menuMap := make(map[int64]*types.MenuItem)
	var rootMenus []*types.MenuItem

	for _, m := range menus {
		menuItem := &types.MenuItem{
			Id:        m.ID,
			ParentId:  m.ParentID,
			MenuName:  m.MenuName,
			MenuPath:  m.MenuPath,
			Component: m.Component,
			Icon:      m.Icon,
			Sort:      m.Sort,
			Status:    m.Status,
			MenuType:  m.MenuType,
			CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
			Children:  make([]types.MenuItem, 0),
		}
		menuMap[m.ID] = menuItem

		// 如果是根菜单
		if m.ParentID == 0 {
			rootMenus = append(rootMenus, menuItem)
		}
	}

	// 构建树形结构
	for _, m := range menus {
		if m.ParentID != 0 {
			if parent, ok := menuMap[m.ParentID]; ok {
				parent.Children = append(parent.Children, *menuMap[m.ID])
			}
		}
	}

	// 转换指针为值
	data := make([]types.MenuItem, 0, len(rootMenus))
	for _, m := range rootMenus {
		data = append(data, *m)
	}

	resp = &types.MenuListResp{
		Data:  data,
		Total: int64(len(menus)),
	}
	return
}
