// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	menu, err := l.svcCtx.MenuRepo.GetByID(req.Id)
	if err != nil || menu == nil {
		return err
	}

	if req.MenuName != "" {
		menu.MenuName = req.MenuName
	}
	if req.MenuPath != "" {
		menu.MenuPath = req.MenuPath
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Sort > 0 {
		menu.Sort = req.Sort
	}
	if req.Status > 0 {
		menu.Status = req.Status
	}
	if req.MenuType > 0 {
		menu.MenuType = req.MenuType
	}

	err = l.svcCtx.MenuRepo.Update(menu)
	if err != nil {
		return nil
	}

	return nil
}
