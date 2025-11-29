// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iface

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApiLogic) UpdateApi(req *types.UpdateApiReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	api, err := l.svcCtx.APIRepo.GetByID(req.Id)
	if err != nil || api == nil {
		return err
	}

	if req.Status > 0 {
		api.Status = req.Status
	}

	err = l.svcCtx.APIRepo.Update(api)
	if err != nil {
		return nil
	}

	return nil
}
