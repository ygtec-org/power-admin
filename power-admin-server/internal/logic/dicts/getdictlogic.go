// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dicts

import (
	"context"

	"power-admin-server/internal/svc"
	"power-admin-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictLogic {
	return &GetDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictLogic) GetDict(req *types.GetDictReq) (resp *types.DictInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
