// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"

	"power-admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConfigLogic {
	return &DeleteConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConfigLogic) DeleteConfig() error {
	id := l.ctx.Value("id").(int64)
	return l.svcCtx.CodegenRepo.DeleteConfig(l.ctx, id)
}
