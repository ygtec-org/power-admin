// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dicts

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictLogic {
	return &DeleteDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictLogic) DeleteDict(req *types.DeleteDictReq) (err error) {
	if req.Id == 0 {
		return nil
	}

	err = l.svcCtx.DictRepo.Delete(req.Id)
	if err != nil {
		return nil
	}

	return nil
}
