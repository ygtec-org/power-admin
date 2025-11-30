// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

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
	// 计算 offset
	//offset := (req.Page - 1) * req.PageSize

	// 获取分页菜单数据
	menus, total, err := l.svcCtx.MenuRepo.All(req.ParentId)
	if err != nil {
		return nil, err
	}

	// 转换为 MenuItem
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
		Total: total, // 使用真实的总数量
	}
	return
}
