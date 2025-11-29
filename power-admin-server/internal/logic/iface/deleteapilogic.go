// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iface

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApiLogic) DeleteApi(req *types.DeleteApiReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	err = l.svcCtx.APIRepo.Delete(req.Id)
	if err != nil {
		return nil
	}

	return nil
}
