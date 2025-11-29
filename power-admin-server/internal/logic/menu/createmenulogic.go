// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"
	"power-admin-server/pkg/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (err error) {
	if req.MenuName == "" {
		return nil
	}

	menu := &models.Menu{
		ParentID:  req.ParentId,
		MenuName:  req.MenuName,
		MenuPath:  req.MenuPath,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		MenuType:  req.MenuType,
		Status:    1,
	}

	err = l.svcCtx.MenuRepo.Create(menu)
	if err != nil {
		return nil
	}

	return nil
}
