// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package codegen

import (
	"context"

	"power-admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHistoryLogic {
	return &DeleteHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteHistoryLogic) DeleteHistory() error {
	id := l.ctx.Value("id").(int64)
	return l.svcCtx.CodegenRepo.DeleteHistory(l.ctx, id)
}
